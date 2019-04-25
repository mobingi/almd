package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptionsFromFile(t *testing.T) {
	o, err := NewOptionsFromFile("testdata/config-test.yaml")
	assert.NoError(t, err)
	assert.Equal(t, "rootuser-subuser", o.ID)
	assert.Equal(t, "token", o.Token)
	assert.Equal(t, "template name", o.TemplateName)
	assert.Equal(t, "k8s cluster name", o.ClusterName)
}
