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
	ReferencesFilePath     = "dataset/references.json.gz"
	ReferencesDownloadLink = "https://github.com/zanfranceschi/rinha-de-backend-2026/raw/refs/heads/main/resources/references.json.gz"
)

type Reference struct {
	Vector []float32 `json:"vector"`
	Label  string    `json:"label"`
}

// ReadReferences load references dataset
func ReadReferences(file string) (error, []Reference) {
	f, err := os.Open(file)
	if err != nil {
		return err, nil
	}
	defer f.Close()
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return err, nil
	}
	defer gzReader.Close()
	data, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return err, nil
	}
	var references []Reference
	err = json.Unmarshal(data, &references)
	if err != nil {
		return err, nil
	}
	return nil, references
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
