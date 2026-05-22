package main

import (
	"github.com/coder/hnsw"
	"os"
)

const (
	GraphPath = "dataset/hnsw.graph"
)

// AddReferences add vectors to the graph struct for hierarchical search
func AddReferences(references []Reference) *hnsw.Graph[int] {
	g := hnsw.NewGraph[int]()
	for i, ref := range references {
		g.Add(hnsw.MakeNode(i+1, ref.Vector))
	}
	return g
}

// SaveGraph saves the graph structure into file system to later retrieval
func SaveGraph(g *hnsw.Graph[int]) error {
	path := GraphPath
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	err = g.Export(out)
	return err
}
