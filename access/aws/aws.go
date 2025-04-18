package aws

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	Bucket = "common-crawl-hosts"
)

type S3Getter struct {
	client     *s3.Client
	bucketName string
	folder     string
}

func New(cfg aws.Config, bucketName string, folder string) *S3Getter {
	return &S3Getter{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
		folder:     folder,
	}
}

func (g *S3Getter) Get(ctx context.Context, fileName string, offset int, length int) ([]byte, error) {
	key := fmt.Sprintf("%s/%s", g.folder, fileName)

	if offset < 0 {
		return nil, errors.New("offset cannot be negative")
	}

	if length <= 0 {
		return nil, errors.New("length must be positive")
	}

	rangeStr := fmt.Sprintf("bytes=%d-%d", offset, offset+length-1)
	input := &s3.GetObjectInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(key),
		Range:  aws.String(rangeStr),
	}

	result, err := g.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3 bucket %s, key %s : %w", g.bucketName, key, err)
	}

	defer func() {
		if err := result.Body.Close(); err != nil {
			log.Printf("error closing S3 Body for file %s: %v", fileName, err)
		}
	}()

	if aws.ToInt64(result.ContentLength) != int64(length) {
		return nil, errors.New("unexpected content length")
	}

	buf := make([]byte, length)

	_, err = io.ReadFull(result.Body, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, nil
}
