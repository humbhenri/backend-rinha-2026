package main

import (
	"encoding/json"
	"os"
	"compress/gzip"
	"io/ioutil"
)

type Reference struct {
	Vector []float32 `json:"vector"`
	Label  string    `json:"label"`
}

func ReadReferences() (error, []Reference) {
	f, err := os.Open("dataset/references.json.gz")
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
