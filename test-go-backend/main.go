package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	massaAPIURL = "https://test.massa.net/api/v2" // Use testnet URL
)

var (
	contractAddress string
	publicKey       ed25519.PublicKey
	privateKey      ed25519.PrivateKey
)

type RPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type RPCResponse struct {
	JsonRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type incrementByNRequest struct {
	Increment int `json:"increment"`
}

func main() {
	contractAddress = "CONTRACT_ADDRESS"

	// Load or generate keys
	var err error
	publicKey, privateKey, err = loadKeys()
	if err != nil {
		log.Fatalf("Failed to load keys: %v", err)
	}

	http.HandleFunc("/readCounter", enableCORS(readCounterHandler))
	http.HandleFunc("/incrementByOne", enableCORS(incrementByOneHandler))
	http.HandleFunc("/incrementByN", enableCORS(incrementByNHandler))

	port := "8080"

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}

func loadKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	// In a real-world scenario, load these securely from environment or secure storage
	privateKeyStr := "PRIVATE_KEY"
	if privateKeyStr == "" {
		return ed25519.GenerateKey(nil)
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, nil, err
	}
	return ed25519.GenerateKey(bytes.NewReader(privateKeyBytes))
}

func readCounterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := callSmartContract("readCounter", "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read counter: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func incrementByOneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := callSmartContract("incrementByOne", "")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to increase counter: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func incrementByNHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req incrementByNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	result, err := callSmartContract("incrementByN", strconv.Itoa(req.Increment))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to increase counter by N: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func callSmartContract(function string, parameter string) (json.RawMessage, error) {
	transaction := map[string]interface{}{
		"function":  function,
		"parameter": parameter,
		"address":   contractAddress,
		"fee":       "100000", // Adjust as needed
	}

	// Serialize the transaction
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	// Sign the transaction
	signature := ed25519.Sign(privateKey, transactionBytes)

	// Create the final transaction object
	finalTransaction := map[string]interface{}{
		"transaction":        transaction,
		"creator_public_key": base64.StdEncoding.EncodeToString(publicKey),
		"signature":          base64.StdEncoding.EncodeToString(signature),
	}

	// Send the transaction
	return sendRPCRequest("execute_smart_contract", []interface{}{finalTransaction})
}

func sendRPCRequest(method string, params []interface{}) (json.RawMessage, error) {
	rpcRequest := RPCRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	requestBody, err := json.Marshal(rpcRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(massaAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(body, &rpcResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}