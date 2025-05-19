package search_test

import (
	"path"
	"testing"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_GetTargets(t *testing.T) {
	t.Parallel()

	rootFolder := "../testdata/sample"

	offsetsGetter := file.NewGetter(rootFolder)
	eOffsets, err := edges.NewOffsets(t.Context(), offsetsGetter)
	require.NoError(t, err)

	edgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesFolder))
	out := edges.NewEdges(edgesGetter, *eOffsets)

	offsetsReversed, err := edges.NewOffsetsReversed(t.Context(), offsetsGetter)
	require.NoError(t, err)

	revEdgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesReversedFolder))
	// revEdgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesReversedFolder)
	in := edges.NewEdges(revEdgesGetter, *offsetsReversed)

	vOffsets, err := vertices.NewOffsets(t.Context(), offsetsGetter)
	require.NoError(t, err)

	verticesGetter := file.NewGetter(path.Join(rootFolder, vertices.Folder))
	// verticesGetter := aws.New(cfg, aws.Bucket, vertices.Folder)
	v := vertices.NewVertices(verticesGetter, *vOffsets)

	searcher := search.NewSearcher(v, out, in)
	results, err := searcher.GetTargets(t.Context(), "target.com")
	require.NoError(t, err)
	assert.Equal(t, []string{"out.com"}, results.Out)
	assert.Equal(t, []string{"in.com"}, results.In)
	assert.Equal(t, "target.com", results.Target)
}
