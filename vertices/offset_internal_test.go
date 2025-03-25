package vertices

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadVerticesOffset(t *testing.T) {
	t.Parallel()
	line := "example.com\t123\t42\tvertices.txt"
	vo, err := loadOffset(line)
	require.NoError(t, err)

	assert.Equal(t, 123, vo.offset)
	assert.Equal(t, "example.com", vo.domain)
	assert.Equal(t, 42, vo.id)
	assert.Equal(t, "vertices.txt", vo.file)
}

func TestLoadVerticesOffset_InvalidLine(t *testing.T) {
	t.Parallel()
	line := "invalid_line"
	_, err := loadOffset(line)
	require.Error(t, err)
}

func TestLoadVerticesOffset_InvalidOffset(t *testing.T) {
	t.Parallel()
	line := "example.com\tinvalid_offset\tvertices.txt"
	_, err := loadOffset(line)
	require.Error(t, err)
}
