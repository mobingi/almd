package watch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	MessageCodeResources = "MESSAGE_CODE_LIST_RESOURCES"
	backendURLEnv        = "BACKEND_URL"
)

type postPayload struct {
	MessageCode string `json:"code"`

	ID           string `json:"id"`
	TemplateName string `json:"templateName"`
	ClusterName  string `json:"clusterName"`

	Pods        []corev1.Pod        `json:"pods,omitempty"`
	Deployments []appsv1.Deployment `json:"deployments,omitempty"`
	Services    []corev1.Service    `json:"services,omitempty"`
}

func report(client clientset.Interface, httpClient *http.Client, eventChan chan struct{}, backendURL, token string) {
	for {
		select {
		case <-eventChan:
			data, _ := json.Marshal(newPostBody(client))
			req := createRequest(backendURL, token, data)
			post := func() (string, error) {
				_, err := httpClient.Do(req)
				return "", err
			}
			util(post)
		}
	}
}

func newPostBody(client clientset.Interface) *postPayload {
	pods, _ := client.CoreV1().Pods(corev1.NamespaceAll).List(metav1.ListOptions{})
	services, _ := client.CoreV1().Services(corev1.NamespaceAll).List(metav1.ListOptions{})
	deployments, _ := client.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{})
	return &postPayload{
		MessageCode: MessageCodeResources,
		Pods:        pods.Items,
		Services:    services.Items,
		Deployments: deployments.Items,
	}
}

func createRequest(backendURL, token string, data []byte) *http.Request {
	authorzation := fmt.Sprintf("Bearer %s", token)
	header := map[string][]string{
		"Authorization": {authorzation},
	}
	req, err := http.NewRequest(http.MethodPost, backendURL, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header = header

	return req
}
