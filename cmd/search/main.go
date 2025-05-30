package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
)

func createSearcher(ctx context.Context) (*search.Searcher, error) {
	rootFolder := "./data"

	offsetsGetter := file.NewGetter(rootFolder)
	edgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesFolder))
	revEdgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesReversedFolder))
	verticesGetter := file.NewGetter(path.Join(rootFolder, vertices.Folder))

	return search.NewSearcher(ctx, offsetsGetter, edgesGetter, revEdgesGetter, verticesGetter)
}

func main() {
	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	args := os.Args
	if len(args) != 2 {
		return fmt.Errorf("expected 2 arguments, got %d", len(args))
	}

	target := args[1]

	// short timeout to load offsets
	ctxCreate, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	searcher, err := createSearcher(ctxCreate)
	if err != nil {
		return err
	}

	result, err := searcher.GetTargets(ctx, target)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData)) //nolint:forbidigo

	return nil
}
