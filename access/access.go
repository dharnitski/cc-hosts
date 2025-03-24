package access

const (
	VerticesOffsetsFile = "vertices.offsets.txt"
	EdgesOffsetsFile    = "edges.offsets.txt"
)

type Getter interface {
	Get(fileName string, offset int, length int) ([]byte, error)
}
