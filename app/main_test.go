package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)

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
	req := httptest.NewRequest(http.MethodPost, "/fraud-score", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
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
