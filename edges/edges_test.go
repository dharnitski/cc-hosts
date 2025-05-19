package edges_test

import (
	"testing"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getEdges(t *testing.T) *edges.Edges {
	t.Helper()

	offsets, err := edges.NewOffsets(t.Context(), file.NewGetter("../testdata/sample"))
	require.NoError(t, err)

	return edges.NewEdges(file.NewGetter("../testdata/sample/edges"), *offsets)
}

func TestEdgesGet(t *testing.T) {
	t.Parallel()

	v := getEdges(t)

	tests := []struct {
		id       string
		expected []string
	}{
		{
			id: "2",
			expected: []string{
				"1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			t.Parallel()
			ids, err := v.Get(t.Context(), tt.id)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, ids)
		})
	}
}
