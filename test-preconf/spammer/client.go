package spammer

import (
	"context"
	"errors"
	"math/big"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	ctx context.Context
	*ethclient.Client
	chainID *big.Int
	logger  *log.Logger
}

func NewEthClient(ctx context.Context, url string, chainID *big.Int, logger *log.Logger) (*EthClient, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &EthClient{ctx, client, chainID, logger}, nil
}

func (ec *EthClient) FetchAssignedSlots() ([]uint64, error) {
	var assignedSlots []uint64
	err := ec.Client.Client().CallContext(ec.ctx, &assignedSlots, "taiko_fetchAssignedSlots")
	if err != nil {
		return nil, err
	}
	return assignedSlots, nil
}

func (ec *EthClient) FetchL1GenesisTimestamp() (uint64, error) {
	var l1GenesisTimestamp uint64
	err := ec.Client.Client().CallContext(ec.ctx, &l1GenesisTimestamp, "taiko_fetchL1GenesisTimestamp")
	if err != nil {
		return 0, err
	}
	return l1GenesisTimestamp, nil
}

func (ec *EthClient) FetchCurrentSlot(now int64) (uint64, uint64, error) {
	l1GenesisTimestamp, err := ec.FetchL1GenesisTimestamp()
	if err != nil {
		return 0, 0, err
	}
	currentSlot, _ := common.CurrentSlotAndEpoch(l1GenesisTimestamp, now)
	_, currentSlotEndTime := common.CurrentSlotStartEndTime(l1GenesisTimestamp, now)
	return currentSlot, currentSlotEndTime, nil
}

func (ec *EthClient) GetNonce(account *Account) (uint64, error) {
	nonce, err := ec.PendingNonceAt(ec.ctx, *account.Address())
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (ec *EthClient) SendTx(account *Account, tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewPreconfSigner(ec.chainID), account.PrivateKey())
	if err != nil {
		ec.logger.Error("Failed to sign transaction", "error", err)
		return nil, err
	}

	res, err := ec.SendTransactionWithResult(ec.ctx, signedTx)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(res, common.Hash{}) {
		ec.logger.Error("Tx rejected", "hash", res)
	} else {
		ec.logger.Info("Submitted tx", "hash", signedTx.Hash(), "nonce", signedTx.Nonce(), "slot deadline", signedTx.Deadline().Uint64())
	}

	return signedTx, nil
}

func (ec *EthClient) LogTx(signedTx *types.Transaction) {
	_, _, err := ec.TransactionByHash(ec.ctx, signedTx.Hash())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			ec.logger.Error("Transaction not found", "tx hash", signedTx.Hash())
			return
		} else {
			ec.logger.Error("Failed to get transaction by hash", "error", err, "tx hash", signedTx.Hash())
		}
	}
}

func (ec *EthClient) LogReceipt(signedTx *types.Transaction) {
	txReceipt, err := ec.TransactionReceipt(ec.ctx, signedTx.Hash())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			ec.logger.Error("Transaction receipt not found", "tx hash ", signedTx.Hash())
			return
		} else {
			ec.logger.Error("Failed to get transaction receipt", "error", err)
		}
	}
	ec.logger.Warn("Transaction receipt", "tx hash", signedTx.Hash(), "block number", txReceipt.BlockNumber, "status", txReceipt.Status, "cumulative gas used", txReceipt.CumulativeGasUsed, "effective gas price", txReceipt.EffectiveGasPrice, "gas used", txReceipt.GasUsed)
}
