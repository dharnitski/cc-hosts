package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Getter struct {
	client     *s3.Client
	bucketName string
}

func New(cfg aws.Config, bucketName string) *S3Getter {
	return &S3Getter{
		client:     s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}
}

func (g *S3Getter) Get(fileName string, offset int, length int) ([]byte, error) {
	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative")
	}
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}

	rangeStr := fmt.Sprintf("bytes=%d-%d", offset, offset+length-1)
	input := &s3.GetObjectInput{
		Bucket: aws.String(g.bucketName),
		Key:    aws.String(fileName),
		Range:  aws.String(rangeStr),
	}

	result, err := g.client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	buf := make([]byte, length)
	_, err = result.Body.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return buf, nil
}
