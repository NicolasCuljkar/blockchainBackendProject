package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Wallet struct {
	WalletAddress   string `json:"wallet_address"`
	CurrencyCode    string `json:"currency_code"`
	CurrencyBalance string `json:"currency_balance"`
}

type WalletResp struct {
	Blockchain string `json:"blockchain"`
	PinCode    string `json:"pin_code"`
}

func CreateWallet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	walletResp := WalletResp{}
	if err := json.Unmarshal(body, &walletResp); err != nil {
		log.Println("Failed to unmarshal payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if walletResp.Blockchain == "" || walletResp.PinCode == "" {
		fmt.Println("Empty values")
		return
	}

	newWallet := &Wallet{}
	newWallet.WalletAddress = "8613417vyg67"
	newWallet.CurrencyCode = "ETH"
	newWallet.CurrencyBalance = "1.08"

	// On écrit dans walletDB.json
	walletByte, err := json.Marshal(newWallet)
	//fmt.Println(string(walletByte))

	w.WriteHeader(http.StatusOK)
	w.Write(walletByte)

}

func main() {
	fmt.Println("Starting server on port 8001 ...")
	http.HandleFunc("/wallet", CreateWallet)
	log.Fatal(http.ListenAndServe(":8001", nil))
}
