package access

import "context"

const (
	VerticesOffsetsFile     = "vertices.offsets.txt"
	EdgesOffsetsFile        = "edges.offsets.txt"
	EdgesReversedOffsetFile = "edges-reversed.offsets.txt"
)

type Getter interface {
	Get(ctx context.Context, fileName string, offset int, length int) ([]byte, error)
}
