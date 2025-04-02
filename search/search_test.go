package search_test

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"testing"

	"github.com/dharnitski/cc-hosts/access/file"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/testdata"
	"github.com/dharnitski/cc-hosts/vertices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearcher_GetTargets(t *testing.T) {
	t.Parallel()

	rootFolder := "../data"

	// cfg, err := config.LoadDefaultConfig(t.Context())
	// require.NoError(t, err)
	eOffsets, err := edges.NewOffsets()
	require.NoError(t, err)

	edgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesFolder))
	out := edges.NewEdges(edgesGetter, *eOffsets)

	offsetsReversed, err := edges.NewOffsetsReversed()
	require.NoError(t, err)

	revEdgesGetter := file.NewGetter(path.Join(rootFolder, edges.EdgesReversedFolder))
	// revEdgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesReversedFolder)
	in := edges.NewEdges(revEdgesGetter, *offsetsReversed)

	vOffsets, err := vertices.NewOffsets()
	require.NoError(t, err)

	verticesGetter := file.NewGetter(path.Join(rootFolder, vertices.Folder))
	// verticesGetter := aws.New(cfg, aws.Bucket, vertices.Folder)
	v := vertices.NewVertices(verticesGetter, *vOffsets)

	searcher := search.NewSearcher(v, out, in)
	results, err := searcher.GetTargets(t.Context(), "binaryedge.io")
	require.NoError(t, err)
	assert.Equal(t, []string{"40fy.io", "app.binaryedge.io", "blog.binaryedge.io", "cloudflare.com", "coalitioninc.com", "cyberfables.io", "d1ehrggk1349y0.cloudfront.net", "facebook.com", "fonts.googleapis.com", "github.com", "linkedin.com", "maps.googleapis.com", "slack.binaryedge.io", "support.cloudflare.com", "twitter.com"}, results.Out)
	assert.Equal(t, []string{}, results.In)
	assert.Equal(t, "binaryedge.io", results.Target)
}

func TestSearcher_Missed(t *testing.T) {
	// t.Skip()
	t.Parallel()

	inputs := testdata.GetInputs()
	// inputs = append(inputs, testdata.GetExpected()...)

	eOffsets, err := edges.NewOffsets()
	require.NoError(t, err)

	e := edges.NewEdges(file.NewGetter("../data/edges"), *eOffsets)

	reversedOffsets, err := edges.NewOffsetsReversed()
	require.NoError(t, err)

	reversed := edges.NewEdges(file.NewGetter("../data/edges_reversed"), *reversedOffsets)

	vOffsets, err := vertices.NewOffsets()
	require.NoError(t, err)

	v := vertices.NewVertices(file.NewGetter("../data/vertices"), *vOffsets)

	// TODO: Use in and out edges
	searcher := search.NewSearcher(v, e, reversed)

	out := []search.Result{}

	for _, input := range inputs {
		results, err := searcher.GetTargets(t.Context(), input)
		assert.NoError(t, err) //nolint:testifylint

		if results == nil {
			log.Printf("no results for %s\n", input)

			continue
		}

		out = append(out, *results)
	}
	// save out to JSON file
	jsonData, err := json.MarshalIndent(out, "", "    ")
	require.NoError(t, err)

	// Write to file
	err = os.WriteFile("output.json", jsonData, 0o644)
	require.NoError(t, err)
}
