package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	ReferencesFilePath         = "dataset/references.json.gz"
	ApprovedReferencesFilePath = "dataset/references.json"
	ReferencesDownloadLink     = "https://github.com/zanfranceschi/rinha-de-backend-2026/raw/refs/heads/main/resources/references.json.gz"
)

type Reference struct {
	Vector []float32 `json:"vector"`
	Label  string    `json:"label"`
}

type ApprovedReferences []bool

// ReadReferences load references dataset
func ReadReferences(file string) ([]Reference, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	data, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}
	var references []Reference
	err = json.Unmarshal(data, &references)
	if err != nil {
		return nil, err
	}
	return references, nil
}

// DownloadReferences downloads the references file
func DownloadReferences() error {
	err := os.Mkdir("dataset", 0755)
	if err != nil && !strings.Contains(err.Error(), "file exists") {
		return err
	}
	f, err := os.Create(ReferencesFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(ReferencesDownloadLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// SaveReferences save the array of bools representing that a reference at index n is approved
func SaveReferences(references []Reference) error {
	f, err := os.Create(ApprovedReferencesFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var approved ApprovedReferences = make([]bool, len(references) + 1)
	for n, ref := range references {
		approved[n+1] = ref.Label != "fraud" // index starts at 1
	}
	bytes, err := json.Marshal(approved)
	if err != nil {
		return err
	}
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// LoadReferences reads the boolean array of approved references (true - approved, false - fraud) from the file system
func LoadReferences() (ApprovedReferences, error) {
	f, err := os.Open(ApprovedReferencesFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var approved ApprovedReferences
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &approved)
	if err != nil {
		return nil, err
	}
	return approved, nil
}
