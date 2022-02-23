package scy_test

import (
	"context"
	_ "github.com/viant/scy/kms/blowfish"
	"github.com/stretchr/testify/assert"
	"github.com/viant/scy"
	"github.com/viant/scy/cred"
	"path"
	"testing"
)

func TestService_Load(t *testing.T) {

	basePath := "/tmp/" // os.TempDir()
	var testCases = []struct {
		description string
		secret      *scy.Secret
		resource    *scy.Resource
		expect      interface{}
	}{
		{
			description: "raw secret with local fs and key",
			resource:    scy.NewResource("key", path.Join(basePath, "secret.sec"), "blowfish://default"),
			secret:      scy.NewSecret("this is secret", nil),
			expect:      "this is secret",
		},
		{
			description: "securable secrets",
			resource:    scy.NewResource(cred.Basic{}, path.Join(basePath, "json.sec"), "blowfish://default"),
			secret:      scy.NewSecret(&cred.Basic{Username: "Bob", Password: "ch@nge!Me"}, nil),
			expect:      &cred.Basic{Username: "Bob", Password: "ch@nge!Me"},
		},
	}

	for _, testCase := range testCases {
		srv := scy.New()
		ctx := context.Background()
		testCase.secret.Resource = testCase.resource
		err := srv.Store(ctx, testCase.secret)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		secret, err := srv.Load(ctx, testCase.resource)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}
		assert.EqualValues(t, testCase.expect, secret.Target, testCase.description)
	}

}
