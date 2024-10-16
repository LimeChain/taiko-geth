package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func sendSingleTx() {
	url := "http://127.0.0.1:28545"

	chainID := big.NewInt(167001)

	alice := map[string]string{
		"privateKey": "39725efee3fb28614de3bacaffe4cc4bd8c436257e2c8bb887c4b5c4be45e76d",
		"from":       "0xE25583099BA105D9ec0A67f5Ae86D90e50036425",
	}

	bob := map[string]string{
		"privateKey": "53321db7c1e331d93a11a41d16f004d7ff63972ec8ec7c25db329728ceeb1710",
		"from":       "0x614561D2d143621E126e87831AEF287678B442b8",
	}

	// Connect to the Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Alice
	alicePrivateKey, err := crypto.HexToECDSA(alice["privateKey"])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	alicePublicKey := alicePrivateKey.Public()
	alicePublicKeyECDSA, ok := alicePublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	aliceAddress := crypto.PubkeyToAddress(*alicePublicKeyECDSA)

	// Bob
	bobPrivateKey, err := crypto.HexToECDSA(bob["privateKey"])
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	bobPublicKey := bobPrivateKey.Public()
	bobPublicKeyECDSA, ok := bobPublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	bobAddress := crypto.PubkeyToAddress(*bobPublicKeyECDSA)

	// Get the nonce
	nonce, err := client.PendingNonceAt(context.Background(), aliceAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	txFields := types.InclusionPreconfirmationTx{
		ChainID:    chainID,
		Nonce:      nonce,
		GasTipCap:  big.NewInt(10_000_000), // maxPriorityFeePerGas
		GasFeeCap:  big.NewInt(10_000_002), // maxFeePerGas
		Gas:        21_000,
		To:         &bobAddress,
		Value:      big.NewInt(1_000_000_000_000_000), // in wei (1 eth = 1_000_000_000_000_000_000)
		Data:       []byte{},
		AccessList: types.AccessList{},
		Deadline:   big.NewInt(123456),
		V:          big.NewInt(0),
		R:          big.NewInt(0),
		S:          big.NewInt(0),
	}
	fmt.Println(txFields)

	tx := types.NewTx(&txFields)

	buf := bytes.NewBuffer(nil)
	tx.EncodeRLP(buf)
	fmt.Printf("RLP: %x\n", buf.Bytes())

	fmt.Println("Tx hash:", tx.Hash().String())

	signedTx, err := types.SignTx(tx, types.NewPreconfSigner(chainID), alicePrivateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	v, r, s := signedTx.RawSignatureValues()
	fmt.Println("V:", v, "R:", r, "S:", s)

	b, _ := signedTx.MarshalBinary()
	fmt.Printf("Serialized Tx: %x\n", b)

	fmt.Println("Tx hash:", signedTx.Hash().String())

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Submitted Tx hash: %s\n", signedTx.Hash())
}
