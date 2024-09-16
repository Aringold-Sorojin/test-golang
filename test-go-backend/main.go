package main

import (
	"log"
	"fmt"
	"encoding/json"
	"encoding/binary"
	"strconv"
	"unicode"
    "unicode/utf8"
    "bytes"
	"net/http"
	"github.com/edatts/go-massa"
)

type incrementByNRequest struct {
	Increment int64 `json:"increment"`
}

const (
	massaAPIURL = "https://buildnet.massa.net/api/v2:33035" // Use testnet URL
	senderSecretKey = "S1VWZwi3hLzmBqGdb7m2aMurRcZAtKpHPKLcr3gJRV6NueZQZSk"
	senderPassword = "password"
	contractAddress="AS1WMGJdaHeAhZpcj93QgvQASwsrrwXhw7ZYJbLYHd7jmU2FV75i"
)

var (
	massaClient = massa.NewClient()
	apiClient = massa.NewApiClient()
	senderAddr string
)

func main() {
	massaClient.Init(massaAPIURL)
	apiClient.Init(massaAPIURL)

	senderAddr, _ = massaClient.ImportFromPriv(senderSecretKey, senderPassword)

	// Load or generate keys
	http.HandleFunc("/readCounter", enableCORS(readCounterHandler))
	http.HandleFunc("/incrementByOne", enableCORS(incrementByOneHandler))
	http.HandleFunc("/incrementByN", enableCORS(incrementByNHandler))

	port := "8080"

	log.Printf("Server starting on port %s...\n", port)
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

func unicodeToASCIIString(input []byte) string {
    var output bytes.Buffer
    
    for len(input) > 0 {
        r, size := utf8.DecodeRune(input)
        if r == utf8.RuneError {
            input = input[1:]
            continue
        }
        
        if r <= 127 && unicode.IsPrint(r) {
            output.WriteByte(byte(r))
        }
        
        input = input[size:]
    }
    
    return output.String()
}

func readCounter() (string, error) {
	var params []byte

	callData := massa.CallData{
		Fee:            10_000_000,
		MaxGas:         10_000_000,
		Coins:          0,
		TargetAddress:  contractAddress,
		TargetFunction: "readCounter",
		Parameter:      params,
	}

	res, err := apiClient.ReadSC(senderAddr, callData)
	if err != nil {
		return "0", err
	}

	return unicodeToASCIIString(res), nil
}

func readCounterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := readCounter()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read counter: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json.RawMessage(result))
}

func incrementByOneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	counter, err := readCounter()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to increase counter by one: %v", err), http.StatusInternalServerError)
		return
	}

	current_counter, _ := strconv.ParseFloat(counter, 64)
	counter = strconv.FormatInt(int64(current_counter) + 1, 10)

	var params []byte

	opId, err := massaClient.CallSC(senderAddr, contractAddress, "incrementByOne", params, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Operation ID: %s", opId)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json.RawMessage(counter))
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

	counter, err := readCounter()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to increase counter by N: %v", err), http.StatusInternalServerError)
		return
	}

	current_counter, _ := strconv.ParseFloat(counter, 64)
	counter = strconv.FormatInt(int64(current_counter) + req.Increment, 10)

	var params []byte
	params = binary.LittleEndian.AppendUint32(params, uint32(req.Increment))

	opId, err := massaClient.CallSC(senderAddr, contractAddress, "incrementByN", params, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Operation ID: %s", opId)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json.RawMessage(counter))
}