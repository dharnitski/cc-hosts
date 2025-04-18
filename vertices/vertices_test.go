package vertices_test

import (
	"testing"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getVertices(t *testing.T) *vertices.Vertices {
	t.Helper()

	offsets, err := vertices.NewOffsets()
	require.NoError(t, err)

	return vertices.NewVertices(file.NewGetter("../data/vertices"), *offsets)
}

func TestVerticesGetByDomain(t *testing.T) {
	t.Parallel()

	v := getVertices(t)

	tests := []string{
		"aaa.11111",
		"ae.regards",
		"com.example",
		"org.example",
		"abbott.at",
		"com.amoblog.resmedairsense10autoset94814",
		"zw.zzs.th.ac.lpru.arounduniversity.ixiz.qoo",
	}

	for _, domain := range tests {
		t.Run(domain, func(t *testing.T) {
			t.Parallel()
			vertice, err := v.GetByDomain(t.Context(), domain)
			require.NoError(t, err)
			require.NotNil(t, vertice, domain)
			assert.Equal(t, domain, vertice.Domain())
			assert.GreaterOrEqual(t, vertice.ID(), "0", domain)
		})
	}
}

func TestVerticesGetNil(t *testing.T) {
	t.Parallel()

	v := getVertices(t)

	tests := []string{
		"random.random",
		"com.dharnitski",
	}

	for _, domain := range tests {
		t.Run(domain, func(t *testing.T) {
			t.Parallel()
			vertice, err := v.GetByDomain(t.Context(), domain)
			require.NoError(t, err)
			assert.Nil(t, vertice, domain)
		})
	}
}

func TestVerticesGetByID(t *testing.T) {
	t.Parallel()

	v := getVertices(t)

	tests := []string{
		"0",
		"119",
		"283704017",
		"283704060",
	}

	for _, id := range tests {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			vertice, err := v.GetByID(t.Context(), id)
			require.NoError(t, err)
			require.NotNil(t, vertice, id)
			assert.Equal(t, id, vertice.ID())
		})
	}
}

func TestVerticesGetByIDs(t *testing.T) {
	t.Parallel()

	v := getVertices(t)

	ids := []string{
		"0",
		"119",
		"283704017",
		"283704060",
	}
	vertices, err := v.GetByIDs(t.Context(), ids)
	require.NoError(t, err)
	require.NotNil(t, vertices)
	assert.Len(t, vertices, 4)
}
