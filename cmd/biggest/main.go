package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/vertices"
)

const (
	savingTrashload = 1_0000
)

func main() {
	ctx := context.Background()
	biggestIDs := make(map[string]int)
	edgesFolder := "data/edges"
	err := loadBiggestHosts(edgesFolder, biggestIDs)
	if err != nil {
		fmt.Printf("Error loading biggest hosts: %v\n", err)
		return
	}

	err = convertAndSave(ctx, biggestIDs, "biggest.json")
	if err != nil {
		fmt.Printf("Error converting and saving: %v\n", err)
		return
	}

	biggestIDs = make(map[string]int)
	edgesFolder = "data/edges_reversed"
	err = loadBiggestHosts(edgesFolder, biggestIDs)
	if err != nil {
		fmt.Printf("Error loading biggest hosts: %v\n", err)
		return
	}

	err = convertAndSave(ctx, biggestIDs, "biggest.reversed.json")
	if err != nil {
		fmt.Printf("Error converting and saving: %v\n", err)
		return
	}
}

func convertAndSave(ctx context.Context, biggestIDs map[string]int, outFile string) error {
	fmt.Printf("Getting Domains for IDs\n")

	offsets := vertices.Offsets{}
	err := offsets.Load(fmt.Sprintf("data/%s", access.VerticesOffsetsFile))
	if err != nil {
		return fmt.Errorf("Error loading offsets: %v", err)
	}
	vertices := vertices.NewVertices(file.NewGetter("data/vertices"), offsets)

	biggest := make(map[string]int)
	for id, counter := range biggestIDs {
		vertice, err := vertices.GetByID(ctx, id)
		if err != nil {
			fmt.Printf("Error getting vertice by ID %s: %v\n", id, err)
			continue
		}
		if vertice == nil {
			fmt.Printf("Vertice %s not found\n", id)
			continue
		}
		domain := vertice.ReversedDomain()
		biggest[domain] = counter
	}

	file, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(biggest)
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}
	return nil
}

func loadBiggestHosts(edgesFolder string, biggest map[string]int) error {
	fmt.Printf("Loading  Edges from %s folder\n", edgesFolder)
	entries, err := os.ReadDir(edgesFolder)
	if err != nil {
		return fmt.Errorf("Error reading directory %q: %v", edgesFolder, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(edgesFolder, entry.Name())
		fmt.Printf("Processing Edges file: %s\n", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file %q: %v", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		err = processOneEdgesFile(scanner, biggest)
		if err != nil {
			return fmt.Errorf("Error processing file %q: %v", filePath, err)
		}
	}
	return nil
}

func processOneEdgesFile(scanner *bufio.Scanner, biggest map[string]int) error {
	id := ""
	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		edge, err := edges.LoadEdge(line)
		if err != nil {
			return fmt.Errorf("Invalid line: %s\n", line)
		}
		// first line
		if id == "" {
			id = edge.FromID()
			continue
		}
		// continuation of the same Vertice
		if id == edge.FromID() {
			counter++
			continue
		}
		// data for next Vertice
		if edge.FromID() != id {
			if counter > savingTrashload {
				biggest[id] = biggest[id] + counter
			}
			counter = 0
			id = edge.FromID()
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: %v\n", err)
	}
	return nil
}
