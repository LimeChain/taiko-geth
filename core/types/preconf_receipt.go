package types

import (
	"github.com/ethereum/go-ethereum/common"
)

// Pre-confirmation receipts are similar to normal receipts, but
// they also contain additional fields that are not expected to
// be resolved later on and the storage is also different.

// PreconfReceipt extends Receipt by including From and To fields.
type PreconfReceipt struct {
	Receipt
	From common.Address
	To   common.Address
}
