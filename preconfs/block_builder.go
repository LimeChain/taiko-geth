package preconfs

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/taiko"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

// implemented by ethapi.Backend
type EthBackendAPI interface {
	CurrentHeader() *types.Header
	ChainDb() ethdb.Database
}

// implemented by catalyst.ConsensusAPI
type EthEngineAPI interface {
	ForkchoiceUpdatedV2(update engine.ForkchoiceStateV1, params *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)
	GetPayloadV2(payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)
	NewPayloadV2(params engine.ExecutableData) (engine.PayloadStatusV1, error)
}

// implemented by ethapi.TransactionAPI
type EthTransactionAPI interface {
	GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error)
}

type BlockBuilder struct {
	backendAPI     EthBackendAPI
	transactionAPI EthTransactionAPI
	engineAPI      EthEngineAPI
	rpc            *Client // for L1 and L2 conrtact calls
}

func NewBlockBuilder(
	ctx context.Context,
	backendAPI EthBackendAPI,
	transactionAPI EthTransactionAPI,
	engineAPI EthEngineAPI,
	taikoL2Address common.Address,
	l1Endpoint string,
	l2Endpoint string,
	timeout time.Duration,
) (*BlockBuilder, error) {
	l1Client, err := NewEthClient(ctx, l1Endpoint, timeout)
	if err != nil {
		log.Error("Failed to connect to L1 endpoint, retrying", "endpoint", l1Endpoint, "err", err)
		return nil, err
	}

	l2Client, err := NewEthClient(ctx, l2Endpoint, timeout)
	if err != nil {
		log.Error("Failed to connect to L2 endpoint, retrying", "endpoint", l2Endpoint, "err", err)
		return nil, err
	}
	taikoL2, err := bindings.NewTaikoL2Client(taikoL2Address, l2Client)
	if err != nil {
		return nil, err
	}

	return &BlockBuilder{
		backendAPI:     backendAPI,
		transactionAPI: transactionAPI,
		engineAPI:      engineAPI,
		rpc: &Client{
			L1:      l1Client,
			L2:      l2Client,
			TaikoL2: taikoL2,
		},
	}, nil
}

// BuildVirtualBlock builds virtual block that provides pre-confirmation receipts for the contained TXs.
func (b *BlockBuilder) BuildVirtualBlock(ctx context.Context, txList types.Transactions) error {
	head := b.backendAPI.CurrentHeader()
	log.Info("Preconfs: current head", "number", head.Number, "hash", head.Hash())

	// parameters of the TaikoL2.anchor transaction
	l1Origin, err := b.HeadL1Origin()
	if err != nil {
		return err
	}
	l1Height := l1Origin.L1BlockHeight
	l1Hash := l1Origin.L1BlockHash

	baseFeeInfo, err := b.rpc.TaikoL2.GetBasefee(
		&bind.CallOpts{BlockNumber: head.Number, Context: ctx},
		l1Height.Uint64(),
		uint32(head.GasUsed),
	)
	if err != nil {
		return fmt.Errorf("failed to get L2 baseFee")
	}

	anchorConstructor, err := NewAnchorTxConstructor(b.rpc, b.transactionAPI)
	if err != nil {
		return err
	}
	anchorTx, err := anchorConstructor.AssembleAnchorTx(
		ctx,
		l1Height,
		l1Hash,
		new(big.Int).Add(head.Number, common.Big1),
		baseFeeInfo.Basefee,
		head.GasUsed,
	)
	if err != nil {
		return fmt.Errorf("failed to create TaikoL2.anchor transaction: %w", err)
	}
	log.Info("Anchor tx", "hash", anchorTx.Hash().String())

	// Insert the anchor transaction at the head of the transactions list
	txList = append([]*types.Transaction{anchorTx}, txList...)
	txListBytes, err := rlp.EncodeToBytes(txList)
	if err != nil {
		return fmt.Errorf("failed to encode transactions: %w", err)
	}

	fc := &engine.ForkchoiceStateV1{HeadBlockHash: head.Hash()}

	timestamp := uint64(time.Now().Unix())
	coinbase := common.Address{}

	attributes := &engine.PayloadAttributes{
		Timestamp:             timestamp,
		Random:                common.Hash{},
		SuggestedFeeRecipient: coinbase,
		Withdrawals:           make(types.Withdrawals, 0),
		BlockMetadata: &engine.BlockMetadata{
			HighestBlockID: head.Number,
			Beneficiary:    coinbase,
			GasLimit:       uint64(21000) + taiko.AnchorGasLimit,
			Timestamp:      timestamp,
			TxList:         txListBytes,
			MixHash:        common.Hash{},
			ExtraData:      []byte{},
		},
		BaseFeePerGas: baseFeeInfo.Basefee,
		L1Origin: &rawdb.L1Origin{
			BlockID:       head.Number,
			L2BlockHash:   common.Hash{}, // Will be set by taiko-geth.
			L1BlockHeight: l1Height,
			L1BlockHash:   l1Hash,
		},
		VirtualBlock: true,
	}

	// Start building payload
	fcRes, err := b.engineAPI.ForkchoiceUpdatedV2(*fc, attributes)
	if err != nil {
		return fmt.Errorf("failed to update fork choice: %w", err)
	}
	if fcRes.PayloadStatus.Status != engine.VALID {
		return fmt.Errorf("unexpected ForkchoiceUpdate response status: %s", fcRes.PayloadStatus.Status)
	}
	if fcRes.PayloadID == nil {
		return errors.New("empty payload ID")
	}
	log.Info("Preconfs: started building payload")

	// Get the built payload
	payload, err := b.engineAPI.GetPayloadV2(*fcRes.PayloadID)
	if err != nil {
		log.Info("Preconfs: failed to get payload")
		return fmt.Errorf("failed to get payload: %w", err)
	}
	log.Info("Preconfs: get built payload")

	// Execute the payload
	execStatus, err := b.engineAPI.NewPayloadV2(*payload.ExecutionPayload)
	if err != nil {
		log.Info("Preconfs: failed to create a new payload")
		return fmt.Errorf("failed to create a new payload: %w", err)
	}
	if execStatus.Status != engine.VALID {
		return fmt.Errorf("unexpected NewPayload response status: %s", execStatus.Status)
	}
	log.Info("Preconfs: executed payload", "hash", payload.ExecutionPayload.BlockHash.String(), "txs", len(payload.ExecutionPayload.Transactions))

	// fc = &engine.ForkchoiceStateV1{
	// 	HeadBlockHash:      payload.ExecutionPayload.BlockHash,
	// 	SafeBlockHash:      payload.ExecutionPayload.BlockHash,
	// 	FinalizedBlockHash: payload.ExecutionPayload.BlockHash,
	// }

	// // Update the fork choice
	// fcRes, err = pb.engineAPI.ForkchoiceUpdatedV2(*fc, nil)
	// if err != nil {
	// 	return err
	// }
	// if fcRes.PayloadStatus.Status != engine.VALID {
	// 	return fmt.Errorf("unexpected ForkchoiceUpdate response status: %s", fcRes.PayloadStatus.Status)
	// }

	return nil
}

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (b *BlockBuilder) HeadL1Origin() (*rawdb.L1Origin, error) {
	blockID, err := rawdb.ReadHeadL1Origin(b.backendAPI.ChainDb())
	if err != nil {
		return nil, err
	}

	if blockID == nil {
		return nil, ethereum.NotFound
	}

	l1Origin, err := rawdb.ReadL1Origin(b.backendAPI.ChainDb(), blockID)
	if err != nil {
		return nil, err
	}

	if l1Origin == nil {
		return nil, ethereum.NotFound
	}

	return l1Origin, nil
}
