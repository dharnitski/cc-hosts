package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dharnitski/cc-hosts/vertices"
)

const (
	dataFolder     = "data"
	verticesFolder = dataFolder + "/vertices"
)

func main() {
	err := createVerticesIndex()
	if err != nil {
		log.Fatal("Vertices Error: ", err)
	}
}

func createVerticesIndex() error {
	fmt.Printf("Loading  Vertices from %s folder\n", verticesFolder)
	// entries are sorted by filename
	entries, err := os.ReadDir(verticesFolder)
	if err != nil {
		return fmt.Errorf("Error reading directory %q: %v", verticesFolder, err)
	}
	offsets := vertices.Offsets{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(verticesFolder, entry.Name())
		fmt.Printf("Processing Vertices file: %s\n", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file %q: %v", filePath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		items, err := processOneVerticesFile(scanner, entry.Name())
		if err != nil {
			return fmt.Errorf("Error processing file %q: %v", filePath, err)
		}
		offsets.Append(items)
	}
	if offsets.Len() > 0 {
		err := offsets.Validate()
		if err != nil {
			return fmt.Errorf("Error validating offsets: %v", err)
		}
		saveFile := fmt.Sprintf("%s/vertices.offsets.txt", dataFolder)
		err = offsets.Save(saveFile)
		if err != nil {
			return fmt.Errorf("Error saving offsets: %v", err)
		}
		fmt.Printf("Saved %d Vertices offsets to %s\n", offsets.Len(), saveFile)
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
	id := ""
	for scanner.Scan() {
		// read bytes to properly calculate offset
		bytes := scanner.Bytes()
		// +1 for newline. scanner returns the line without delimiter
		tokenLength := len(bytes) + 1

		line := string(bytes)
		vertice, err := vertices.LoadVertice(line)
		if err != nil {
			return nil, fmt.Errorf("Invalid line: %s\n", line)
		}
		domain = vertice.Domain()
		id = vertice.ID()
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
		return result, fmt.Errorf("Error reading file: %v\n", err)
	}
	// save the last offset
	result = append(result, vertices.NewOffset(offset, domain, id, fileName))
	return result, nil
}
