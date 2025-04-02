package vertices

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dharnitski/cc-hosts/offsets"
)

const (
	// fileChunkSize is the size of the chunk of the file to be read in bytes.
	FileChunkSize = 1024 * 32 // 32 KB
)

type Offset struct {
	// offset in bytes to find the domain in sorted file
	offset int
	// domain name in reverse domain format
	// sample: com.example
	domain string
	// vertice id
	// in file ot 0 based line number
	id int
	// vertices file name without path
	// TODO: this is not memory efficient structure, the same string repeated many times and uses memory for copies
	file string
}

func NewOffset(offset int, domain string, id int, file string) Offset {
	return Offset{offset: offset, domain: domain, id: id, file: file}
}

// save in format "domain \t offset \t file".
func (v Offset) String() string {
	return fmt.Sprintf("%s\t%d\t%d\t%s", v.domain, v.offset, v.id, v.file)
}

func (v Offset) Offset() int {
	return v.offset
}

func (v Offset) Domain() string {
	return v.domain
}

func loadOffset(line string) (Offset, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 4 {
		return Offset{}, fmt.Errorf("invalid line: %s, %d parts", line, len(parts))
	}

	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return Offset{}, fmt.Errorf("invalid offset: %s", parts[1])
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return Offset{}, fmt.Errorf("invalid id: %s", parts[2])
	}

	return Offset{offset: offset, domain: parts[0], id: id, file: parts[3]}, nil
}

type Offsets struct {
	offsets []Offset
}

func NewOffsets() (*Offsets, error) {
	result := &Offsets{
		offsets: make([]Offset, 0),
	}
	reader := bytes.NewReader(offsets.Vertices)
	err := result.loadFromReader(reader)

	return result, err
}

func (v *Offsets) Append(offsets []Offset) {
	v.offsets = append(v.offsets, offsets...)
}

func (v *Offsets) Items() []Offset {
	return v.offsets
}

func (v *Offsets) Len() int {
	return len(v.offsets)
}

func (v *Offsets) Save(fileName string) error {
	file, err := os.Create(fileName) //nolint:gosec
	if err != nil {
		return fmt.Errorf("error creating file %q: %w", fileName, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file %s: %v", fileName, err)
		}
	}()

	for _, offset := range v.offsets {
		_, err := file.WriteString(offset.String() + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file %q: %w", fileName, err)
		}
	}

	return nil
}

func (v *Offsets) Load(fileName string) error {
	v.offsets = make([]Offset, 0)

	file, err := os.Open(fileName) //nolint:gosec
	if err != nil {
		return fmt.Errorf("error opening file %q: %w", fileName, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file %s: %v", fileName, err)
		}
	}()

	return v.loadFromReader(file)
}

func (v *Offsets) loadFromReader(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		offset, err := loadOffset(scanner.Text())
		if err != nil {
			return fmt.Errorf("error loading offset: %w", err)
		}

		v.offsets = append(v.offsets, offset)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func (v *Offsets) Validate() error {
	if v.Len() == 0 {
		return errors.New("no offsets found")
	}

	previousOffset := 0
	previousDomain := ""
	previousFile := ""
	previousID := -1

	for _, offset := range v.offsets {
		// offset
		if offset.offset < 0 {
			return fmt.Errorf("invalid offset: %d", offset.offset)
		}
		// we reset offset when we change file
		if previousFile == offset.file {
			if offset.offset <= previousOffset {
				return fmt.Errorf("offset goes down: %d, previous %d", offset.offset, previousOffset)
			}
		}

		previousOffset = offset.offset

		// domain
		if offset.domain == "" {
			return errors.New("empty domain")
		}

		if offset.domain <= previousDomain {
			return fmt.Errorf("domain goes down: %s, previous %s", offset.domain, previousDomain)
		}

		previousDomain = offset.domain

		// id
		id := offset.id
		if id <= previousID {
			return fmt.Errorf("ID goes down: %d, previous %d", id, previousID)
		}

		previousID = id

		// file
		if offset.file == "" {
			return errors.New("empty file")
		}

		previousFile = offset.file
	}

	return nil
}

// return from and to offsets for domain to fetch data from file.
func (v *Offsets) FindForDomain(domain string) (Offset, Offset) {
	items := v.offsets
	if len(items) == 0 {
		return Offset{}, Offset{}
	}

	// Binary search implementation
	left := 0
	right := len(items) - 1

	// If domain is outside our range, return appropriate bounds
	if domain < items[left].domain {
		return Offset{}, items[left]
	}

	if domain > items[right].domain {
		return items[right], Offset{}
	}

	// Binary search
	for left <= right {
		mid := left + (right-left)/2

		if items[mid].domain == domain {
			// Exact match found
			return items[mid], items[mid]
		}

		if items[mid].domain < domain {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// At this point, left > right, and the domain was not found
	// right is the greatest index with domain < target
	// left is the smallest index with domain > target
	return items[right], items[left]
}

func (v *Offsets) FindForID(id int) (Offset, Offset) {
	items := v.offsets
	if len(items) == 0 {
		return Offset{}, Offset{}
	}

	// Binary search implementation
	left := 0
	right := len(items) - 1

	// If id is outside our range, return appropriate bounds
	if id < items[left].id {
		return Offset{}, items[left]
	}

	if id > items[right].id {
		return items[right], Offset{}
	}

	// Binary search
	for left <= right {
		mid := left + (right-left)/2

		if items[mid].id == id {
			// Exact match found
			return items[mid], items[mid]
		}

		if items[mid].id < id {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// At this point, left > right, and the id was not found
	// right is the greatest index with id < target
	// left is the smallest index with id > target
	return items[right], items[left]
}
