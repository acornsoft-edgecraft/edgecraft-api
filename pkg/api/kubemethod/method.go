/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package kubemethod

import (
	"context"
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/k8s"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	//zeroDuration time.Duration = 0

	// openstack_capi_group     string = "infrastructure.cluster.x-k8s.io"
	// openstack_capi_version   string = "v1alpha3"
	// openstack_capi_resources string = "openstackclusters"

	openstack_cluster_group     string = "cluster.x-k8s.io"
	openstack_cluster_version   string = "v1alpha3"
	openstack_cluster_resources string = "clusters"
)

// checkProvisioningPhase - 전달된 unstructured 에서 Statue/Ready 상태 검증
func checkProvisioningPhase(item *unstructured.Unstructured) (string, error) {
	val, exists, err := unstructured.NestedString(item.Object, "status", "phase")
	if err != nil {
		return "", err
	} else if !exists {
		return "", nil
	}

	return val, nil
}

// Apply - 지정한 YAML 문자열 정보를 Kubernetes에 적용
func Apply(clusterName, yaml string) error {
	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClient("")
	if err != nil {
		return err
	}

	res, err := dynamicClient.OpenstackProvisionPost(strings.NewReader(yaml))
	if err != nil {
		return err
	}

	logger.Infof("Dynamic Apply processed: %v", res)

	return nil
}

// GetKubeconfig - 지정한 클러스터에 대한 Kubeconfig 추출
func GetKubeconfig(namespace, clusterName, keyName string) (string, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient("")
	if err != nil {
		return "", err
	}

	secret, err := apiClient.CoreV1().Secrets(namespace).Get(context.TODO(), clusterName+"-kubeconfig", metaV1.GetOptions{})
	if err != nil {
		return "", err
	}

	if secret.Data[keyName] == nil {
		return "", nil
	}

	return string(secret.Data[keyName]), nil
}

// GetProvisioned - 지정한 클러스터에 대한 Provision 상태 검증.
func GetProvisioned(namespace, clusterName string) (string, error) {
	// Get kubernetes client
	//dynamicClient, err := config.HostCluster.GetDynamicClient("")
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema("", openstack_cluster_group, openstack_cluster_version, openstack_cluster_resources)
	if err != nil {
		return "", err
	}

	// checking the clsuter ready
	dynamicClient.SetNamespace(namespace)
	data, err := dynamicClient.Get(clusterName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	} else if data != nil {
		// TODO: Checking Provision
		return checkProvisioningPhase(data)
	}
	return "", nil
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

// GetNodeList - 해당 클러스터의 노드 리스트 반환
func GetNodeList(clusterId string) ([]k8s.Node, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient(clusterId)
	if err != nil {
		return nil, err
	}

	nodes, err := apiClient.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return k8s.ConvertToNodeList(nodes)
}
