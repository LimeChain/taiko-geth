package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	url         = "http://127.0.0.1:28545"
	chainID     = big.NewInt(167001) // mainnet
	currentSlot = big.NewInt(52625)

	accounts = []map[string]string{
		{
			"privateKey":  "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31",
			"fromAddress": "0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"toAddress":   "0xE25583099BA105D9ec0A67f5Ae86D90e50036425",
		},
		{
			"privateKey":  "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d",
			"fromAddress": "0xE25583099BA105D9ec0A67f5Ae86D90e50036425",
			"toAddress":   "0x8943545177806ED17B9F23F0a21ee5948eCaa776",
		},
	}

	value = big.NewInt(1_000) // in wei (1 eth = 1_000_000_000_000_000_000)
	gas   = uint64(21_000)
	data  = make([]byte, 0)
)

func main() {
	// Connect to the Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Iterate over each account
	for _, account := range accounts {
		privateKey, err := crypto.HexToECDSA(account["privateKey"])
		if err != nil {
			log.Fatalf("Failed to load private key: %v", err)
		}
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("Error casting public key to ECDSA")
		}
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			log.Fatalf("Failed to get nonce: %v", err)
		}

		fmt.Println()
		fmt.Println("Account:", fromAddress)

		accountsTxs := map[string][]interface{}{
			"0x8943545177806ED17B9F23F0a21ee5948eCaa776": {
				types.InclusionPreconfirmationTx{Nonce: nonce, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 5, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.DynamicFeeTx{Nonce: nonce + 6}, // does not have immediate receipt
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(-1))},  // past deadline
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(1))},   // not assigned slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 9, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(100))}, // not assigned slot, in next L1 epoch
			},
			"0xE25583099BA105D9ec0A67f5Ae86D90e50036425": {
				types.InclusionPreconfirmationTx{Nonce: nonce, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 1, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.InclusionPreconfirmationTx{Nonce: nonce + 4, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
				types.DynamicFeeTx{Nonce: nonce + 5}, // does not have immediate receipt
				types.InclusionPreconfirmationTx{Nonce: nonce + 6, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(-1))},  // past deadline
				types.InclusionPreconfirmationTx{Nonce: nonce + 7, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(1))},   // not assigned slot
				types.InclusionPreconfirmationTx{Nonce: nonce + 8, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(100))}, // not assigned slot, in next L1 epoch
			},
		}

		accountTxs := accountsTxs[account["fromAddress"]]
		for i := 0; i < len(accountTxs); i++ {
			fmt.Println()

			// gasPrice, err := client.SuggestGasPrice(context.Background())
			// if err != nil {
			// 	log.Fatalf("Failed to suggest gas price: %v", err)
			// }
			toAddress := common.HexToAddress(account["toAddress"])

			var tx *types.Transaction
			switch txData := accountTxs[i].(type) {
			case types.InclusionPreconfirmationTx:
				txData.To = &toAddress
				txData.Value = value
				txData.Gas = gas
				txData.GasFeeCap = big.NewInt(6 * 1_000)
				txData.GasTipCap = big.NewInt(4 * 1_000)
				txData.Data = data
				tx = types.NewTx(&txData)
			case types.DynamicFeeTx:
				txData.To = &toAddress
				txData.Value = value
				txData.Gas = gas
				txData.GasFeeCap = big.NewInt(5 * 1_000)
				txData.GasTipCap = big.NewInt(3 * 1_000)
				txData.Data = data
				tx = types.NewTx(&txData)
			}

			signedTx, err := types.SignTx(tx, types.NewPreconfSigner(chainID), privateKey)
			if err != nil {
				log.Fatalf("Failed to sign transaction: %v", err)
			}

			err = client.SendTransaction(context.Background(), signedTx)
			if err != nil {
				log.Fatalf("Failed to send transaction: %v", err)
			}

			fmt.Printf("Submitted Tx hash: %s", signedTx.Hash())

			_, _, err = client.TransactionByHash(context.Background(), signedTx.Hash())
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					fmt.Printf("Tx Hash [%s]: Transaction not found", signedTx.Hash())
					continue
				} else {
					log.Fatalf("Failed to get transaction by hash %v, hash %s", err, signedTx.Hash())
				}
			}

			txReceipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					fmt.Printf("Tx Hash [%s]: Transaction receipt not found", signedTx.Hash())
					continue
				} else {
					log.Fatalf("Failed to get transaction receipt: %v", err)
				}
			}

			fmt.Printf("Transaction receipt: TxHash [%s], Block Number: [%d], Status: [%d], Cumulative Gas Used: [%d], EffectiveGasPrice: [%d], GasUsed: [%d]", signedTx.Hash(), txReceipt.BlockNumber, txReceipt.Status,
				txReceipt.CumulativeGasUsed, txReceipt.EffectiveGasPrice, txReceipt.GasUsed)
		}
	}
}
