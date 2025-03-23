package vertices

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
	// domain name in reverse domain format
	// sample: com.example
	domain string
	// vertice id
	// in file ot 0 based line number
	id string
	// vertices file name without path
	file string
}

func NewOffset(offset int, domain, id, file string) Offset {
	return Offset{offset: offset, domain: domain, id: id, file: file}
}

// save in format "domain \t offset \t file"
func (v Offset) String() string {
	return fmt.Sprintf("%s\t%d\t%s\t%s", v.domain, v.offset, v.id, v.file)
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
		return Offset{}, fmt.Errorf("Invalid line: %s, %d parts", line, len(parts))
	}
	offset, err := strconv.Atoi(parts[1])
	if err != nil {
		return Offset{}, fmt.Errorf("Invalid offset: %s", parts[1])
	}
	return Offset{offset: offset, domain: parts[0], id: parts[2], file: parts[3]}, nil
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
	previousDomain := ""
	previousFile := ""
	previousID := ""
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

		// domain
		if offset.domain == "" {
			return fmt.Errorf("Empty domain")
		}
		if offset.domain <= previousDomain {
			return fmt.Errorf("Domain goes down: %s, previous %s", offset.domain, previousDomain)
		}
		previousDomain = offset.domain

		// id
		if offset.id == "" {
			return fmt.Errorf("Empty id")
		}
		if offset.id <= previousID {
			return fmt.Errorf("ID goes down: %s, previous %s", offset.id, previousID)
		}
		previousID = offset.id

		// file
		if offset.file == "" {
			return fmt.Errorf("Empty file")
		}
		previousFile = offset.file
	}
	return nil
}

// return from and to offsets for domain to fetch data from file
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
