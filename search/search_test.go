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
	edgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesFolder))
	revEdgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesReversedFolder))
	verticesGetter := file.NewGetter(path.Join(rootFolder, vertices.Folder))
	ctx := t.Context()
	searcher, err := search.NewSearcher(ctx, offsetsGetter, edgesGetter, revEdgesGetter, verticesGetter)
	require.NoError(t, err)
	results, err := searcher.GetTargets(ctx, "target.com")
	require.NoError(t, err)
	assert.Equal(t, []string{"out.com"}, results.Out)
	assert.Equal(t, []string{"in.com"}, results.In)
	assert.Equal(t, "target.com", results.Target)
}
