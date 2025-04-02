package edges

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/dharnitski/cc-hosts/access"
)

const (
	EdgesFolder         = "edges"
	EdgesReversedFolder = "edges_reversed"
	DefaultMaxSize      = 10_000
)

type Edge struct {
	// source vertice id
	fromID string
	// target vertice id
	toID string
}

func (v *Edge) FromID() string {
	return v.fromID
}

func (v *Edge) ToID() string {
	return v.toID
}

func LoadEdge(line string) (*Edge, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid line: %s, %d parts", line, len(parts))
	}
	return &Edge{fromID: parts[0], toID: parts[1]}, nil
}

type Edges struct {
	// offsets to find edges in edges files
	offsets Offsets
	getter  access.Getter
}

func NewEdges(getter access.Getter, offsets Offsets) *Edges {
	return &Edges{
		offsets: offsets,
		getter:  getter,
	}
}

// for source vertice id return list of target vertice ids
func (v *Edges) Get(ctx context.Context, fromID string) ([]string, error) {
	offsets := v.offsets.FindForFromID(fromID)

	type result struct {
		edges []string
		err   error
	}

	results := make(chan result, len(offsets))
	var wg sync.WaitGroup

	for file, offset := range offsets {
		wg.Add(1)
		go func(file string, offset TwoOffsets) {
			defer wg.Done()

			buffer, err := v.getter.Get(ctx, file, offset.From.offset, offset.To.offset-offset.From.offset)
			if err != nil {
				results <- result{nil, err}
				return
			}
			edges, err := findEdges(buffer, fromID)
			results <- result{edges, err}
		}(file, offset)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	allEdges := make([]string, 0)
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		allEdges = append(allEdges, res.edges...)
	}

	sort.Strings(allEdges)
	return allEdges, nil
}

func findEdges(buffer []byte, fromID string) ([]string, error) {
	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)
	results := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		vertice, err := LoadEdge(line)
		if err != nil {
			return nil, err
		}
		if vertice.fromID == fromID {
			results = append(results, vertice.toID)
			if len(results) >= DefaultMaxSize {
				break
			}
		} else {
			// items sorted and we can break after we reach items with different fromID
			if len(results) > 0 {
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	return results, nil
}
