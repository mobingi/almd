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

	"github.com/mobingi/oceand/pkg/util"
)

const (
	MessageCodeResources = "MESSAGE_CODE_LIST_RESOURCES"
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

const backendURLEnv = "BACKEND_URL"

func report(client clientset.Interface, eventChan chan struct{}, token string) {
	backendURL := util.ReadEnvOrDie(backendURLEnv)
	for {
		select {
		case <-eventChan:
			data, _ := json.Marshal(newPostBody(client))
			sendData(backendURL, token, data)
		}
	}
}

func newPostBody(client clientset.Interface) *postPayload {
	pods, _ := client.CoreV1().Pods(corev1.NamespaceAll).List(metav1.ListOptions{})
	services, _ := client.CoreV1().Services(corev1.NamespaceAll).List(metav1.ListOptions{})
	deployments, _ := client.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{})
	fmt.Println(pods)
	return &postPayload{
		MessageCode: MessageCodeResources,
		Pods:        pods.Items,
		Services:    services.Items,
		Deployments: deployments.Items,
	}
}

func sendData(backendURL, token string, data []byte) error {
	authorzation := fmt.Sprintf("Bearer %s", token)
	header := map[string][]string{
		"Authorization": {authorzation},
	}
	req, err := http.NewRequest(http.MethodPost, backendURL, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.Header = header

	_, err = http.DefaultClient.Do(req)
	// TODO rewrite it
	if err != nil {
		fmt.Println("error:", err)
	}

	return err
}
