package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	referencesFilePath     = "dataset/references.json.gz"
	referencesDownloadLink = "https://github.com/zanfranceschi/rinha-de-backend-2026/raw/refs/heads/main/resources/references.json.gz"
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
	f, err := os.Create(referencesFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Println("References file output created")
	resp, err := http.Get(referencesDownloadLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Println("References data downloaded")
	_, err = io.Copy(f, resp.Body)
	return err
}
