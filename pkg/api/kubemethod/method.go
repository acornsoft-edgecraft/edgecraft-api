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
)

// // SplitResourcesFromYaml - Yaml 포맷 문자열에 존재하는 리소스들을 분린 반환
// func SplitResourcesFromYaml(yamlString string) ([][]byte, error) {
// 	dec := yaml.NewDecoder(strings.NewReader(yamlString))
// 	var res [][]byte
// 	for {
// 		var value interface{}
// 		err := dec.Decode(&value)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		valueBytes, err := yaml.Marshal(value)
// 		if err != nil {
// 			return nil, err
// 		}
// 		res = append(res, valueBytes)
// 	}

// 	return res, nil
// }

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
