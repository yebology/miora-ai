// Command register-schema registers the Miora trading reputation schema on EAS (Base Sepolia).
//
// Usage:
//
//	go run cmd/register-schema/main.go
//
// Required env vars (from .env):
//
//	EAS_RPC_URL              — Base Sepolia RPC endpoint
//	EAS_SCHEMA_REGISTRY_ADDRESS — SchemaRegistry contract (default: 0x4200000000000000000000000000000000000020)
//	EAS_ATTESTER_PRIVATE_KEY — Private key for signing the registration tx
//
// Output:
//
//	Prints the schema UID — copy this to EAS_SCHEMA_UID in .env
package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"miora-ai/app/clients"
	"miora-ai/constants"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// Schema: "uint8 score,string recommendation,uint32 totalTransactions,string chain"
const schema = "uint8 score,string recommendation,uint32 totalTransactions,string chain"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env: ", err)
	}

	rpcURL := os.Getenv("EAS_RPC_URL")
	registryAddr := os.Getenv("EAS_SCHEMA_REGISTRY_ADDRESS")
	privateKeyHex := os.Getenv("EAS_ATTESTER_PRIVATE_KEY")

	if rpcURL == "" || privateKeyHex == "" {
		log.Fatal("Missing required env vars: EAS_RPC_URL, EAS_ATTESTER_PRIVATE_KEY")
	}
	if registryAddr == "" {
		registryAddr = "0x4200000000000000000000000000000000000020"
	}

	// Connect to Base Sepolia
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to RPC: ", err)
	}
	defer client.Close()

	// Parse private key
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyHex, "0x"))
	if err != nil {
		log.Fatal("Invalid private key: ", err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	sender := crypto.PubkeyToAddress(*publicKey)

	fmt.Printf("Registering EAS schema on Base Sepolia...\n")
	fmt.Printf("  Sender:   %s\n", sender.Hex())
	fmt.Printf("  Registry: %s\n", registryAddr)
	fmt.Printf("  Schema:   %s\n", schema)
	fmt.Printf("  Revocable: true\n\n")

	// Check balance
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, err := client.BalanceAt(ctx, sender, nil)
	if err != nil {
		log.Fatal("Failed to get balance: ", err)
	}
	fmt.Printf("  Balance:  %s wei\n\n", balance.String())

	if balance.Cmp(big.NewInt(0)) == 0 {
		log.Fatal("Wallet has no ETH. Get Base Sepolia ETH from a faucet first.\n  Faucet: https://www.coinbase.com/faucets/base-ethereum-goerli-faucet")
	}

	// Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(clients.SchemaRegistryABIJSON))
	if err != nil {
		log.Fatal("Failed to parse ABI: ", err)
	}

	// Pack register(schema, resolver=address(0), revocable=true)
	zeroAddr := common.Address{}
	input, err := parsedABI.Pack("register", schema, zeroAddr, true)
	if err != nil {
		log.Fatal("Failed to pack register call: ", err)
	}

	// Get nonce and gas price
	nonce, err := client.PendingNonceAt(ctx, sender)
	if err != nil {
		log.Fatal("Failed to get nonce: ", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("Failed to get gas price: ", err)
	}

	registry := common.HexToAddress(registryAddr)
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From: sender,
		To:   &registry,
		Data: input,
	})
	if err != nil {
		log.Fatal("Failed to estimate gas: ", err)
	}

	// Build and sign tx
	chainID := big.NewInt(constants.BaseSepoliaChainID)
	tx := types.NewTransaction(nonce, registry, big.NewInt(0), gasLimit, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal("Failed to sign tx: ", err)
	}

	// Send tx
	if err := client.SendTransaction(ctx, signedTx); err != nil {
		log.Fatal("Failed to send tx: ", err)
	}

	txHash := signedTx.Hash().Hex()
	fmt.Printf("Transaction sent: %s\n", txHash)
	fmt.Printf("Waiting for confirmation...\n\n")

	// Wait for receipt
	schemaUID := waitForSchemaUID(client, signedTx.Hash())

	fmt.Printf("✅ Schema registered successfully!\n\n")
	fmt.Printf("  Schema UID: %s\n", schemaUID)
	fmt.Printf("  Tx Hash:    %s\n", txHash)
	fmt.Printf("  Explorer:   https://base-sepolia.easscan.org/schema/view/%s\n\n", schemaUID)
	fmt.Printf("👉 Add this to your .env:\n")
	fmt.Printf("   EAS_SCHEMA_UID=%s\n", schemaUID)
}

func waitForSchemaUID(client *ethclient.Client, txHash common.Hash) string {
	// event Registered(bytes32 indexed uid, address indexed registerer, SchemaRecord schema)
	registeredSig := crypto.Keccak256Hash([]byte("Registered(bytes32,address,(bytes32,address,bool,string))"))

	for i := 0; i < 60; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		receipt, err := client.TransactionReceipt(ctx, txHash)
		cancel()

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if receipt.Status == 0 {
			log.Fatal("Transaction reverted!")
		}

		for _, logEntry := range receipt.Logs {
			if len(logEntry.Topics) > 0 && logEntry.Topics[0] == registeredSig {
				// Schema UID is the first indexed topic (Topics[1])
				if len(logEntry.Topics) > 1 {
					return logEntry.Topics[1].Hex()
				}
			}
		}

		log.Fatal("Registered event not found in receipt logs")
	}

	log.Fatal("Timeout waiting for transaction receipt")
	return ""
}
