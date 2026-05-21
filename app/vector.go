package main

import (
	"fmt"
	"github.com/coder/hnsw"
	"math"
	"time"
)

var normalization = map[string]float32{
	"max_amount":              10000,
	"max_installments":        12,
	"amount_vs_avg_ratio":     10,
	"max_minutes":             1440,
	"max_km":                  1000,
	"max_tx_count_24h":        20,
	"max_merchant_avg_amount": 10000,
}

var mccRisk = map[string]float32{
	"5411": 0.15,
	"5812": 0.30,
	"5912": 0.20,
	"5944": 0.45,
	"7801": 0.80,
	"7802": 0.75,
	"7995": 0.85,
	"4511": 0.35,
	"5311": 0.25,
	"5999": 0.50,
}

const (
	M              = 32
	efConstruction = 400
	efSearch       = 100
	K              = 5
	dims           = 14
)

// clamp keeps value between 0 and 1
func clamp(x float32) float32 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func dayOfWeek(t time.Time) float32 {
	fmt.Println(t.Weekday())
	if t.Weekday() == 0.0 {
		return 1
	}
	return (float32(t.Weekday()) - 1) / 6
}

func hourOfDay(t time.Time) float32 {
	return float32(t.Hour()) / 23.0
}

func minutesSinceLastTx(p Payload) float32 {
	if p.LastTransaction == nil {
		return -1
	}
	minutes := float32(time.Now().Sub(p.LastTransaction.Timestamp).Minutes())
	return clamp(minutes / normalization["max_minutes"])
}

func kmFromLastTx(p Payload) float32 {
	if p.LastTransaction == nil {
		return -1
	}
	return clamp(p.LastTransaction.KmFromCurrent / normalization["max_km"])
}

// Transform the payload into a 14-dimensional vector, following the normalization formulas.
func Transform(p Payload) []float32 {
	var isOnline float32 = 0
	if p.Terminal.IsOnline {
		isOnline = 1
	}
	var cardPresent float32 = 0
	if p.Terminal.CardPresent {
		cardPresent = 1
	}
	var unknownMerchant float32 = 1
	for _, mid := range p.Customer.KnownMerchants {
		if mid == p.Merchant.Id {
			unknownMerchant = 0
			break
		}
	}
	result := []float32{
		clamp(p.Transaction.Amount / normalization["max_amount"]),
		clamp(float32(p.Transaction.Installments) / normalization["max_installments"]),
		clamp((p.Transaction.Amount / p.Customer.AvgAmount) / normalization["amount_vs_avg_ratio"]),
		hourOfDay(p.Transaction.RequestedAt),
		dayOfWeek(p.Transaction.RequestedAt),
		minutesSinceLastTx(p),
		kmFromLastTx(p),
		clamp(p.Terminal.KmFromHome / normalization["max_km"]),
		clamp(float32(p.Customer.TxCount24h) / normalization["max_tx_count_24h"]),
		isOnline,
		cardPresent,
		unknownMerchant,
		mccRisk[p.Merchant.Mcc],
		clamp(p.Merchant.AvgAmount / normalization["max_merchant_avg_amount"]),
	}
	var ratio float32 = 10000
	for i, x := range result {
		result[i] = float32(math.Round(float64(x*ratio))) / ratio
	}
	return result
}

// AddReferences add vectors to the graph struct for hierarchical search
func AddReferences(references []Reference) *hnsw.Graph[int] {
	// var zero hnsw.Point = make([]float32, dims)
	// h := hnsw.New(M, efConstruction, zero)
	// h.Grow(len(references))
	// for i, ref := range references {
	// 	log.Printf("Add vector %v with index %d\n", ref.Vector, uint32(i+1))
	// 	h.Add(ref.Vector, uint32(i+1)) // ID must start from 1
	// }
	// return
	g := hnsw.NewGraph[int]()
	for i, ref := range references {
		g.Add(hnsw.MakeNode(i+1, ref.Vector))
	}
	return g
}

// SearchVector search a vector in the graph struct
func SearchVector(vector []float32, g *hnsw.Graph[int]) {
	if found := g.Search(vector, K); found != nil {
		fmt.Println("Found")
	}
}
