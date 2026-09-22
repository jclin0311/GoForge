package discovery

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceDiscovery struct {
	clientset *kubernetes.Clientset
	namespace string
}

func NewServiceDiscovery(namespace string, inCluster bool) (*ServiceDiscovery, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return &ServiceDiscovery{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

func (sd *ServiceDiscovery) GetServiceEndpoint(serviceName string) (string, error) {
	ctx := context.Background()
	service, err := sd.clientset.CoreV1().Services(sd.namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get service %s: %w", serviceName, err)
	}

	if len(service.Spec.Ports) == 0 {
		return "", fmt.Errorf("service %s has no ports defined", serviceName)
	}

	port := service.Spec.Ports[0].Port
	endpoint := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, sd.namespace, port)

	return endpoint, nil
}

func (sd *ServiceDiscovery) ListServices() ([]string, error) {
	ctx := context.Background()
	services, err := sd.clientset.CoreV1().Services(sd.namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var serviceNames []string
	for _, service := range services.Items {
		serviceNames = append(serviceNames, service.Name)
	}

	return serviceNames, nil
}
