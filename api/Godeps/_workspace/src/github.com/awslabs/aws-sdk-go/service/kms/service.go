package kms

import (
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/aws"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/protocol/jsonrpc"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/internal/signer/v4"
)

type KMS struct { // KMS is a client for KMS.
	*aws.Service
}

// Used for custom service initialization logic
var initService func(*aws.Service)

// Used for custom request initialization logic
var initRequest func(*aws.Request)

// New returns a new KMS client.
func New(config *aws.Config) *KMS {
	if config == nil {
		config = &aws.Config{}
	}

	service := &aws.Service{
		Config:       aws.DefaultConfig.Merge(config),
		ServiceName:  "kms",
		APIVersion:   "2014-11-01",
		JSONVersion:  "1.1",
		TargetPrefix: "TrentService",
	}
	service.Initialize()

	// Handlers
	service.Handlers.Sign.PushBack(v4.Sign)
	service.Handlers.Build.PushBack(jsonrpc.Build)
	service.Handlers.Unmarshal.PushBack(jsonrpc.Unmarshal)
	service.Handlers.UnmarshalMeta.PushBack(jsonrpc.UnmarshalMeta)
	service.Handlers.UnmarshalError.PushBack(jsonrpc.UnmarshalError)

	// Run custom service initialization if present
	if initService != nil {
		initService(service)
	}

	return &KMS{service}
}

// newRequest creates a new request for a KMS operation and runs any
// custom request initialization.
func (c *KMS) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := aws.NewRequest(c.Service, op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
