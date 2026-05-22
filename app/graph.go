package main

import (
	"fmt"
	"github.com/coder/hnsw"
)

const (
	K = 5
)

// AddReferences add vectors to the graph struct for hierarchical search
func AddReferences(references []Reference) *hnsw.Graph[int] {
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
