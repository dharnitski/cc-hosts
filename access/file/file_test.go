package file_test

import (
	"testing"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffsetsFile(t *testing.T) {
	t.Parallel()

	getter := file.NewGetter("../../data/vertices")
	// offset from data/vertices.offsets.txt
	buffer, err := getter.Get("part-00000-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt", 2097161, 16)
	require.NoError(t, err)

	assert.Equal(t, "88296	ae.regards", string(buffer))
}
