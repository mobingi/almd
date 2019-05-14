package watch

import (
	"io/ioutil"
	"net"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	pkgutil "github.com/mobingi/oceand/pkg/util"
)

const chanSize = 100

// Watch will watch pod,service,deployment from apiserver and report to backend server
// it will return immediately
func Watch(id, token string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	httpClient := newHTTPClient()
	backendURL := pkgutil.ReadEnvOrDie(backendURLEnv)
	accessToken := getAccessToken(httpClient, backendURL, id, token)

	podEventChan, serviceEventChan, deploymentEventChan, err := newEventChans(client)
	if err != nil {
		return err
	}

	eventChan := make(chan struct{}, chanSize)
	go func() {
		for {
			select {
			case <-podEventChan:
				eventChan <- struct{}{}
			case <-serviceEventChan:
				eventChan <- struct{}{}
			case <-deploymentEventChan:
				eventChan <- struct{}{}
			}
		}
	}()

	go report(client, httpClient, eventChan, backendURL, accessToken)

	return nil
}

// new kubernetes client, use serviceaccount for auth
func newClient() (clientset.Interface, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
		host       = "kubernetes.default"
		port       = "443"
	)

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		Host: "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: rootCAFile,
		},
		BearerToken: string(token),
	}

	return clientset.NewForConfig(config)
}

func newEventChans(client clientset.Interface) (<-chan watch.Event, <-chan watch.Event, <-chan watch.Event, error) {
	podWatcher, err := client.CoreV1().Pods(corev1.NamespaceAll).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	serviceWatcher, err := client.CoreV1().Services(corev1.NamespaceAll).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, nil, nil, err
	}
	deploymentWatcher, err := client.AppsV1().Deployments(corev1.NamespaceAll).Watch(metav1.ListOptions{})
	if err != nil {
		return nil, nil, nil, err
	}

	podEventChan := podWatcher.ResultChan()
	serviceEventChan := serviceWatcher.ResultChan()
	deploymentEventChan := deploymentWatcher.ResultChan()

	return podEventChan, serviceEventChan, deploymentEventChan, nil
}
