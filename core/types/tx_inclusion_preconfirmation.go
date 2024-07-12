// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// InclusionPreconfirmationTx is the data of EIP-2930 access list transactions.
type InclusionPreconfirmationTx struct {
	ChainID    *big.Int        // destination chain ID
	Nonce      uint64          // nonce of sender account
	GasPrice   *big.Int        // wei per gas
	Gas        uint64          // gas limit
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int        // wei amount
	Data       []byte          // contract invocation input data
	AccessList AccessList      // EIP-2930 access list
	Deadline   *big.Int        // check if the deadline has passed before committing to it
	V, R, S    *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *InclusionPreconfirmationTx) copy() TxData {
	cpy := &InclusionPreconfirmationTx{
		Nonce:    tx.Nonce,
		Deadline: new(big.Int),
		To:       copyAddressPtr(tx.To),
		Data:     common.CopyBytes(tx.Data),
		Gas:      tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasPrice:   new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	if tx.Deadline != nil {
		cpy.Deadline.Set(tx.Deadline)
	}
	return cpy
}

// accessors for innerTx.
func (tx *InclusionPreconfirmationTx) txType() byte           { return InclusionPreconfirmationTxType }
func (tx *InclusionPreconfirmationTx) chainID() *big.Int      { return tx.ChainID }
func (tx *InclusionPreconfirmationTx) accessList() AccessList { return tx.AccessList }
func (tx *InclusionPreconfirmationTx) data() []byte           { return tx.Data }
func (tx *InclusionPreconfirmationTx) gas() uint64            { return tx.Gas }
func (tx *InclusionPreconfirmationTx) gasPrice() *big.Int     { return tx.GasPrice }
func (tx *InclusionPreconfirmationTx) gasTipCap() *big.Int    { return tx.GasPrice }
func (tx *InclusionPreconfirmationTx) gasFeeCap() *big.Int    { return tx.GasPrice }
func (tx *InclusionPreconfirmationTx) value() *big.Int        { return tx.Value }
func (tx *InclusionPreconfirmationTx) nonce() uint64          { return tx.Nonce }
func (tx *InclusionPreconfirmationTx) to() *common.Address    { return tx.To }
func (tx *InclusionPreconfirmationTx) deadline() *big.Int     { return tx.Deadline }

func (tx *InclusionPreconfirmationTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	fmt.Println("InclusionPreconfirmationTx effectiveGasPrice")
	if baseFee == nil {
		return dst.Set(tx.GasPrice)
	}
	preconfirmationBaseFeePerGas := dst.Mul(new(uint256.Int).SetUint64(params.InclusionPreconfirmationFeePremium).ToBig(), baseFee)
	tip := preconfirmationBaseFeePerGas.Div(preconfirmationBaseFeePerGas, new(uint256.Int).SetUint64(100).ToBig())
	return tip.Add(tip, baseFee)
}

func (tx *InclusionPreconfirmationTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *InclusionPreconfirmationTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *InclusionPreconfirmationTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *InclusionPreconfirmationTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *DynamicFeeTx) deadline() *big.Int {
	return big.NewInt(0)
}

func (tx *LegacyTx) deadline() *big.Int {
	return big.NewInt(0)
}

func (tx *AccessListTx) deadline() *big.Int {
	return big.NewInt(0)
}

func (tx *BlobTx) deadline() *big.Int {
	return big.NewInt(0)
}
