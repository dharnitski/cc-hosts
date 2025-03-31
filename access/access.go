package access

import "context"

type Getter interface {
	Get(ctx context.Context, fileName string, offset int, length int) ([]byte, error)
}
