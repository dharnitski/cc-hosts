package vertices

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/offsets"
)

const (
	// fileChunkSize is the size of the chunk of the file to be read in bytes.
	FileChunkSize = 1024 * 32 // 32 KB
)

type Offset struct {
	// offset in bytes to find the domain in sorted file
	// points to string start in `file`for `domain`
	offset int
	// domain name in reverse domain format
	// sample: com.example
	domain string
	// vertex id
	// in file it is 0 based line number
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

func NewOffsets(ctx context.Context, getter access.Getter) (*Offsets, error) {
	result := &Offsets{
		offsets: make([]Offset, 0),
	}

	data, err := getter.Get(ctx, offsets.VerticesOffsetsFile, 0, 0)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)
	err = result.loadFromReader(reader)

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

// Load offsets from file.
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
	return find(v.offsets, domain, func(a, b string) bool { return a < b }, func(o Offset) string { return o.domain })
}

func (v *Offsets) FindForID(id int) (Offset, Offset) {
	return find(v.offsets, id, func(a, b int) bool { return a < b }, func(o Offset) int { return o.id })
}

// Generic binary search function.
func find[T comparable](items []Offset, target T, less func(T, T) bool, getField func(Offset) T) (Offset, Offset) {
	if len(items) == 0 {
		return Offset{}, Offset{}
	}

	// Binary search implementation
	left := 0
	right := len(items) - 1

	// If target is outside our range, return appropriate bounds
	if less(target, getField(items[left])) {
		return Offset{}, items[left]
	}

	if less(getField(items[right]), target) {
		return items[right], Offset{}
	}

	// Binary search
	for left <= right {
		mid := left + (right-left)/2

		if getField(items[mid]) == target {
			// Exact match found
			return items[mid], items[mid]
		}

		if less(getField(items[mid]), target) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// At this point, left > right, and the target was not found
	// right is the greatest index with field < target
	// left is the smallest index with field > target
	return items[right], items[left]
}
