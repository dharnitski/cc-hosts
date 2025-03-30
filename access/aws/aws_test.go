package aws_test

import (
	"github.com/dharnitski/cc-hosts/access"
	"github.com/dharnitski/cc-hosts/access/aws"
)

// Verify that S3Getter implements Getter interface
var _ access.Getter = (*aws.S3Getter)(nil)
