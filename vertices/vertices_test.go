package vertices_test

import (
	"fmt"
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getVertices(t *testing.T) vertices.Vertices {
	offsets := vertices.Offsets{}
	err := offsets.Load(fmt.Sprintf("../data/%s", access.VerticesOffsetsFile))
	require.NoError(t, err)
	return vertices.NewVertices("../data/vertices", offsets)
}

func TestVerticesGet(t *testing.T) {
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
			vertice, err := v.Get(domain)
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
			vertice, err := v.Get(domain)
			require.NoError(t, err)
			assert.Nil(t, vertice, domain)
		})
	}
}
