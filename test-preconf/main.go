package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
)

var (
	url     = "http://127.0.0.1:28545"
	chainID = big.NewInt(167001) // mainnet

	gasPriceMultiplier = big.NewInt(1_000)
	defaultValue       = big.NewInt(1_000_000_000_000) // in wei (1 eth = 1_000_000_000_000_000_000 wei)
	defaultGas         = uint64(100_000)
	defaultData        = make([]byte, 0) //

	god     = Account{privKey: "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"} // 0x8943545177806ED17B9F23F0a21ee5948eCaa776
	alice   = Account{privKey: "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d"} // 0xE25583099BA105D9ec0A67f5Ae86D90e50036425
	bob     = Account{privKey: "53321db7c1e331d93a11a41d16f004d7ff63972ec8ec7c25db329728ceeb1710"} // 0x614561D2d143621E126e87831AEF287678B442b8
	charlie = Account{privKey: "ab63b23eb7941c1251757e24b3d2350d2bc05c3c388d06f8fe6feafefb1e8c70"} // 0xf93Ee4Cf8c6c40b329b0c0626F28333c132CF241
	dave    = Account{privKey: "27515f805127bebad2fb9b183508bdacb8c763da16f54e0678b16e8f28ef3fff"} // 0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6
	eve     = Account{privKey: "7ff1a4c1d57e5e784d327c4c7651e952350bc271f156afb3d00d20f5ef924856"} // 0x2c57d1CFC6d5f8E4182a56b4cf75421472eBAEa4
	fred    = Account{privKey: "5d2344259f42259f82d2c140aa66102ba89b57b4883ee441a8b312622bd42491"} // 0x802dCbE1B1A97554B4F50DB5119E37E8e7336417
	george  = Account{privKey: "3a91003acaf4c21b3953d94fa4a6db694fa69e5242b2e37be05dd82761058899"} // 0x741bFE4802cE1C4b5b00F9Df2F5f179A1C89171A

	accounts = []Account{alice, bob}

	txsForAccount = func(addr common.Address, nonce uint64, currentSlot uint64, assignedSlots []uint64) []interface{} {
		acceptableSlotDeadline0 := calculateFirstAcceptableSlot(currentSlot, assignedSlots)
		acceptableSlotDeadline1 := new(big.Int).Add(acceptableSlotDeadline0, big.NewInt(1))
		acceptableSlotDeadline2 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(2))
		acceptableSlotDeadline3 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(3))
		acceptableSlotDeadline4 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(4))
		acceptableSlotDeadline5 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(5))
		acceptableSlotDeadline6 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(6))
		acceptableSlotDeadline7 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(7))
		acceptableSlotDeadline8 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(8))
		acceptableSlotDeadline9 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(9))
		acceptableSlotDeadline10 := new(big.Int).Add(acceptableSlotDeadline1, big.NewInt(10))
		pastSlotDeadline := new(big.Int).Add(acceptableSlotDeadline0, big.NewInt(-1))
		currentSlotDeadline := new(big.Int).SetUint64(currentSlot)
		offset := big.NewInt(int64(32 - common.SlotIndex(currentSlotDeadline.Uint64())))
		tooFarInFutureSlotDeadline := new(big.Int).Add(currentSlotDeadline, offset)

		return map[common.Address][]interface{}{
			*god.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: alice.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: bob.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: dave.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: eve.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: fred.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: acceptableSlotDeadline1},
			},
			*alice.Address(): {
				// types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: charlie.Address(), Deadline: acceptableSlotDeadline2},

				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: charlie.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 10, To: charlie.Address(), Deadline: acceptableSlotDeadline1},

				types.InclusionPreconfirmationTx{Nonce: nonce + 11, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 12, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 13, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 14, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 15, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 16, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 17, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 18, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 19, To: charlie.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 20, To: charlie.Address(), Deadline: acceptableSlotDeadline2},

				types.InclusionPreconfirmationTx{Nonce: nonce + 21, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 22, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 23, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 24, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 25, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 26, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 27, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 28, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 29, To: charlie.Address(), Deadline: acceptableSlotDeadline3},
				types.InclusionPreconfirmationTx{Nonce: nonce + 30, To: charlie.Address(), Deadline: acceptableSlotDeadline3},

				types.InclusionPreconfirmationTx{Nonce: nonce + 31, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 32, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 33, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 34, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 35, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 36, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 37, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 38, To: charlie.Address(), Deadline: acceptableSlotDeadline4},
				types.InclusionPreconfirmationTx{Nonce: nonce + 39, To: charlie.Address(), Deadline: acceptableSlotDeadline4},

				types.InclusionPreconfirmationTx{Nonce: nonce + 40, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 41, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 42, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 43, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 44, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 45, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 46, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 47, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 48, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 49, To: charlie.Address(), Deadline: acceptableSlotDeadline5},
				types.InclusionPreconfirmationTx{Nonce: nonce + 50, To: charlie.Address(), Deadline: acceptableSlotDeadline5},

				types.InclusionPreconfirmationTx{Nonce: nonce + 51, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 52, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 53, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 54, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 55, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 56, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 57, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 58, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 59, To: charlie.Address(), Deadline: acceptableSlotDeadline6},
				types.InclusionPreconfirmationTx{Nonce: nonce + 60, To: charlie.Address(), Deadline: acceptableSlotDeadline6},

				types.InclusionPreconfirmationTx{Nonce: nonce + 61, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 62, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 63, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 64, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 65, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 66, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 67, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 68, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 69, To: charlie.Address(), Deadline: acceptableSlotDeadline7},
				types.InclusionPreconfirmationTx{Nonce: nonce + 70, To: charlie.Address(), Deadline: acceptableSlotDeadline7},

				types.InclusionPreconfirmationTx{Nonce: nonce + 71, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 72, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 73, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 74, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 75, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 76, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 77, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 78, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 79, To: charlie.Address(), Deadline: acceptableSlotDeadline8},
				types.InclusionPreconfirmationTx{Nonce: nonce + 80, To: charlie.Address(), Deadline: acceptableSlotDeadline8},

				types.InclusionPreconfirmationTx{Nonce: nonce + 81, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 82, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 83, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 84, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 85, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 86, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 87, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 88, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 89, To: charlie.Address(), Deadline: acceptableSlotDeadline9},
				types.InclusionPreconfirmationTx{Nonce: nonce + 90, To: charlie.Address(), Deadline: acceptableSlotDeadline9},

				types.InclusionPreconfirmationTx{Nonce: nonce + 91, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 92, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 93, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 94, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 95, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 96, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 97, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 98, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 99, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
				types.InclusionPreconfirmationTx{Nonce: nonce + 100, To: charlie.Address(), Deadline: acceptableSlotDeadline10},
			},
			*bob.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 2, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 3, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 4, To: alice.Address()},
				types.LegacyTx{Nonce: nonce + 5, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 6, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 7, To: alice.Address()},
				types.AccessListTx{Nonce: nonce + 8, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 9, To: alice.Address()},
			},
			*charlie.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: george.Address(), Deadline: acceptableSlotDeadline2},
				types.DynamicFeeTx{Nonce: nonce + 6, To: george.Address()}, // does not have immediate receipt
			},
			*dave.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: pastSlotDeadline},           // rejected, past slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: currentSlotDeadline},        // rejected, current slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: tooFarInFutureSlotDeadline}, // not assigned slot
			},
			*eve.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 10, To: george.Address(), Deadline: acceptableSlotDeadline1}, // disallowed, preconf with future nonces
				types.InclusionPreconfirmationTx{Nonce: nonce + 11, To: george.Address(), Deadline: acceptableSlotDeadline1}, // disallowed, preconf with future nonces
			},
			*fred.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: acceptableSlotDeadline1},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: acceptableSlotDeadline1},
				types.DynamicFeeTx{Nonce: nonce + 2, To: george.Address()}, // processed, due to higher fees
				types.DynamicFeeTx{Nonce: nonce + 3, To: george.Address()}, // processed, due to higher fees
			},
			*george.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: alice.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: alice.Address()},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: alice.Address(), Deadline: acceptableSlotDeadline1}, // not preconfirmed
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: alice.Address(), Deadline: acceptableSlotDeadline1}, // not preconfirmed
			},
		}[addr]
	}

	txWithDefaults = func(i interface{}) *types.Transaction {
		var tx *types.Transaction

		switch txData := i.(type) {
		case types.LegacyTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasPrice == nil {
				txData.GasPrice = new(big.Int).Mul(big.NewInt(10), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.AccessListTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasPrice == nil {
				txData.GasPrice = new(big.Int).Mul(big.NewInt(7), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.DynamicFeeTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(big.Int).Mul(big.NewInt(8), gasPriceMultiplier)
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(big.Int).Mul(big.NewInt(5), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.BlobTx:
			if txData.Value == nil {
				txData.Value = uint256.NewInt(defaultValue.Uint64())
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(uint256.Int).Mul(uint256.NewInt(9), uint256.NewInt(gasPriceMultiplier.Uint64()))
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(uint256.Int).Mul(uint256.NewInt(7), uint256.NewInt(gasPriceMultiplier.Uint64()))
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		case types.InclusionPreconfirmationTx:
			if txData.Value == nil {
				txData.Value = defaultValue
			}
			txData.Gas = defaultGas
			if txData.GasFeeCap == nil {
				txData.GasFeeCap = new(big.Int).Mul(big.NewInt(1), gasPriceMultiplier)
			}
			if txData.GasTipCap == nil {
				txData.GasTipCap = new(big.Int).Mul(big.NewInt(0), gasPriceMultiplier)
			}
			txData.Data = defaultData
			tx = types.NewTx(&txData)
		}

		return tx
	}
)

type Account struct {
	privKey string
}

func (a *Account) PrivateKey() *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(a.privKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	return privateKey
}

func (a *Account) Address() *common.Address {
	publicKey := a.PrivateKey().Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	addr := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &addr
}

func calculateFirstAcceptableSlot(currentSlot uint64, assignedSlots []uint64) *big.Int {
	sort.Slice(assignedSlots, func(i, j int) bool {
		return assignedSlots[i] < assignedSlots[j]
	})

	fmt.Println("\nAssigned slots:", assignedSlots)
	fmt.Println("\nCurrent slot:", currentSlot)

	var firstAcceptableDeadline uint64
	for _, slot := range assignedSlots {
		if slot >= currentSlot+common.SlotsOffsetInAdvance {
			firstAcceptableDeadline = slot
			break
		}
	}
	if firstAcceptableDeadline == 0 {
		log.Fatal("Failed to find acceptable slot")
	}

	fmt.Println("\nAcceptable slot:", firstAcceptableDeadline)
	return new(big.Int).SetUint64(firstAcceptableDeadline)
}

func logTxReceiptInfo(client *ethclient.Client, signedTx *types.Transaction) {
	txReceipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			fmt.Printf("Tx Hash [%s]: Transaction receipt not found\n", signedTx.Hash())
			return
		} else {
			// log.Fatalf("Failed to get transaction receipt: %v", err)
			fmt.Printf("Failed to get transaction receipt: %v\n", err)
		}
	}
	fmt.Printf("Transaction receipt: TxHash [%s], Block Number: [%d], Status: [%d], Cumulative Gas Used: [%d], EffectiveGasPrice: [%d], GasUsed: [%d]\n", signedTx.Hash(), txReceipt.BlockNumber, txReceipt.Status, txReceipt.CumulativeGasUsed, txReceipt.EffectiveGasPrice, txReceipt.GasUsed)
}

func main() {
	var wg sync.WaitGroup

	// Connect to the Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var assignedSlots []uint64
	err = client.Client().CallContext(context.Background(), &assignedSlots, "taiko_fetchAssignedSlots")
	if err != nil {
		log.Fatalf("Failed to get assigned slots: %v", err)
	}

	var l1GenesisTimestamp uint64
	err = client.Client().CallContext(context.Background(), &l1GenesisTimestamp, "taiko_fetchL1GenesisTimestamp")
	if err != nil {
		log.Fatalf("Failed to get l1 genesis timestamp: %v", err)
	}

	// Iterate over each account
	for _, account := range accounts {
		// wg.Add(1)
		// go
		func(wg *sync.WaitGroup) {
			// defer wg.Done()

			nonce, err := client.PendingNonceAt(context.Background(), *account.Address())
			if err != nil {
				log.Fatalf("Failed to get nonce: %v", err)
			}

			currentSlot, _ := common.CurrentSlotAndEpoch(l1GenesisTimestamp, time.Now().Unix())

			// Get txs for the account
			txs := txsForAccount(*account.Address(), nonce, currentSlot, assignedSlots)
			fmt.Println("\nAccount:", *account.Address(), "Txs count:", len(txs))

			// Iterate over each tx
			for i := 0; i < len(txs); i++ {
				tx := txWithDefaults(txs[i])

				// Sign the tx
				signedTx, err := types.SignTx(tx, types.NewPreconfSigner(chainID), account.PrivateKey())
				if err != nil {
					// log.Fatalf("Failed to sign transaction: %v", err)
					fmt.Printf("Failed to sign transaction: %v\n", err)
				}

				// Send the tx
				err = client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					// log.Fatalf("Failed to send transaction: %v", err)
					fmt.Printf("Failed to send transaction: %v\n", err)
				}
				fmt.Printf("\nSubmitted Tx hash: %s slot deadline: %d\n", signedTx.Hash(), signedTx.Deadline().Uint64())

				// Get the tx
				_, _, err = client.TransactionByHash(context.Background(), signedTx.Hash())
				if err != nil {
					if errors.Is(err, ethereum.NotFound) {
						fmt.Printf("Tx Hash [%s]: Transaction not found\n", signedTx.Hash())
						return
					} else {
						// log.Fatalf("Failed to get transaction by hash %v, hash %s", err, signedTx.Hash())
						fmt.Printf("Failed to get transaction by hash %v, hash %s\n", err, signedTx.Hash())
					}
				}

				// Get the receipt
				logTxReceiptInfo(client, signedTx)
			}

		}(&wg)
		// wg.Wait()

		fmt.Println("\nDone")
	}
}
