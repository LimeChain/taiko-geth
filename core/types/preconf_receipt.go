package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// PreconfReceipt extends Receipt by including additional From and To fields.
type PreconfReceipt struct {
	Receipt
	From common.Address
	To   common.Address
}
