package common

import (
	"math/big"
)

func IncreaseByPercentage(percentage uint64, amount *big.Int) *big.Int {
	premium := new(big.Int).Div(
		new(big.Int).Mul(new(big.Int).SetUint64(percentage), amount),
		new(big.Int).SetUint64(100),
	)
	return new(big.Int).Add(amount, premium)
}
