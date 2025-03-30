package vertices

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/dharnitski/cc-hosts/access"
)

const (
	Concurrency    = 10
	Folder = "vertices"
)

type Vertice struct {
	// vertice id
	id string
	// domain name in reverse domain format
	// sample: com.example
	domain string
}

func (v *Vertice) ID() string {
	return v.id
}

// internal reversed domain format
// sample: com.example
func (v *Vertice) Domain() string {
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
func (v *Vertice) ReversedDomain() string {
	return ReverseDomain(v.domain)
}

func LoadVertice(line string) (*Vertice, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid line: %s, %d parts", line, len(parts))
	}
	return &Vertice{id: parts[0], domain: parts[1]}, nil
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

func (v *Vertices) GetByDomain(domain string) (*Vertice, error) {
	return v.get(domain, searchKeyDomain)
}

func (v *Vertices) GetByID(id string) (*Vertice, error) {
	return v.get(id, searchKeyID)
}

func (v *Vertices) GetByIDs(ids []string) ([]Vertice, error) {
	type result struct {
		vertice *Vertice
		err     error
		index   int
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

			vertice, err := v.GetByID(id)
			resultChan <- result{
				vertice: vertice,
				err:     err,
				index:   idx,
			}
		}(i, id)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Prepare results in order
	results := []Vertice{}
	errs := []error{}
	for res := range resultChan {
		if res.err != nil {
			errs = append(errs, res.err)
		}
		if res.vertice != nil {
			results = append(results, *res.vertice)
		}
	}
	if len(errs) > 0 {
		return results, fmt.Errorf("Errors: %v", errs)
	}
	return results, nil
}

func (v *Vertices) get(key string, searchSwitch searchKey) (*Vertice, error) {
	var from, to Offset
	switch searchSwitch {
	case searchKeyDomain:
		from, to = v.offsets.FindForDomain(key)
	case searchKeyID:
		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf("Invalid ID: %s", key)
		}
		from, to = v.offsets.FindForID(id)
	}
	// if we lucky and Vertice is in offset
	if from.domain == to.domain &&
		from.id == to.id && from.offset == to.offset {
		return &Vertice{id: strconv.Itoa(from.id), domain: from.domain}, nil
	}
	buffer, err := v.getter.Get(from.file, from.offset, to.offset-from.offset)
	if err != nil {
		return nil, err
	}

	return findVertice(buffer, key, searchSwitch)
}

func findVertice(buffer []byte, key string, searchSwitch searchKey) (*Vertice, error) {
	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		vertice, err := LoadVertice(line)
		if err != nil {
			return nil, err
		}
		switch searchSwitch {
		case searchKeyDomain:
			if vertice.domain == key {
				return vertice, nil
			}
		case searchKeyID:
			if vertice.id == key {
				return vertice, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	return nil, nil
}
