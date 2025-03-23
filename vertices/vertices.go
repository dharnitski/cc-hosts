package vertices

import (
	"bufio"
	"bytes"
	"fmt"
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

func NewVertices(folder string, offsets Offsets) Vertices {
	return Vertices{
		folder:  folder,
		offsets: offsets,
		getter:  file.NewGetter(folder),
	}
}

func (v *Vertices) Get(domain string) (*Vertice, error) {
	from, to := v.offsets.FindForDomain(domain)

	// if we lucky and Vertice is in offset
	if from.domain == domain {
		return &Vertice{id: from.id, domain: from.domain}, nil
	}
	if to.domain == domain {
		return &Vertice{id: to.id, domain: to.domain}, nil
	}
	buffer, err := v.getter.Get(from.file, from.offset, to.offset-from.offset)
	if err != nil {
		return nil, err
	}

	return findVertice(buffer, domain)
}

func findVertice(buffer []byte, domain string) (*Vertice, error) {
	reader := bytes.NewReader(buffer)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		vertice, err := LoadVertice(line)
		if err != nil {
			return nil, err
		}
		if vertice.domain == domain {
			return vertice, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %v\n", err)
	}

	return nil, nil
}
