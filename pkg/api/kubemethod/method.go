/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package kubemethod

import (
	"context"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetKubeconfig - 지정한 클러스터에 대한 Kubeconfig 추출
func GetKubeconfig(namespace, secretName, keyName string) (string, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient("")
	if err != nil {
		return "", err
	}

	secret, err := apiClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	}

	if secret.Data[keyName] == nil {
		return "", nil
	}

	return string(secret.Data[keyName]), nil
}

// GetPodList - description
func GetPodList(clusterId, namespace string) (*coreV1.PodList, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient(clusterId)
	if err != nil {
		return nil, err
	}

	return apiClient.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{})
}
