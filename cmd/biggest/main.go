package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		log.Printf("Error loading biggest hosts: %v\n", err)

		return
	}

	err = convertAndSave(ctx, biggestIDs, "biggest.json")
	if err != nil {
		log.Printf("Error converting and saving: %v\n", err)

		return
	}

	biggestIDs = make(map[string]int)
	edgesFolder = "data/edges_reversed"

	err = loadBiggestHosts(edgesFolder, biggestIDs)
	if err != nil {
		log.Printf("Error loading biggest hosts: %v\n", err)

		return
	}

	err = convertAndSave(ctx, biggestIDs, "biggest.reversed.json")
	if err != nil {
		log.Printf("Error converting and saving: %v\n", err)

		return
	}
}

func convertAndSave(ctx context.Context, biggestIDs map[string]int, outFile string) error {
	log.Printf("Getting Domains for IDs\n")

	offsets, err := vertices.NewOffsets()
	if err != nil {
		return fmt.Errorf("error loading offsets: %w", err)
	}

	vertices := vertices.NewVertices(file.NewGetter("data/vertices"), *offsets)

	biggest := make(map[string]int)

	for id, counter := range biggestIDs {
		vertice, err := vertices.GetByID(ctx, id)
		if err != nil {
			log.Printf("Error getting vertice by ID %s: %v\n", id, err)

			continue
		}

		if vertice == nil {
			log.Printf("Vertice %s not found\n", id)

			continue
		}

		domain := vertice.ReversedDomain()
		biggest[domain] = counter
	}

	file, err := os.Create(outFile) //nolint:gosec
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file %s: %v", outFile, err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	err = encoder.Encode(biggest)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	return nil
}

func loadBiggestHosts(edgesFolder string, biggest map[string]int) error {
	log.Printf("Loading  Edges from %s folder\n", edgesFolder)

	entries, err := os.ReadDir(edgesFolder)
	if err != nil {
		return fmt.Errorf("error reading directory %q: %w", edgesFolder, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(edgesFolder, entry.Name())
		log.Printf("Processing Edges file: %s\n", filePath)

		file, err := os.Open(filePath) //nolint:gosec
		if err != nil {
			return fmt.Errorf("error opening file %q: %w", filePath, err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("error closing file %s: %v", filePath, err)
			}
		}()

		scanner := bufio.NewScanner(file)

		err = processOneEdgesFile(scanner, biggest)
		if err != nil {
			return fmt.Errorf("error processing file %q: %w", filePath, err)
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
			return fmt.Errorf("invalid line: %s: %w", line, err)
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
				biggest[id] += counter
			}

			counter = 0
			id = edge.FromID()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}
