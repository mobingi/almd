package watch

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewEventChans(t *testing.T) {
	data, err := ioutil.ReadFile("/Users/n1ce/.kube/config")
	assert.NoError(t, err)
	config, err := clientcmd.RESTConfigFromKubeConfig(data)
	assert.NoError(t, err)
	client, err := clientset.NewForConfig(config)
	assert.NoError(t, err)
	_, _, _, err = newEventChans(client)
	assert.NoError(t, err)
}
