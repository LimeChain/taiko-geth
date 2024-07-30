package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	url = "http://127.0.0.1:28545"

	chainID = big.NewInt(167001) // mainnet

	currentSlot = big.NewInt(52625)

	accounts = []map[string]string{
		{
			"privateKey":  "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31",
			"fromAddress": "0x8943545177806ED17B9F23F0a21ee5948eCaa776",
			"toAddress":   "0x0000777735367b36bC9B61C50022d9D0700dB4Ec",
		},
		{ // golden touch account
			"privateKey":  "92954368afd3caa1f3ce3ead0069c1af414054aefe1ef9aeacc1bf426222ce38",
			"fromAddress": "0x0000777735367b36bC9B61C50022d9D0700dB4Ec",
			"toAddress":   "0x8943545177806ED17B9F23F0a21ee5948eCaa776",
		},
	}

	nonce = uint64(1718)

	accountsTxs = map[string][]interface{}{
		"0x8943545177806ED17B9F23F0a21ee5948eCaa776": {
			types.DynamicFeeTx{Nonce: nonce + 0},
			types.DynamicFeeTx{Nonce: nonce + 1},
			types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},

			types.DynamicFeeTx{Nonce: nonce + 2},
			types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(6))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 4, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(7))},
			types.DynamicFeeTx{Nonce: nonce + 3},
			types.DynamicFeeTx{Nonce: nonce + 4},
			types.InclusionPreconfirmationTx{Nonce: nonce + 5, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(8))},

			types.DynamicFeeTx{Nonce: nonce + 6},
			types.DynamicFeeTx{Nonce: nonce + 7},
		},
		"0x0000777735367b36bC9B61C50022d9D0700dB4Ec": {
			types.DynamicFeeTx{Nonce: nonce + 0},
			types.DynamicFeeTx{Nonce: nonce + 1},
			types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},

			types.DynamicFeeTx{Nonce: nonce + 2},
			types.InclusionPreconfirmationTx{Nonce: nonce + 2, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(5))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 3, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(6))},
			types.InclusionPreconfirmationTx{Nonce: nonce + 4, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(7))},
			types.DynamicFeeTx{Nonce: nonce + 3},
			types.DynamicFeeTx{Nonce: nonce + 4},
			types.InclusionPreconfirmationTx{Nonce: nonce + 5, Deadline: big.NewInt(0).Add(currentSlot, big.NewInt(8))},

			types.DynamicFeeTx{Nonce: nonce + 6},
			types.DynamicFeeTx{Nonce: nonce + 7},
		},
	}

	value    = big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit = uint64(200000)                  // in units
	data     []byte
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
		fmt.Println("Nonce:", nonce)

		// Send pre-configured number of txs
		accountTxs := accountsTxs[account["fromAddress"]]

		for i := 0; i < len(accountTxs); i++ {
			// nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
			// if err != nil {
			// 	log.Fatalf("Failed to get nonce: %v", err)
			// }
			toAddress := common.HexToAddress(account["toAddress"])
			gasPrice, err := client.SuggestGasPrice(context.Background())
			if err != nil {
				log.Fatalf("Failed to suggest gas price: %v", err)
			}

			var tx *types.Transaction
			switch txData := accountTxs[i].(type) {
			case types.InclusionPreconfirmationTx:
				// txData.Nonce = nonce
				txData.To = &toAddress
				txData.Value = value
				txData.Gas = gasLimit
				txData.GasPrice = gasPrice
				txData.Data = data
				tx = types.NewTx(&txData)
			case types.DynamicFeeTx:
				// txData.Nonce = nonce
				txData.To = &toAddress
				txData.Value = value
				txData.Gas = gasLimit
				txData.GasFeeCap = big.NewInt(5)
				txData.GasTipCap = big.NewInt(2)
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

			fmt.Printf("Tx hash: %s\n", signedTx.Hash().Hex())
		}

	}
}
