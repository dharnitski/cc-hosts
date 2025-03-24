package edges

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// fileChunkSize is the size of the chunk of the file to be read in bytes
	FileChunkSize = 1024 * 1024 // 1MB
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

// save in format "domain \t offset \t file"
func (v Offset) String() string {
	return fmt.Sprintf("%s\t%d\t%s", v.id, v.offset, v.file)
}

func loadOffset(line string) (Offset, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 3 {
		return Offset{}, fmt.Errorf("Invalid line: %s, %d parts", line, len(parts))
	}
	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return Offset{}, fmt.Errorf("Invalid offset: %s", parts[1])
	}
	return Offset{offset: offset, id: parts[0], file: parts[2]}, nil
}

type Offsets struct {
	offsets []Offset
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
		return fmt.Errorf("Error creating file %q: %v\n", fileName, err)
	}
	defer file.Close()

	for _, offset := range v.offsets {
		_, err := file.WriteString(offset.String() + "\n")
		if err != nil {
			return fmt.Errorf("Error writing to file %q: %v\n", fileName, err)
		}
	}
	return nil
}

func (v *Offsets) Load(fileName string) error {
	v.offsets = make([]Offset, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error opening file %q: %v\n", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		offset, err := loadOffset(scanner.Text())
		if err != nil {
			return fmt.Errorf("Error loading offset: %v\n", err)
		}
		v.offsets = append(v.offsets, offset)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: %v\n", err)
	}
	return nil
}

func (v *Offsets) Validate() error {
	if v.Len() == 0 {
		return fmt.Errorf("No offsets found")
	}

	previousOffset := 0
	previousFile := ""
	previousID := -1
	for _, offset := range v.offsets {
		// offset
		if offset.offset < 0 {
			return fmt.Errorf("Invalid offset: %d", offset.offset)
		}
		// we reset offset when we change file
		if previousFile == offset.file {
			if offset.offset <= previousOffset {
				return fmt.Errorf("Offset goes down: %d, previous %d", offset.offset, previousOffset)
			}
		}
		previousOffset = offset.offset

		// id
		if offset.id == "" {
			return fmt.Errorf("Empty id")
		}
		id, err := strconv.Atoi(offset.id)
		if err != nil {
			return fmt.Errorf("Error converting id to integer: %v", err)
		}
		// ids are not unique in edges file, we can have multiple offsets pointing to same ID
		if id < previousID {
			return fmt.Errorf("ID goes down: %d, previous %d", id, previousID)
		}
		previousID = id

		// file
		if offset.file == "" {
			return fmt.Errorf("Empty file")
		}
		previousFile = offset.file
	}
	return nil
}
