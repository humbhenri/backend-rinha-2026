package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/coder/hnsw"
)

type RequestTest struct {
	Request            Payload `json:"request"`
	ExpectedApproved   bool    `json:"expected_approved"`
	ExpectedFraudScore float32 `json:"expected_fraud_score"`
}

const testInputPath = "test/resources/test-data.json"

var approvedReferences ApprovedReferences
var savedGraph *hnsw.SavedGraph[int]
var testInput []RequestTest

func makeRequest(payload []byte) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/fraud-score", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestMain(m *testing.M) {
	var err error
	approvedReferences, err = LoadReferences()
	if err != nil {
		panic(err)
	}
	savedGraph, err = hnsw.LoadSavedGraph[int](GraphPath)
	if err != nil {
		panic(err)
	}
	file, err := os.Open(testInputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(data), &testInput)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestRequests(t *testing.T) {
	for _, request := range testInput {
		t.Run(request.Request.Id, func(t *testing.T) {
			bytes, err := json.Marshal(request.Request)
			if err != nil {
				t.Error(err)
			}
			req := makeRequest(bytes)
			w := httptest.NewRecorder()
			FraudScoreHandler(w, req, savedGraph.Graph, approvedReferences)
			if w.Code != http.StatusOK {
				t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
			}
			var response Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Error(err)
			}
			if response.Approved != request.ExpectedApproved {
				t.Errorf("Expected approved = %v but got %v", request.ExpectedApproved, response.Approved)
			}
			if response.FraudScore != request.ExpectedFraudScore {
				t.Errorf("Expected fraud score %v but got %v", request.ExpectedFraudScore, response.FraudScore)
			}
		})
	}
}
