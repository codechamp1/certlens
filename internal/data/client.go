package data

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset kubernetes.Interface
}

type SecretsFetcher interface {
	FetchSecrets(namespace string) (*corev1.SecretList, error)
	FetchSecret(namespace, name string) (*corev1.Secret, error)
}

func newClient(kubeconfig, context string) (*Client, error) {
	config, err := buildConfigWithContext(context, kubeconfig)

	if err != nil {
		return nil, fmt.Errorf("cant build the k8s config with the used kubeconfig and context: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return nil, fmt.Errorf("cant build the k8s client with the used kubeconfig: %w", err)
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func NewDefaultSecretsFetcher(kubeconfig, context string) (SecretsFetcher, error) {
	client, err := newClient(kubeconfig, context)

	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	return client, nil
}

func (c Client) FetchSecrets(namespace string) (*corev1.SecretList, error) {
	secrets, err := c.clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("error listing secrets in namespace %s: %w", namespace, err)
	}

	return secrets, nil
}

func (c Client) FetchSecret(namespace, name string) (*corev1.Secret, error) {
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})

	if err != nil {
		return nil, fmt.Errorf("error fetching secret %s in namespace %s: %w", name, namespace, err)
	}

	return secret, nil
}

func buildConfigWithContext(context string, kubeconfigPath string) (*rest.Config, error) {
	var loadingRules *clientcmd.ClientConfigLoadingRules
	if kubeconfigPath != "" {
		loadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	} else {
		loadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loadingRules,
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
