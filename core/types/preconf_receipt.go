package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// Pre-confirmation receipts are similar to normal receipts, but
// they also contain additional fields that are not expected to
// be resolved later on and the storage is also different.

type PreconfReceipt struct {
	Receipt
	// Additional fields that can not be resolved later
	From common.Address
	To   common.Address
}
