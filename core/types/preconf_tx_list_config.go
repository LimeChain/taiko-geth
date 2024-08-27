package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// CHANGE(limechain):

type TxListConfig struct {
	Beneficiary          common.Address // L1 proposer address
	BaseFee              *big.Int       // base fee calculated in the protocol contract
	BlockMaxGasLimit     uint64         // hard-coded in the protocol contract config
	MaxBytesPerTxList    uint64
	MaxTransactionsLists uint64
	Locals               []string // TODO(limechain): []common.Address
}
