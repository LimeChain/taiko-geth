package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
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
	url                = "http://127.0.0.1:28545"
	chainID            = big.NewInt(167001) // mainnet
	l1GenesisTimestamp = uint64(1726642054)

	currentSlot, _ = common.CurrentSlotAndEpoch(l1GenesisTimestamp, time.Now().Unix())

	pastSlotDeadline        = new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(-1))
	currentSlotDeadline     = new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(0))
	nextSlotDeadline        = new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(3))
	notAssignedSlotDeadline = new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(2))
	nextEpochDeadline       = new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(100))

	gasPriceMultiplier = big.NewInt(10_000)
	defaultValue       = big.NewInt(1_000_000_000) // in wei (1 eth = 1_000_000_000_000_000_000 wei)
	defaultGas         = uint64(21_000)
	defaultData        = make([]byte, 0) // 100_000

	god     = Account{privKey: "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"} // 0x8943545177806ED17B9F23F0a21ee5948eCaa776
	alice   = Account{privKey: "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d"} // 0xE25583099BA105D9ec0A67f5Ae86D90e50036425
	bob     = Account{privKey: "53321db7c1e331d93a11a41d16f004d7ff63972ec8ec7c25db329728ceeb1710"} // 0x614561D2d143621E126e87831AEF287678B442b8
	charlie = Account{privKey: "ab63b23eb7941c1251757e24b3d2350d2bc05c3c388d06f8fe6feafefb1e8c70"} // 0xf93Ee4Cf8c6c40b329b0c0626F28333c132CF241
	dave    = Account{privKey: "27515f805127bebad2fb9b183508bdacb8c763da16f54e0678b16e8f28ef3fff"} // 0xAe95d8DA9244C37CaC0a3e16BA966a8e852Bb6D6
	eve     = Account{privKey: "7ff1a4c1d57e5e784d327c4c7651e952350bc271f156afb3d00d20f5ef924856"} // 0x2c57d1CFC6d5f8E4182a56b4cf75421472eBAEa4
	fred    = Account{privKey: "5d2344259f42259f82d2c140aa66102ba89b57b4883ee441a8b312622bd42491"} // 0x802dCbE1B1A97554B4F50DB5119E37E8e7336417
	george  = Account{privKey: "3a91003acaf4c21b3953d94fa4a6db694fa69e5242b2e37be05dd82761058899"} // 0x741bFE4802cE1C4b5b00F9Df2F5f179A1C89171A

	accounts = []Account{alice}

	txsForAccount = func(addr common.Address, nonce uint64) []interface{} {
		return map[common.Address][]interface{}{
			*god.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: alice.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: bob.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: charlie.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: dave.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: eve.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: fred.Address(), Deadline: nextSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: nextSlotDeadline},
			},
			*alice.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: george.Address()},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.DynamicFeeTx{Nonce: nonce + 5, To: alice.Address()},                                                  // does not have immediate receipt
				// types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: alice.Address(), Deadline: pastSlotDeadline},        // past deadline
				// types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: alice.Address(), Deadline: notAssignedSlotDeadline}, // not assigned slot
				// types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: alice.Address(), Deadline: nextEpochDeadline},       // not assigned slot, in next L1 epoch
				// types.DynamicFeeTx{Nonce: nonce + 0, To: alice.Address()},
				// types.DynamicFeeTx{Nonce: nonce + 1, To: alice.Address()},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: alice.Address(), Deadline: nextSlotDeadline},
				// types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: alice.Address(), Deadline: nextSlotDeadline},
			},
			*bob.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: george.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: george.Address()},
				types.DynamicFeeTx{Nonce: nonce + 2, To: george.Address()},
				types.LegacyTx{Nonce: nonce + 3, To: george.Address()},
				types.LegacyTx{Nonce: nonce + 4, To: george.Address()},
				types.AccessListTx{Nonce: nonce + 5, To: george.Address()},
			},
			*charlie.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, To: george.Address(), Deadline: currentSlotDeadline},
				types.DynamicFeeTx{Nonce: nonce + 6, To: george.Address()},                                                  // does not have immediate receipt
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: george.Address(), Deadline: pastSlotDeadline},        // past deadline
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: george.Address(), Deadline: notAssignedSlotDeadline}, // not assigned slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, To: george.Address(), Deadline: nextEpochDeadline},       // not assigned slot, in next L1 epoch
			},
			*dave.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: nextSlotDeadline}, // processed with higher priority
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: nextSlotDeadline}, // processed with higher priority
				types.DynamicFeeTx{Nonce: nonce + 0, To: george.Address()},
				types.DynamicFeeTx{Nonce: nonce + 1, To: george.Address()},
			},
			*eve.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, To: george.Address(), Deadline: currentSlotDeadline},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, To: george.Address(), Deadline: currentSlotDeadline},
				types.DynamicFeeTx{Nonce: nonce + 5, To: george.Address()},                                                  // does not have immediate receipt
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, To: george.Address(), Deadline: pastSlotDeadline},        // past deadline
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, To: george.Address(), Deadline: notAssignedSlotDeadline}, // not assigned slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, To: george.Address(), Deadline: nextEpochDeadline},       // not assigned slot, in next L1 epoch
			},
			*fred.Address(): {
				types.DynamicFeeTx{Nonce: nonce + 0, To: george.Address()},                                              // does not have immediate receipt
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, To: george.Address(), Deadline: currentSlotDeadline}, // not preconfirmed
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, To: george.Address(), Deadline: currentSlotDeadline}, // not preconfirmed
			},
			*george.Address(): {
				types.InclusionPreconfirmationTx{Nonce: nonce + 0, To: george.Address(), Deadline: new(big.Int).Add(new(big.Int).SetUint64(currentSlot), big.NewInt(3))},
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

func main() {
	var wg sync.WaitGroup

	// Connect to the Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Iterate over each account
	for _, account := range accounts {
		nonce, err := client.PendingNonceAt(context.Background(), *account.Address())
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}

		// Get txs for the account
		txs := txsForAccount(*account.Address(), nonce)
		fmt.Println("\nAccount:", *account.Address(), "Txs count:", len(txs))

		// Iterate over each tx
		for i := 0; i < len(txs); i++ {
			// wg.Add(1)
			// go
			func(wg *sync.WaitGroup, tx *types.Transaction) {
				// defer wg.Done()

				// Sign the tx
				signedTx, err := types.SignTx(tx, types.NewPreconfSigner(chainID), account.PrivateKey())
				if err != nil {
					log.Fatalf("Failed to sign transaction: %v", err)
				}

				// Send the tx
				err = client.SendTransaction(context.Background(), signedTx)
				if err != nil {
					log.Fatalf("Failed to send transaction: %v", err)
				}
				fmt.Printf("Submitted Tx hash: %s slot deadline: %d\n", signedTx.Hash(), signedTx.Deadline().Uint64())

				// Get the tx
				_, _, err = client.TransactionByHash(context.Background(), signedTx.Hash())
				if err != nil {
					if errors.Is(err, ethereum.NotFound) {
						fmt.Printf("Tx Hash [%s]: Transaction not found\n", signedTx.Hash())
						return
					} else {
						log.Fatalf("Failed to get transaction by hash %v, hash %s", err, signedTx.Hash())
					}
				}

				// Get the receipt
				txReceipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
				if err != nil {
					if errors.Is(err, ethereum.NotFound) {
						fmt.Printf("Tx Hash [%s]: Transaction receipt not found\n", signedTx.Hash())
						return
					} else {
						log.Fatalf("Failed to get transaction receipt: %v", err)
					}
				}

				fmt.Printf("Transaction receipt: TxHash [%s], Block Number: [%d], Status: [%d], Cumulative Gas Used: [%d], EffectiveGasPrice: [%d], GasUsed: [%d]\n", signedTx.Hash(), txReceipt.BlockNumber, txReceipt.Status, txReceipt.CumulativeGasUsed, txReceipt.EffectiveGasPrice, txReceipt.GasUsed)
			}(&wg, txWithDefaults(txs[i]))
		}

		// wg.Wait()
		fmt.Println("Done")
	}
}
