package file

import (
	"context"
	"fmt"
	"log"
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

	file, err := os.OpenFile(fullName, os.O_RDONLY, 0o644) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file %s: %v", fileName, err)
		}
	}()

	buffer := make([]byte, length)

	ret, err := file.Seek(int64(offset), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek file: %w", err)
	}

	if ret != int64(offset) {
		return nil, fmt.Errorf("failed to seek file: expected %d bytes, read %d", offset, ret)
	}

	n, err := file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if n != length {
		return nil, fmt.Errorf("failed to read file: expected %d bytes, read %d", length, n)
	}

	return buffer, nil
}
