package file_test

import (
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ access.Getter = (*file.Getter)(nil)

func TestOffsetsFile(t *testing.T) {
	t.Parallel()

	getter := file.NewGetter("../../testdata/sample/vertices")
	buffer, err := getter.Get(t.Context(), "1.txt", 9, 9)
	require.NoError(t, err)

	assert.Equal(t, "1	com.out", string(buffer))
}

func TestOffsetsFileAll(t *testing.T) {
	t.Parallel()

	getter := file.NewGetter("../../testdata/sample/vertices")
	buffer, err := getter.Get(t.Context(), "1.txt", 0, 0)
	require.NoError(t, err)

	assert.Len(t, string(buffer), 32)
}
