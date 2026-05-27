package main

import (
	"encoding/json"
	"github.com/coder/hnsw"
	"log"
	"net/http"
	"time"
	"flag"
)

var setupFlag = flag.Bool("s", false, "download and preprocess references file and exit")

type Transaction struct {
	Amount       float32   `json:"amount"`
	Installments int       `json:"installments"`
	RequestedAt  time.Time `json:"requested_at"`
}

type Customer struct {
	AvgAmount      float32  `json:"avg_amount"`
	TxCount24h     int      `json:"tx_count_24h"`
	KnownMerchants []string `json:"known_merchants"`
}

type Merchant struct {
	Id        string  `json:"id"`
	Mcc       string  `json:"mcc"`
	AvgAmount float32 `json:"avg_amount"`
}

type Terminal struct {
	IsOnline    bool    `json:"is_online"`
	CardPresent bool    `json:"card_present"`
	KmFromHome  float32 `json:"km_from_home"`
}

type LastTransaction struct {
	Timestamp     time.Time `json:"timestamp"`
	KmFromCurrent float32   `json:"km_from_current"`
}

type Payload struct {
	Id              string           `json:"id"`
	Transaction     *Transaction     `json:"transaction"`
	Customer        *Customer        `json:"customer"`
	Merchant        *Merchant        `json:"merchant"`
	Terminal        *Terminal        `json:"terminal"`
	LastTransaction *LastTransaction `json:"last_transaction"`
}

type Response struct {
	Approved   bool    `json:"approved"`
	FraudScore float32 `json:"fraud_score"`
}

func calcFraudScore(neighbors []hnsw.Node[int], approvedReferences ApprovedReferences) float32 {
	log.Printf("Neighbors count: %d\n", len(neighbors))
	frauds := 0
	for _, neighbor := range neighbors {
		index := neighbor.Key
		approved := approvedReferences[index]
		log.Printf("reference index %d, label %v\n", index, approved)
		if !approved {
			frauds++
		}
	}
	return float32(frauds) / 5.0
}

func FraudScoreHandler(w http.ResponseWriter, r *http.Request, graph *hnsw.Graph[int], approvedReferences ApprovedReferences) {
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
	vector := Transform(p)
	log.Printf("Payload transformed into vector %v\n", vector)
	neighbors := graph.Search(vector, 5)
	fraudScore := calcFraudScore(neighbors, approvedReferences)
	approved := fraudScore < 0.6
	log.Printf("Transaction approved: %v. Fraud score: %v", approved, fraudScore)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := Response{Approved: approved, FraudScore: fraudScore}
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	if *setupFlag {
		log.Println("Executing setup")
		err := DownloadReferences()
		if err != nil {
			panic(err)
		}
		log.Println("References data downloaded")
		references, err := ReadReferences(ReferencesFilePath)
		if err != nil {
			panic(err)
		}
		log.Println("References file read")
		err = SaveReferences(references)
		if err != nil {
			panic(err)
		}
		g := AddReferences(references)
		log.Println("Graph struct done")
		err = SaveGraph(g)
		if err != nil {
			panic(err)
		}
		log.Println("Graph struct saved to file system, setup completed")
		return
	}
 	// log.Println("Reading references file")
	// references, err := ReadReferences(ReferencesFilePath)
	// if err != nil {
	// 	panic(err)
	// }
	log.Println("Read approved references file")
	approvedReferences, err := LoadReferences()
	if err != nil {
		panic(err)
	}
	savedGraph, err := hnsw.LoadSavedGraph[int](GraphPath)
	if err != nil {
	 	panic(err)
	}
	log.Println("References added to the graph struct")
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/fraud-score", func(w http.ResponseWriter, r *http.Request) {
		FraudScoreHandler(w, r, savedGraph.Graph, approvedReferences)
	})
	log.Println("App running on port 6969")
	http.ListenAndServe(":6969", nil)
}
