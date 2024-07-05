package preconfs

import (
	"context"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/bindings/encoding"
	"github.com/ethereum/go-ethereum/common"
	consensus "github.com/ethereum/go-ethereum/consensus/taiko"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

// AnchorTxConstructor is responsible for assembling the anchor transaction (TaikoL2.anchor) in
// each L2 block, which must be the first transaction, and its sender must be the golden touch account.
type AnchorTxConstructor struct {
	signer           *FixedKSigner
	rpc              *Client
	l2TransactionAPI EthTransactionAPI
}

// New creates a new AnchorConstructor instance.
func NewAnchorTxConstructor(rpc *Client, l2TransactionAPI EthTransactionAPI) (*AnchorTxConstructor, error) {
	signer, err := NewFixedKSigner("0x" + encoding.GoldenTouchPrivKey)
	if err != nil {
		return nil, fmt.Errorf("invalid golden touch private key %s", encoding.GoldenTouchPrivKey)
	}
	return &AnchorTxConstructor{
		signer,
		rpc,
		l2TransactionAPI,
	}, nil
}

// AssembleAnchorTx assembles a signed TaikoL2.anchor transaction.
func (c *AnchorTxConstructor) AssembleAnchorTx(
	ctx context.Context,
	l1Height *big.Int,
	l1Hash common.Hash,
	l2Height *big.Int, // height of the L2 block which includes the TaikoL2.anchor transaction
	baseFee *big.Int,
	parentGasUsed uint64,
) (*types.Transaction, error) {
	opts, err := c.transactOpts(ctx, l2Height, baseFee)
	if err != nil {
		return nil, err
	}

	l1Header, err := c.rpc.L1.HeaderByHash(ctx, l1Hash)
	if err != nil {
		return nil, err
	}

	// log.Info(
	// 	"Anchor arguments",
	// 	"l2Height", l2Height,
	// 	"l1Height", l1Height,
	// 	"l1Hash", l1Hash,
	// 	"stateRoot", l1Header.Root,
	// 	"baseFee", utils.WeiToGWei(baseFee),
	// 	"gasUsed", parentGasUsed,
	// )

	return c.rpc.TaikoL2.Anchor(opts, l1Hash, l1Header.Root, l1Height.Uint64(), uint32(parentGasUsed))
}

// transactOpts is a utility method to create some transact options of the anchor transaction in given L2 block with
// golden touch account's private key.
func (c *AnchorTxConstructor) transactOpts(
	ctx context.Context,
	l2Height *big.Int,
	baseFee *big.Int,
) (*bind.TransactOpts, error) {
	var (
		signer       = types.LatestSignerForChainID(c.rpc.L2.chainID)
		parentHeight = new(big.Int).Sub(l2Height, common.Big1)
	)

	// Get the nonce of golden touch account at the specified parentHeight.
	nonce, err := c.AccountNonce(ctx, consensus.GoldenTouchAccount, parentHeight)
	if err != nil {
		return nil, err
	}

	log.Info(
		"Golden touch account nonce",
		"address", consensus.GoldenTouchAccount,
		"nonce", nonce,
		"parent", parentHeight,
	)

	return &bind.TransactOpts{
		From: consensus.GoldenTouchAccount,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != consensus.GoldenTouchAccount {
				return nil, bind.ErrNotAuthorized
			}
			signature, err := c.signTxPayload(signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
		Nonce:     new(big.Int).SetUint64(nonce),
		Context:   ctx,
		GasFeeCap: baseFee,
		GasTipCap: common.Big0,
		GasLimit:  consensus.AnchorGasLimit,
		NoSend:    true,
	}, nil
}

// signTxPayload calculates an ECDSA signature for an anchor transaction.
func (c *AnchorTxConstructor) signTxPayload(hash []byte) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}

	// Try k = 1.
	sig, ok := c.signer.SignWithK(new(secp256k1.ModNScalar).SetInt(1))(hash)
	if !ok {
		// Try k = 2.
		sig, ok = c.signer.SignWithK(new(secp256k1.ModNScalar).SetInt(2))(hash)
		if !ok {
			log.Crit("Failed to sign TaikoL2.anchor transaction using K = 1 and K = 2")
		}
	}

	return sig[:], nil
}

// AccountNonce fetches the nonce of the given L2 account at a specified height.
func (c *AnchorTxConstructor) AccountNonce(ctx context.Context, account common.Address, height *big.Int) (uint64, error) {
	blockNumber := rpc.BlockNumber(height.Int64())

	result, err := c.l2TransactionAPI.GetTransactionCount(ctx, account, rpc.BlockNumberOrHash{BlockNumber: &blockNumber})
	if err != nil {
		return 0, err
	}

	return uint64(*result), err
}
