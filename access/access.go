package access

import "context"

type Getter interface {
	// Get bytes from source with offset and length
	// if offset is 0 and length is 0, return the whole file
	Get(ctx context.Context, fileName string, offset int, length int) ([]byte, error)
}
