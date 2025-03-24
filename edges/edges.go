package edges

import (
	"fmt"
	"strings"
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
