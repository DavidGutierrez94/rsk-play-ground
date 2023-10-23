package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"strconv"

	"github.com/davidgutierrez94/rsk-play-ground/internal/blockchain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}
type WalletAddressResponse struct {
	UserID string `json:"userID"`
	Wallet string `json:"wallet"`
	Index  uint32 `json:"index"`
}

type SendTransactionResponse struct {
	UserID string `json:"userID"`
	Index  uint32 `json:"index"`
	from   string `json:"from"`
	to     string `json:"to"`
	amount uint32 `json:"amount"`
	txId   string `json:"txId"`
	status string `json:"status"`
}

func getWalletAddress(w http.ResponseWriter, r *http.Request) {
	b := &blockchain.Blockchain{}
	userID := r.URL.Query().Get("userID")
	index, _ := strconv.ParseUint(r.URL.Query().Get("index"), 10, 32)
	address, err := b.CreateDepositAddress(userID, uint32(index))
	if err != nil {
		http.Error(w, "Failed to create deposit address", http.StatusInternalServerError)
		return
	}
	response := WalletAddressResponse{
		UserID: userID,
		Wallet: address,
		Index:  uint32(index),
	}

	// Marshal the response struct to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func getTransaction(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("txHash")
	w.Write([]byte(txHash))

}

func sendTransaction(w http.ResponseWriter, r *http.Request) {
	b := &blockchain.Blockchain{}
	userID := r.URL.Query().Get("userID")
	index, _ := strconv.ParseUint(r.URL.Query().Get("index"), 10, 32)
	recipientAddress := r.URL.Query().Get("recipientAddress")
	amount, _ := strconv.ParseUint(r.URL.Query().Get("amount"), 10, 32)

	RskNodeURL := "https://rsk.getblock.io/b9ef1fcf-8731-42e5-a7af-39c44b7bf055/testnet/"

	// Replace with your Ethereum private key
	privateKey, _ := b.GetPrivateKey(uint32(index))

	// Create an Ethereum client
	client, err := ethclient.Dial(RskNodeURL)
	if err != nil {
		http.Error(w, "Failed to connect to the Ethereum client", http.StatusInternalServerError)
		return
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		http.Error(w, "Failed to get nonce", http.StatusInternalServerError)
		return
	}

	value := big.NewInt(int64(amount)) // Amount to send in Wei (0 for a simple transaction)
	gasLimit := uint64(21000)          // Gas limit for a simple transaction
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		http.Error(w, "Failed to suggest gas price", http.StatusInternalServerError)
		return
	}
	/*
		if !common.IsHexAddress(recipientAddress) {
			http.Error(w, "Invalid recipient address", http.StatusBadRequest)
			return
		}*/

	toAddress := common.HexToAddress(recipientAddress)

	// Replace with the actual transaction data (if any)
	data := []byte{}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		http.Error(w, "Failed to get chain ID", http.StatusInternalServerError)
		return
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		http.Error(w, "Failed to sign transaction", http.StatusInternalServerError)
		return
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		http.Error(w, "Failed to send transaction", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Transaction sent:", signedTx.Hash().Hex())

	response := SendTransactionResponse{
		UserID: userID,
		Index:  uint32(index),
		from:   fromAddress.Hex(),
		to:     toAddress.Hex(),
		amount: uint32(amount),
		txId:   signedTx.Hash().Hex(),
		status: "pending",
	}

	// Marshal the response struct to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {

	r := mux.NewRouter()
	fmt.Println("Starting server on port 8000...")
	r.HandleFunc("/getWalletAddress", getWalletAddress).Methods("GET")
	r.HandleFunc("/getTransaction", getTransaction).Methods("GET")
	r.HandleFunc("/sendTransaction", sendTransaction).Methods("POST")
	http.ListenAndServe(":8000", r)
	fmt.Println("Server started on port 8000")

}
