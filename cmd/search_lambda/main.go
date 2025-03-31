package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dharnitski/cc-hosts/access/aws"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
)

var searcher *search.Searcher

type Request struct {
	Domain string `json:"domain"`
}

func HandleRequest(ctx context.Context, event *Request) (*search.Result, error) {
	if event == nil {
		return &search.Result{}, nil
	}
	return searcher.GetTargets(ctx, event.Domain)
}

func createSearcher(ctx context.Context) (*search.Searcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	eOffsets, err := edges.NewOffsets()
	if err != nil {
		return nil, err
	}
	edgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesFolder)
	out := edges.NewEdges(edgesGetter, *eOffsets)

	reversedOffsets, err := edges.NewOffsetsReversed()
	if err != nil {
		return nil, err
	}
	revEdgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesReversedFolder)
	in := edges.NewEdges(revEdgesGetter, *reversedOffsets)

	vOffsets, err := vertices.NewOffsets()
	if err != nil {
		return nil, err
	}
	verticesGetter := aws.New(cfg, aws.Bucket, vertices.Folder)
	v := vertices.NewVertices(verticesGetter, *vOffsets)

	searcher := search.NewSearcher(v, out, in)
	return searcher, nil
}

func main() {
	var err error
	ctx := context.Background()
	searcher, err = createSearcher(ctx)
	if err != nil {
		panic(err)
	}
	lambda.Start(HandleRequest)
}
