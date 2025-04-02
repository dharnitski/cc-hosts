package edges

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dharnitski/cc-hosts/offsets"
)

const (
	// fileChunkSize is the size of the chunk of the file to be read in bytes.
	FileChunkSize = 1024 * 128 // 128 KB
)

type Offset struct {
	// offset in bytes to find the domain in sorted file
	offset int
	// vertice id
	// in file ot 0 based line number
	id string
	// vertices file name without path
	file string
}

func NewOffset(offset int, id, file string) Offset {
	return Offset{offset: offset, id: id, file: file}
}

// save in format "domain \t offset \t file".
func (v Offset) String() string {
	return fmt.Sprintf("%s\t%d\t%s", v.id, v.offset, v.file)
}

func (v Offset) Offset() int {
	return v.offset
}

func loadOffset(line string) (Offset, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 3 {
		return Offset{}, fmt.Errorf("invalid line: %s, %d parts", line, len(parts))
	}

	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return Offset{}, fmt.Errorf("invalid offset: %s", parts[1])
	}

	return Offset{offset: offset, id: parts[0], file: parts[2]}, nil
}

type Offsets struct {
	offsets []Offset
}

func NewOffsets() (*Offsets, error) {
	result := &Offsets{
		offsets: make([]Offset, 0),
	}
	reader := bytes.NewReader(offsets.Edges)
	err := result.loadFromReader(reader)

	return result, err
}

func NewOffsetsReversed() (*Offsets, error) {
	result := &Offsets{
		offsets: make([]Offset, 0),
	}
	reader := bytes.NewReader(offsets.EdgesReversed)
	err := result.loadFromReader(reader)

	return result, err
}

func (v *Offsets) Append(offsets []Offset) {
	v.offsets = append(v.offsets, offsets...)
}

func (v Offsets) Items() []Offset {
	return v.offsets
}

func (v Offsets) Len() int {
	return len(v.offsets)
}

func (v Offsets) Save(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating file %q: %w", fileName, err)
	}
	defer file.Close()

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

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %q: %w", fileName, err)
	}

	defer file.Close()

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

		// id
		if offset.id == "" {
			return errors.New("empty id")
		}

		id, err := strconv.Atoi(offset.id)
		if err != nil {
			return fmt.Errorf("error converting id to integer: %w", err)
		}
		// each file has sorted IDs but files are not sorted itself
		// each file can store IDs from 0 to max ID
		if previousFile == offset.file {
			// ids are not unique in edges file, we can have multiple offsets pointing to same ID
			if id < previousID {
				return fmt.Errorf("ID goes down: %d, previous %d", id, previousID)
			}
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

type TwoOffsets struct {
	From Offset
	To   Offset
}

// map of offsets with key as file name.
func (v *Offsets) offsetsMap() map[string][]Offset {
	result := make(map[string][]Offset)
	for _, offset := range v.offsets {
		result[offset.file] = append(result[offset.file], offset)
	}

	return result
}

// return from and to offsets for domain to fetch data from file.
func (v *Offsets) FindForFromID(fromID string) map[string]TwoOffsets {
	items := v.offsets
	if len(items) == 0 {
		return map[string]TwoOffsets{}
	}

	offsetsMap := v.offsetsMap()
	grouppedOffsets := make(map[string]TwoOffsets)

	for file, offsets := range offsetsMap {
		from, to := findFromFomIDInFile(fromID, offsets)
		grouppedOffsets[file] = TwoOffsets{from, to}
	}

	return grouppedOffsets
}

func findFromFomIDInFile(fromID string, offsets []Offset) (Offset, Offset) {
	items := offsets
	if len(items) == 0 {
		return Offset{}, Offset{}
	}

	inID, err := strconv.Atoi(fromID)
	if err != nil {
		return Offset{}, Offset{}
	}

	left := items[0]
	right := items[len(items)-1]

	for _, offset := range items {
		id, err := strconv.Atoi(offset.id)
		if err != nil {
			return Offset{}, Offset{}
		}

		if id < inID {
			left = offset

			continue
		}

		if id > inID {
			right = offset

			break
		}
	}

	// At this point, left > right, and the domain was not found
	// right is the greatest index with domain < target
	// left is the smallest index with domain > target
	return left, right
}
