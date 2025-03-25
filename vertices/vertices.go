package vertices

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/file"
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

func (v *Vertice) Domain() string {
	return v.domain
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
	// folder with vertices files
	folder string
	getter access.Getter
}

func NewVertices(folder string, offsets Offsets) *Vertices {
	return &Vertices{
		folder:  folder,
		offsets: offsets,
		getter:  file.NewGetter(folder),
	}
}

type verticeMatch func(Vertice, string) bool

func matchByDomain(v Vertice, domain string) bool {
	return v.domain == domain
}

func matchByID(v Vertice, id string) bool {
	return v.id == id
}

func (v *Vertices) GetByDomain(domain string) (*Vertice, error) {
	return v.get(domain, matchByDomain)
}

func (v *Vertices) GetByID(id string) (*Vertice, error) {
	return v.get(id, matchByID)
}

func (v *Vertices) get(key string, matchFn verticeMatch) (*Vertice, error) {
	from, to := v.offsets.FindForDomain(key)
	// if we lucky and Vertice is in offset
	if from.domain == to.domain &&
		from.id == to.id && from.offset == to.offset {
		return &Vertice{id: strconv.Itoa(from.id), domain: from.domain}, nil
	}
	buffer, err := v.getter.Get(from.file, from.offset, to.offset-from.offset)
	if err != nil {
		return nil, err
	}

	return findVertice(buffer, key, matchFn)
}

func findVertice(buffer []byte, key string, matchFn verticeMatch) (*Vertice, error) {
	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		vertice, err := LoadVertice(line)
		if err != nil {
			return nil, err
		}

		match := matchFn(*vertice, key)
		if match {
			return vertice, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	return nil, nil
}
