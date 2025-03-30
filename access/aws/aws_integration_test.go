package aws_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/dharnitski/cc-hosts/access/aws"
	"github.com/dharnitski/cc-hosts/vertices"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOffsetBucket(t *testing.T) {
	t.Skip()
	t.Parallel()
	cfg, err := config.LoadDefaultConfig(t.Context())
	require.NoError(t, err)
	getter := aws.New(cfg, aws.Bucket, vertices.Folder)

	buffer, err := getter.Get(t.Context(), "part-00000-4ba7987d-67a0-4f7d-b410-1d92df440699-c000.txt", 2097161, 16)
	require.NoError(t, err)

	assert.Equal(t, "88296	ae.regards", string(buffer))
}
