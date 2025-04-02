package edges

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOffset(t *testing.T) {
	t.Parallel()

	line := "id123\t123\tedges.txt"
	offset, err := loadOffset(line)
	require.NoError(t, err)
	assert.Equal(t, 123, offset.offset)
	assert.Equal(t, "id123", offset.id)
	assert.Equal(t, "edges.txt", offset.file)
}

func TestLoadOffset_InvalidLine(t *testing.T) {
	t.Parallel()

	line := "invalid_line"
	_, err := loadOffset(line)
	assert.Error(t, err)
}

func TestLoadOffset_InvalidOffset(t *testing.T) {
	t.Parallel()

	line := "id123\tinvalid_offset\tedges.txt"
	_, err := loadOffset(line)
	assert.Error(t, err)
}
