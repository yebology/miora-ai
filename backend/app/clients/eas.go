// Package clients provides the EAS (Ethereum Attestation Service) client.
//
// This client interacts with the EAS and SchemaRegistry contracts on Base Sepolia
// to publish and query trading reputation attestations on-chain.
//
// Contract addresses (Base Sepolia predeploys):
//   - EAS:            0x4200000000000000000000000000000000000021
//   - SchemaRegistry: 0x4200000000000000000000000000000000000020
//
// Schema format: "uint8 score,string recommendation,uint32 totalTransactions,string chain"
package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"miora-ai/app/interfaces"
	"miora-ai/constants"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EASClient implements interfaces.IEASClient for Base Sepolia.
type EASClient struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	attester   common.Address
	easAddr    common.Address
	schemaUID  [32]byte
	chainID    *big.Int
	easABI     abi.ABI
}

// NewEASClient creates a new EAS client connected to Base Sepolia.
// schemaUID should be a hex string (with or without 0x prefix).
// If schemaUID is empty, attestation calls will fail — register schema first.
func NewEASClient(rpcURL, easContractAddr, schemaUID, attesterPrivateKey string) (*EASClient, error) {
	if rpcURL == "" || attesterPrivateKey == "" {
		return &EASClient{}, nil // Return empty client if not configured
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Base Sepolia RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(attesterPrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid attester private key: %w", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	attester := crypto.PubkeyToAddress(*publicKey)

	easAddr := common.HexToAddress(easContractAddr)

	var schemaUIDBytes [32]byte
	if schemaUID != "" {
		uidBytes := common.FromHex(schemaUID)
		copy(schemaUIDBytes[:], uidBytes)
	}

	parsedABI, err := abi.JSON(strings.NewReader(EASABIJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse EAS ABI: %w", err)
	}

	return &EASClient{
		client:     client,
		privateKey: privateKey,
		attester:   attester,
		easAddr:    easAddr,
		schemaUID:  schemaUIDBytes,
		chainID:    big.NewInt(constants.BaseSepoliaChainID),
		easABI:     parsedABI,
	}, nil
}

// Attest publishes a trading reputation attestation on-chain via EAS.
// The attestation data is ABI-encoded: (uint8 score, string recommendation, uint32 totalTransactions, string chain).
func (e *EASClient) Attest(recipient string, score uint8, recommendation string, totalTxns uint32, chain string) (string, string, error) {
	if e.client == nil {
		return "", "", fmt.Errorf("EAS client not configured — set EAS_RPC_URL and EAS_ATTESTER_PRIVATE_KEY")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ABI-encode the attestation data payload
	dataArgs := abi.Arguments{
		{Type: mustType("uint8")},
		{Type: mustType("string")},
		{Type: mustType("uint32")},
		{Type: mustType("string")},
	}
	encodedData, err := dataArgs.Pack(score, recommendation, totalTxns, chain)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode attestation data: %w", err)
	}

	recipientAddr := common.HexToAddress(recipient)

	// Build AttestationRequestData tuple
	requestData := struct {
		Recipient      common.Address
		ExpirationTime uint64
		Revocable      bool
		RefUID         [32]byte
		Data           []byte
		Value          *big.Int
	}{
		Recipient:      recipientAddr,
		ExpirationTime: 0,    // No expiration
		Revocable:      true, // Allow score updates
		RefUID:         [32]byte{},
		Data:           encodedData,
		Value:          big.NewInt(0),
	}

	// Pack the attest() call
	input, err := e.easABI.Pack("attest", struct {
		Schema [32]byte
		Data   struct {
			Recipient      common.Address
			ExpirationTime uint64
			Revocable      bool
			RefUID         [32]byte
			Data           []byte
			Value          *big.Int
		}
	}{
		Schema: e.schemaUID,
		Data:   requestData,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to pack attest call: %w", err)
	}

	// Get nonce and gas price
	nonce, err := e.client.PendingNonceAt(ctx, e.attester)
	if err != nil {
		return "", "", fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get gas price: %w", err)
	}

	// Estimate gas
	gasLimit, err := e.client.EstimateGas(ctx, ethereum.CallMsg{
		From:  e.attester,
		To:    &e.easAddr,
		Data:  input,
		Value: big.NewInt(0),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to estimate gas: %w", err)
	}

	// Build and sign transaction
	tx := types.NewTransaction(nonce, e.easAddr, big.NewInt(0), gasLimit, gasPrice, input)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(e.chainID), e.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	if err := e.client.SendTransaction(ctx, signedTx); err != nil {
		return "", "", fmt.Errorf("failed to send attestation tx: %w", err)
	}

	txHash := signedTx.Hash().Hex()

	// Wait for receipt to get the attestation UID from logs
	uid, err := e.waitForAttestationUID(ctx, signedTx.Hash())
	if err != nil {
		return "", txHash, fmt.Errorf("tx sent (%s) but failed to get attestation UID: %w", txHash, err)
	}

	return uid, txHash, nil
}

// waitForAttestationUID waits for the transaction receipt and extracts the attestation UID
// from the Attested event log.
func (e *EASClient) waitForAttestationUID(ctx context.Context, txHash common.Hash) (string, error) {
	attestedSig := crypto.Keccak256Hash([]byte("Attested(address,address,bytes32,bytes32)"))

	for i := range 60 { // Wait up to 60 seconds
		_ = i
		receipt, err := e.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if receipt.Status == 0 {
			return "", fmt.Errorf("attestation transaction reverted")
		}

		for _, log := range receipt.Logs {
			if log.Topics[0] == attestedSig {
				if len(log.Data) >= 32 {
					uid := common.BytesToHash(log.Data[:32])
					return uid.Hex(), nil
				}
			}
		}

		return "", fmt.Errorf("Attested event not found in receipt logs")
	}

	return "", fmt.Errorf("timeout waiting for transaction receipt")
}

// GetAttestation retrieves an attestation by UID from the EAS contract.
func (e *EASClient) GetAttestation(uid string) (*interfaces.AttestationData, error) {
	if e.client == nil {
		return nil, fmt.Errorf("EAS client not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uidBytes := common.FromHex(uid)
	var uidArr [32]byte
	copy(uidArr[:], uidBytes)

	input, err := e.easABI.Pack("getAttestation", uidArr)
	if err != nil {
		return nil, fmt.Errorf("failed to pack getAttestation call: %w", err)
	}

	result, err := e.client.CallContract(ctx, ethereum.CallMsg{
		To:   &e.easAddr,
		Data: input,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call getAttestation: %w", err)
	}

	output, err := e.easABI.Unpack("getAttestation", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack getAttestation result: %w", err)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("empty result from getAttestation")
	}

	attestStruct, ok := output[0].(struct {
		Uid            [32]byte       `json:"uid"`
		Schema         [32]byte       `json:"schema"`
		Time           uint64         `json:"time"`
		ExpirationTime uint64         `json:"expirationTime"`
		RevocationTime uint64         `json:"revocationTime"`
		RefUID         [32]byte       `json:"refUID"`
		Recipient      common.Address `json:"recipient"`
		Attester       common.Address `json:"attester"`
		Revocable      bool           `json:"revocable"`
		Data           []byte         `json:"data"`
	})
	if !ok {
		return nil, fmt.Errorf("unexpected attestation struct type")
	}

	// Decode the custom data payload
	dataArgs := abi.Arguments{
		{Type: mustType("uint8")},
		{Type: mustType("string")},
		{Type: mustType("uint32")},
		{Type: mustType("string")},
	}

	decoded, err := dataArgs.Unpack(attestStruct.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation data: %w", err)
	}

	return &interfaces.AttestationData{
		UID:               common.BytesToHash(attestStruct.Uid[:]).Hex(),
		Schema:            common.BytesToHash(attestStruct.Schema[:]).Hex(),
		Attester:          attestStruct.Attester.Hex(),
		Recipient:         attestStruct.Recipient.Hex(),
		Time:              attestStruct.Time,
		ExpirationTime:    attestStruct.ExpirationTime,
		RevocationTime:    attestStruct.RevocationTime,
		Revocable:         attestStruct.Revocable,
		Score:             decoded[0].(uint8),
		Recommendation:    decoded[1].(string),
		TotalTransactions: decoded[2].(uint32),
		Chain:             decoded[3].(string),
	}, nil
}

// mustType creates an ABI type or panics. Used for known-good type strings.
func mustType(t string) abi.Type {
	typ, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(fmt.Sprintf("invalid ABI type %q: %v", t, err))
	}
	return typ
}
