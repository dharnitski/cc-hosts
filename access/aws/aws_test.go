package aws_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	aws_sdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Verify that S3Getter implements Getter interface.
var _ access.Getter = (*aws.S3Getter)(nil)

func newAwsConfig(t *testing.T) aws_sdk.Config {
	t.Helper()
	awsConfig, err := config.LoadDefaultConfig(t.Context())
	require.NoError(t, err)

	awsMock := middleware.InitializeMiddlewareFunc("mock", func(
		ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
	) (
		middleware.InitializeOutput, middleware.Metadata, error,
	) {
		switch v := in.Parameters.(type) {
		case *s3.GetObjectInput:
			assert.Equal(t, "test-folder/test-file", *v.Key)
			assert.Equal(t, "test-bucket", *v.Bucket)
			assert.Equal(t, "bytes=5-14", *v.Range)

			return middleware.InitializeOutput{
				Result: &s3.GetObjectOutput{
					Body:          io.NopCloser(bytes.NewReader([]byte("1234567890"))),
					ContentLength: aws_sdk.Int64(10),
				},
			}, middleware.Metadata{}, nil
		default:
			panic(fmt.Sprintf("no support for %T", v))
		}
	})

	awsConfig.APIOptions = append(awsConfig.APIOptions, func(stack *middleware.Stack) error {
		return stack.Initialize.Add(awsMock, middleware.Before)
	})

	return awsConfig
}

func TestS3Getter(t *testing.T) {
	t.Parallel()

	awsConfig := newAwsConfig(t)

	s3Getter := aws.New(awsConfig, "test-bucket", "test-folder")

	buffer, err := s3Getter.Get(t.Context(), "test-file", 5, 10)
	require.NoError(t, err)

	assert.Equal(t, "1234567890", string(buffer))
}
