package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// Cursor that keeps track of the latest virtual block used for preconfirmed txs
type PreconfBlockCursor struct {
	Hash            common.Hash
	Number          uint64
	ProposedTxCount uint64
	SkipExecutedTx  bool
}
