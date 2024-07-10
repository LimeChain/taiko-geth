package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (ec *Client) HeadL1Origin(ctx context.Context) (*rawdb.L1Origin, error) {
	var res *rawdb.L1Origin

	if err := ec.c.CallContext(ctx, &res, "taiko_headL1Origin"); err != nil {
		return nil, err
	}

	return res, nil
}

// L1OriginByID returns the L2 block's corresponding L1 origin.
func (ec *Client) L1OriginByID(ctx context.Context, blockID *big.Int) (*rawdb.L1Origin, error) {
	var res *rawdb.L1Origin

	if err := ec.c.CallContext(ctx, &res, "taiko_l1OriginByID", hexutil.EncodeBig(blockID)); err != nil {
		return nil, err
	}

	return res, nil
}

// GetSyncMode returns the current sync mode of the L2 node.
func (ec *Client) GetSyncMode(ctx context.Context) (string, error) {
	var res string

	if err := ec.c.CallContext(ctx, &res, "taiko_getSyncMode"); err != nil {
		return "", err
	}

	return res, nil
}

// GetPreconfBlockCursor returns the current preconf block cursor.
func (ec *Client) GetPreconfBlockCursor(ctx context.Context) (*types.PreconfBlockCursor, error) {
	log.Warn("Client get")
	var res *types.PreconfBlockCursor
	if err := ec.c.CallContext(ctx, &res, "taiko_getPreconfBlockCursor"); err != nil {
		return nil, err
	}
	return res, nil
}

// UpdatePreconfBlockCursor updates the current preconf block cursor.
func (ec *Client) UpdatePreconfBlockCursor(ctx context.Context, hash *common.Hash, number *big.Int, proposedTxCount *big.Int, skipProposedTx *bool) error {
	log.Warn("Client update")
	var res bool
	if err := ec.c.CallContext(ctx, &res, "taiko_updatePreconfBlockCursor",
		hash.Hex(),
		hexutil.EncodeBig(number),
		hexutil.EncodeBig(proposedTxCount),
		skipProposedTx,
	); err != nil {
		return err
	}
	return nil
}

// TODO: not needed, remove it
// DeletePreconfBlockCursor deletes the current preconf block cursor.
func (ec *Client) DeletePreconfBlockCursor(ctx context.Context) error {
	log.Warn("Client delete")
	var res bool
	if err := ec.c.CallContext(ctx, &res, "taiko_deletePreconfBlockCursor"); err != nil {
		return err
	}
	return nil
}
