package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/offsets"
	"github.com/dharnitski/cc-hosts/vertices"
)

const (
	dataFolder = "data"
)

var (
	edgesForwardFolder  = path.Join(dataFolder, edges.EdgesFolder)
	edgesReversedFolder = path.Join(dataFolder, edges.EdgesReversedFolder)
	verticesFolder      = path.Join(dataFolder, vertices.Folder)
)

func main() {
	err := createVerticesIndex()
	if err != nil {
		log.Fatal("Vertices Error: ", err)
	}

	err = createForwardEdgesIndex()
	if err != nil {
		log.Fatal("Edges Forward Error: ", err)
	}

	err = createBackwardEdgesIndex()
	if err != nil {
		log.Fatal("Edges Backward Error: ", err)
	}
}

func createVerticesIndex() error {
	fmt.Printf("Loading  Vertices from %s folder\n", verticesFolder)
	// entries are sorted by filename
	entries, err := os.ReadDir(verticesFolder)
	if err != nil {
		return fmt.Errorf("error reading directory %q: %w", verticesFolder, err)
	}

	results := vertices.Offsets{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(verticesFolder, entry.Name())
		fmt.Printf("Processing Vertices file: %s\n", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file %q: %w", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		items, err := processOneVerticesFile(scanner, entry.Name())
		if err != nil {
			return fmt.Errorf("error processing file %q: %w", filePath, err)
		}

		results.Append(items)
	}

	if results.Len() > 0 {
		err := results.Validate()
		if err != nil {
			return fmt.Errorf("error validating offsets: %w", err)
		}

		saveFile := fmt.Sprintf("%s/%s", offsets.Folder, offsets.VerticesOffsetsFile)

		err = results.Save(saveFile)
		if err != nil {
			return fmt.Errorf("error saving offsets: %w", err)
		}

		fmt.Printf("Saved %d Vertices offsets to %s\n", results.Len(), saveFile)
	}

	return nil
}

func processOneVerticesFile(scanner *bufio.Scanner, fileName string) ([]vertices.Offset, error) {
	result := make([]vertices.Offset, 0)
	// bytes offset in file
	offset := 0
	firstLine := true
	lastSavedOffset := 0
	domain := ""
	id := -1

	for scanner.Scan() {
		// read bytes to properly calculate offset
		bytes := scanner.Bytes()
		// +1 for newline. scanner returns the line without delimiter
		tokenLength := len(bytes) + 1

		line := string(bytes)
		vertice, err := vertices.LoadVertice(line)

		if err != nil {
			return nil, fmt.Errorf("invalid line: %q: %w", line, err)
		}

		domain = vertice.Domain()
		sid := vertice.ID()
		id, err = strconv.Atoi(sid)

		if err != nil {
			return nil, fmt.Errorf("invalid ID: %q: %w", sid, err)
		}

		if firstLine {
			firstLine = false

			result = append(result, vertices.NewOffset(offset, domain, id, fileName))
			lastSavedOffset = offset
		}

		if offset-lastSavedOffset >= vertices.FileChunkSize {
			result = append(result, vertices.NewOffset(offset, domain, id, fileName))
			lastSavedOffset = offset
		}

		offset += tokenLength
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("error reading file: %w", err)
	}
	// save the last offset
	result = append(result, vertices.NewOffset(offset, domain, id, fileName))

	return result, nil
}

func createForwardEdgesIndex() error {
	return createEdgesIndex(edgesForwardFolder, offsets.EdgesOffsetsFile)
}

func createBackwardEdgesIndex() error {
	return createEdgesIndex(edgesReversedFolder, offsets.EdgesReversedOffsetFile)
}

func createEdgesIndex(edgesFolder string, outFile string) error {
	fmt.Printf("Loading  Edges from %s folder\n", edgesFolder)
	// entries are sorted by filename
	entries, err := os.ReadDir(edgesFolder)
	if err != nil {
		return fmt.Errorf("error reading directory %q: %w", edgesFolder, err)
	}

	results := edges.Offsets{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(edgesFolder, entry.Name())
		fmt.Printf("Processing Edges file: %s\n", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file %q: %w", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		items, err := processOneEdgesFile(scanner, entry.Name())
		if err != nil {
			return fmt.Errorf("error processing file %q: %w", filePath, err)
		}

		results.Append(items)
	}

	if results.Len() > 0 {
		err := results.Validate()
		if err != nil {
			return fmt.Errorf("error validating offsets: %w", err)
		}

		saveFile := fmt.Sprintf("%s/%s", offsets.Folder, outFile)

		err = results.Save(saveFile)
		if err != nil {
			return fmt.Errorf("error saving offsets: %w", err)
		}

		fmt.Printf("Saved %d edges offsets to %s\n", results.Len(), saveFile)
	}

	return nil
}

func processOneEdgesFile(scanner *bufio.Scanner, fileName string) ([]edges.Offset, error) {
	result := make([]edges.Offset, 0)
	// bytes offset in file
	offset := 0
	firstLine := true
	lastSavedOffset := 0
	id := ""

	for scanner.Scan() {
		// read bytes to properly calculate offset
		bytes := scanner.Bytes()
		// +1 for newline. scanner returns the line without delimiter
		tokenLength := len(bytes) + 1

		line := string(bytes)
		edge, err := edges.LoadEdge(line)

		if err != nil {
			return nil, fmt.Errorf("invalid line: %q: %w", line, err)
		}

		id = edge.FromID()

		if firstLine {
			firstLine = false

			result = append(result, edges.NewOffset(offset, id, fileName))
			lastSavedOffset = offset
		}

		if offset-lastSavedOffset >= edges.FileChunkSize {
			result = append(result, edges.NewOffset(offset, id, fileName))
			lastSavedOffset = offset
		}

		offset += tokenLength
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("error reading file: %w", err)
	}
	// save the last offset
	result = append(result, edges.NewOffset(offset, id, fileName))

	return result, nil
}
