package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dharnitski/cc-hosts/access/aws"
	"github.com/dharnitski/cc-hosts/edges"
	"github.com/dharnitski/cc-hosts/search"
	"github.com/dharnitski/cc-hosts/vertices"
)

var searcher *search.Searcher //nolint:gochecknoglobals

type Request struct {
	Domain string `json:"domain"`
}

// handler for basic lambda function.
func HandleRequest(ctx context.Context, event *Request) (*search.Result, error) {
	if event == nil {
		return &search.Result{}, nil
	}

	return searcher.GetTargets(ctx, event.Domain)
}

// handler for API Gateway proxy for lambda function.
func HandleGateway(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	domain, ok := request.PathParameters["domain"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing domain parameter in path",
		}, nil
	}

	response, err := searcher.GetTargets(ctx, domain)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error marshalling response",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonResponse),
	}, nil
}

func createSearcher(ctx context.Context) (*search.Searcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// folder is empty - expect offset files in the root of the bucket
	offsetsGetter := aws.New(cfg, aws.Bucket, "")
	edgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesFolder)
	revEdgesGetter := aws.New(cfg, aws.Bucket, edges.EdgesReversedFolder)
	verticesGetter := aws.New(cfg, aws.Bucket, vertices.Folder)

	return search.NewSearcher(ctx, offsetsGetter, edgesGetter, revEdgesGetter, verticesGetter)
}

func main() {
	var err error

	ctx := context.Background()
	// short timeout to load offsets
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	searcher, err = createSearcher(ctx)
	if err != nil {
		panic(err)
	}

	lambda.Start(HandleGateway)
}
