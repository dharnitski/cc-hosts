package search_test

import (
	"fmt"
	"testing"

	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_GetTargets(t *testing.T) {
	t.Parallel()

	eOffsets := edges.Offsets{}
	err := eOffsets.Load(fmt.Sprintf("../data/%s", access.EdgesOffsetsFile))
	require.NoError(t, err)
	e := edges.NewEdges("../data/edges", eOffsets)

	vOffsets := vertices.Offsets{}
	err = vOffsets.Load(fmt.Sprintf("../data/%s", access.VerticesOffsetsFile))
	require.NoError(t, err)
	v := vertices.NewVertices("../data/vertices", vOffsets)

	searcher := search.NewSearcher(v, e)
	results, err := searcher.GetTargets(t.Context(), "aaa.org")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, results)
}
