package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/coder/hnsw"
)

func makeRequest(payload []byte) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/fraud-score", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestFraudScoreHandler(t *testing.T) {
	payload := []byte(`{
    "id": "tx-1329056812",
    "transaction": {
      "amount": 41.12,
      "installments": 2,
      "requested_at": "2026-03-11T18:45:53Z"
    },
    "customer": {
      "avg_amount": 82.24,
      "tx_count_24h": 3,
      "known_merchants": [
        "MERC-003",
        "MERC-016"
      ]
    },
    "merchant": {
      "id": "MERC-016",
      "mcc": "5411",
      "avg_amount": 60.25
    },
    "terminal": {
      "is_online": false,
      "card_present": true,
      "km_from_home": 29.2331036248
    },
    "last_transaction": null
  }`)
	req := makeRequest(payload)
	w := httptest.NewRecorder()
	err, references := ReadReferences("test/resources/example-references.json.gz")
	if err != nil {
		t.Error(err)
	}
	graph := AddReferences(references)
	FraudScoreHandler(w, req, graph, references)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}
	if !response.Approved {
		t.Errorf("Expected transaction approved but got %v", response.Approved)
	}
}

func TestExpectedFraud(t *testing.T) {
	payload := []byte(`
{
        "id": "tx-853640735",
        "transaction": {
          "amount": 5293.06,
          "installments": 8,
          "requested_at": "2028-09-19T03:34:29Z"
        },
        "customer": {
          "avg_amount": 60.14,
          "tx_count_24h": 11,
          "known_merchants": [
            "MERC-009",
            "MERC-001"
          ]
        },
        "merchant": {
          "id": "MERC-087",
          "mcc": "7995",
          "avg_amount": 21.57
        },
        "terminal": {
          "is_online": false,
          "card_present": false,
          "km_from_home": 265.7823290829
        },
        "last_transaction": {
          "timestamp": "2024-01-04T03:43:32Z",
          "km_from_current": 722.9372664641
        }
      }
`)
	req := makeRequest(payload)
	w := httptest.NewRecorder()
	err, references := ReadReferences(ReferencesFilePath)
	if err != nil {
		t.Error(err)
	}
	savedGraph, err := hnsw.LoadSavedGraph[int](GraphPath)
	if err != nil {
	 	t.Error(err)
	}
	FraudScoreHandler(w, req, savedGraph.Graph, references)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}
	if response.Approved {
		t.Errorf("Expected transaction not approved but got %v", response.Approved)
	}
}
