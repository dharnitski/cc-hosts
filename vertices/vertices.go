package vertices

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/dharnitski/cc-hosts/access"
)

const (
	Concurrency = 100
	Folder      = "vertices"
)

type Vertex struct {
	// vertex id
	id string
	// domain name in reverse domain format
	// sample: com.example
	domain string
}

func (v *Vertex) ID() string {
	return v.id
}

// internal reversed domain format
// sample: com.example
func (v *Vertex) Domain() string {
	return v.domain
}

func ReverseDomain(domain string) string {
	parts := strings.Split(domain, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}

// ReversedDomain returns the domain as we use it in browser
// sample: example.com
func (v *Vertex) ReversedDomain() string {
	return ReverseDomain(v.domain)
}

func LoadVertex(line string) (*Vertex, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid line: %s, %d parts", line, len(parts))
	}

	return &Vertex{id: parts[0], domain: parts[1]}, nil
}

type Vertices struct {
	// offsets to find vertices in vertices files
	offsets Offsets
	getter  access.Getter
}

func NewVertices(getter access.Getter, offsets Offsets) *Vertices {
	return &Vertices{
		offsets: offsets,
		getter:  getter,
	}
}

type searchKey string

const (
	searchKeyDomain searchKey = "domain"
	searchKeyID     searchKey = "id"
)

func (v *Vertices) GetByDomain(ctx context.Context, domain string) (*Vertex, error) {
	return v.get(ctx, domain, searchKeyDomain)
}

func (v *Vertices) GetByID(ctx context.Context, id string) (*Vertex, error) {
	return v.get(ctx, id, searchKeyID)
}

func (v *Vertices) GetByIDs(ctx context.Context, ids []string) ([]Vertex, error) {
	type result struct {
		vertex *Vertex
		err    error
		index  int
	}

	resultChan := make(chan result, len(ids))

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, Concurrency)

	for i, id := range ids {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(idx int, id string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			vertex, err := v.GetByID(ctx, id)
			resultChan <- result{
				vertex: vertex,
				err:    err,
				index:  idx,
			}
		}(i, id)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(semaphore)
		close(resultChan)
	}()

	// Prepare results in order
	results := []Vertex{}
	errs := []error{}

	for res := range resultChan {
		if res.err != nil {
			errs = append(errs, res.err)
		}

		if res.vertex != nil {
			results = append(results, *res.vertex)
		}
	}

	if len(errs) > 0 {
		return results, fmt.Errorf("errors: %v", errs)
	}

	return results, nil
}

func (v *Vertices) get(ctx context.Context, key string, searchSwitch searchKey) (*Vertex, error) {
	var from, to Offset

	switch searchSwitch {
	case searchKeyDomain:
		from, to = v.offsets.FindForDomain(key)
	case searchKeyID:
		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("invalid ID: %s", key)
		}

		from, to = v.offsets.FindForID(id)
	}
	// if we lucky and Vertex is in offset
	if from.domain == to.domain &&
		from.id == to.id && from.offset == to.offset {
		return &Vertex{id: strconv.Itoa(from.id), domain: from.domain}, nil
	}

	buffer, err := v.getter.Get(ctx, from.file, from.offset, to.offset-from.offset)
	if err != nil {
		return nil, err
	}

	return findVertex(buffer, key, searchSwitch)
}

func findVertex(buffer []byte, key string, searchSwitch searchKey) (*Vertex, error) {
	reader := bytes.NewReader(buffer)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		vertex, err := LoadVertex(line)
		if err != nil {
			return nil, err
		}

		switch searchSwitch {
		case searchKeyDomain:
			if vertex.domain == key {
				return vertex, nil
			}
		case searchKeyID:
			if vertex.id == key {
				return vertex, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file for key %s: %w", key, err)
	}

	return nil, nil //nolint:nilnil
}
