package k8s

import (
	"context"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func NewClientFromFile(kubeConfPath string) (*Client, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	kc, err := os.ReadFile(filepath.Join(wd, kubeConfPath))
	if err != nil {
		return nil, err
	}

	return NewClient(string(kc))
}

// NewClient creates a new Kubernetes client from a kubeconfig YAML string
func NewClient(kubeConfigYaml string) (*Client, error) {
	// Load kubeconfig
	kubeConf, err := clientcmd.Load([]byte(kubeConfigYaml))
	if err != nil {
		return nil, err
	}

	// Construct Restclient Config
	restConf, err := clientcmd.BuildConfigFromKubeconfigGetter("",
		func() (*clientcmdapi.Config, error) {
			return kubeConf, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Initialize client
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}

	return &Client{
		kubeconf: kubeConf,
		restconf: restConf,
		client:   client,
	}, nil
}

type Client struct {
	kubeconf *clientcmdapi.Config
	restconf *rest.Config
	client   *kubernetes.Clientset
}

// ServerVersion returns the server version of the Kubernetes cluster
func (c *Client) ServerVersion() (string, error) {
	si, err := c.client.ServerVersion()
	if err != nil {
		return "", err
	}

	return si.String(), nil
}

// ListNamespaces lists all namespaces in the cluster
func (c *Client) ListNamespaces() ([]string, error) {
	namespaces, err := c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaceNames []string
	for _, namespace := range namespaces.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}
	return namespaceNames, nil
}

// ListPods lists all pods in the given namespace
func (c *Client) ListPods(namespace string) ([]string, error) {
	pods, err := c.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}

// ListContainers lists all containers in the given pod
func (c *Client) ListContainers(namespace, podName string) ([]string, error) {
	pod, err := c.client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var containerNames []string
	for _, container := range pod.Spec.Containers {
		containerNames = append(containerNames, container.Name)
	}
	return containerNames, nil
}
