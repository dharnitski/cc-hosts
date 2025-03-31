package offsets

import _ "embed"

const (
	Folder                  = "offsets"
	VerticesOffsetsFile     = "vertices.offsets.txt"
	EdgesOffsetsFile        = "edges.offsets.txt"
	EdgesReversedOffsetFile = "edges-reversed.offsets.txt"
)

//go:embed vertices.offsets.txt
var Vertices []byte

//go:embed edges.offsets.txt
var Edges []byte

//go:embed edges-reversed.offsets.txt
var EdgesReversed []byte
