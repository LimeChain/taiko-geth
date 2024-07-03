// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package engine

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
)

var _ = (*payloadAttributesMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (p PayloadAttributes) MarshalJSON() ([]byte, error) {
	type PayloadAttributes struct {
		Timestamp             hexutil.Uint64      `json:"timestamp"             gencodec:"required"`
		Random                common.Hash         `json:"prevRandao"            gencodec:"required"`
		SuggestedFeeRecipient common.Address      `json:"suggestedFeeRecipient" gencodec:"required"`
		Withdrawals           []*types.Withdrawal `json:"withdrawals"`
		BeaconRoot            *common.Hash        `json:"parentBeaconBlockRoot"`
		BaseFeePerGas         *big.Int            `json:"baseFeePerGas" gencodec:"required"`
		BlockMetadata         *BlockMetadata      `json:"blockMetadata" gencodec:"required"`
		L1Origin              *rawdb.L1Origin     `json:"l1Origin"      gencodec:"required"`
		VirtualBlock          bool                `json:"virtualBlock" gencodec:"required"`
	}
	var enc PayloadAttributes
	enc.Timestamp = hexutil.Uint64(p.Timestamp)
	enc.Random = p.Random
	enc.SuggestedFeeRecipient = p.SuggestedFeeRecipient
	enc.Withdrawals = p.Withdrawals
	enc.BeaconRoot = p.BeaconRoot
	enc.BaseFeePerGas = p.BaseFeePerGas
	enc.BlockMetadata = p.BlockMetadata
	enc.L1Origin = p.L1Origin
	enc.VirtualBlock = p.VirtualBlock
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (p *PayloadAttributes) UnmarshalJSON(input []byte) error {
	type PayloadAttributes struct {
		Timestamp             *hexutil.Uint64     `json:"timestamp"             gencodec:"required"`
		Random                *common.Hash        `json:"prevRandao"            gencodec:"required"`
		SuggestedFeeRecipient *common.Address     `json:"suggestedFeeRecipient" gencodec:"required"`
		Withdrawals           []*types.Withdrawal `json:"withdrawals"`
		BeaconRoot            *common.Hash        `json:"parentBeaconBlockRoot"`
		BaseFeePerGas         *big.Int            `json:"baseFeePerGas" gencodec:"required"`
		BlockMetadata         *BlockMetadata      `json:"blockMetadata" gencodec:"required"`
		L1Origin              *rawdb.L1Origin     `json:"l1Origin"      gencodec:"required"`
		VirtualBlock          *bool               `json:"virtualBlock" gencodec:"required"`
	}
	var dec PayloadAttributes
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Timestamp == nil {
		return errors.New("missing required field 'timestamp' for PayloadAttributes")
	}
	p.Timestamp = uint64(*dec.Timestamp)
	if dec.Random == nil {
		return errors.New("missing required field 'prevRandao' for PayloadAttributes")
	}
	p.Random = *dec.Random
	if dec.SuggestedFeeRecipient == nil {
		return errors.New("missing required field 'suggestedFeeRecipient' for PayloadAttributes")
	}
	p.SuggestedFeeRecipient = *dec.SuggestedFeeRecipient
	if dec.Withdrawals != nil {
		p.Withdrawals = dec.Withdrawals
	}
	if dec.BeaconRoot != nil {
		p.BeaconRoot = dec.BeaconRoot
	}
	if dec.BaseFeePerGas == nil {
		return errors.New("missing required field 'baseFeePerGas' for PayloadAttributes")
	}
	p.BaseFeePerGas = dec.BaseFeePerGas
	if dec.BlockMetadata == nil {
		return errors.New("missing required field 'blockMetadata' for PayloadAttributes")
	}
	p.BlockMetadata = dec.BlockMetadata
	if dec.L1Origin == nil {
		return errors.New("missing required field 'l1Origin' for PayloadAttributes")
	}
	p.L1Origin = dec.L1Origin
	if dec.VirtualBlock == nil {
		return errors.New("missing required field 'virtualBlock' for PayloadAttributes")
	}
	p.VirtualBlock = *dec.VirtualBlock
	return nil
}
