// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package ecsiface_test

import (
	"testing"

	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/service/ecs"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/awslabs/aws-sdk-go/service/ecs/ecsiface"
	"github.com/convox/rack/api/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	assert.Implements(t, (*ecsiface.ECSAPI)(nil), ecs.New(nil))
}
