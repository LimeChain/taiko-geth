package rawdb

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

const anchorTxIndexOffset = 1

// Preconf receipts are stored under the tx hash to
// allow for retrieval without canonical block data.

func ReadPreconfReceipt(db ethdb.Database, txHash common.Hash) *types.PreconfReceipt {
	data, _ := db.Get(preconfTxReceiptKey(txHash.Bytes()))
	if len(data) == 0 {
		return nil
	}

	receipt := new(types.PreconfReceipt)
	if err := rlp.Decode(bytes.NewBuffer(data), receipt); err != nil {
		log.Error("Invalid preconf tx receipt RLP", "err", err)
		return nil
	}

	return receipt
}

func WritePreconfReceipt(db ethdb.Database, receipt *types.Receipt, from *common.Address, to *common.Address) {
	data := bytes.NewBuffer(nil)

	preconfReceipt := &types.PreconfReceipt{
		Receipt: *receipt,
		From:    *from,
		To:      *to,
	}
	// There is anchor tx expected at the beginning of each block,
	// so the tx index is offset by 1.
	preconfReceipt.Receipt.TransactionIndex += anchorTxIndexOffset

	err := rlp.Encode(data, preconfReceipt)
	if err != nil {
		log.Crit("Failed to RLP encode preconf tx receipt", "err", err)
	}

	if err := db.Put(preconfTxReceiptKey(receipt.TxHash.Bytes()), data.Bytes()); err != nil {
		log.Crit("Failed to store preconf tx receipt", "err", err)
	}
}
