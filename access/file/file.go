package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Getter struct {
	folder string
}

func NewGetter(folder string) *Getter {
	return &Getter{folder: folder}
}

func (f *Getter) Get(ctx context.Context, fileName string, offset int, length int) ([]byte, error) {
	fullName := filepath.Join(f.folder, fileName)
	file, err := os.OpenFile(fullName, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, length)
	ret, err := file.Seek(int64(offset), 0)
	if err != nil {
		return nil, fmt.Errorf("Failed to seek file: %w", err)
	}
	if ret != int64(offset) {
		return nil, fmt.Errorf("Failed to seek file: expected %d bytes, read %d", offset, ret)
	}
	n, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %w", err)
	}
	if n != length {
		return nil, fmt.Errorf("Failed to read file: expected %d bytes, read %d", length, n)
	}
	return buffer, nil
}
