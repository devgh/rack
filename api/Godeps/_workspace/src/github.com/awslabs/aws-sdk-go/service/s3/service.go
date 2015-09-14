package s3

import (
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/protocol/restxml"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/signer/v4"
)

type S3 struct // S3 is a client for Amazon S3.
{
	*aws.Service
}

// Used for custom service initialization logic
var initService func(*aws.Service)

// Used for custom request initialization logic
var initRequest func(*aws.Request)

// New returns a new S3 client.
func New(config *aws.Config) *S3 {
	if config == nil {
		config = &aws.Config{}
	}

	service := &aws.Service{
		Config:      aws.DefaultConfig.Merge(config),
		ServiceName: "s3",
		APIVersion:  "2006-03-01",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(restxml.Build)
	service.Handlers.Unmarshal.PushBack(restxml.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(restxml.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(restxml.UnmarshalError)

	// Run custom service initialization if present
	if initService != nil {
		initService(service)
	}

	return &S3{service}
}

// newRequest creates a new request for a S3 operation and runs any
// custom request initialization.
func (c *S3) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
