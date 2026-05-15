package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Transaction struct {
	Amount float32 `json:"amount"`
	Installments int `json:"installments"`
	RequestedAt time.Time `json:"requested_at"`
}

type Customer struct {
	AvgAmount float32 `json:"avg_amount"`
	TxCount24h int `json:"tx_count_24h"`
	KnownMerchants []string `json:"known_merchants"`
}

type Merchant struct {
	Id string `json:"id"`
	Mcc string `json:"mcc"`
	AvgAmount float32 `json:"avg_amount"`
}

type Terminal struct {
	IsOnline bool `json:"is_online"`
	CardPresent bool `json:"card_present"`
	KmFromHome float32 `json:"km_from_home"`
}

type LastTransaction struct {
	Timestamp time.Time `json:"timestamp"`
	KmFromCurrent float32 `json:"km_from_current"`
}

type Payload struct {
	Id string `json:"id"`
	Transaction Transaction `json:"transaction"`
	Customer Customer `json:"customer"`
	Merchant Merchant `json:"merchant"`
	Terminal Terminal `json:"terminal"`
	LastTransaction LastTransaction `json:"last_transaction"`
}

type Response struct {
	Approved bool `json:"approved"`
	FraudScore float32 `json:"fraud_score"`
}

func fraudScoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var p Payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{}
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/fraude-score", fraudScoreHandler)
	http.ListenAndServe(":6969", nil)
}
