/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package kubemethod

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model/k8s"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

const (
	openstack_cluster_group     string = "cluster.x-k8s.io"
	openstack_cluster_version   string = "v1alpha3"
	openstack_cluster_resources string = "clusters"

	// Maser, Controlplanes
	openstack_controlplane_group          string = "controlplane.cluster.x-k8s.io"
	openstack_controlplane_version        string = "v1alpha3"
	openstack_controlplane_k3s_version    string = "v1beta1"
	openstack_controlplane_mk8s_version   string = "v1beta1"
	openstack_controlplane_resources      string = "kubeadmcontrolplanes"
	openstack_controlplane_k3s_resources  string = "kthreescontrolplanes"
	openstack_controlplane_mk8s_resources string = "microk8scontrolplanes"

	// Worker, MachineDeployments
	openstack_machinedeploy_group     string = "cluster.x-k8s.io"
	openstack_machinedeploy_version   string = "v1alpha3"
	openstack_machinedeploy_resources string = "machinedeployments"

	// Worker, MachineTemplate
	openstack_machinetemplate_group     string = "infrastructure.cluster.x-k8s.io"
	openstack_machinetemplate_version   string = "v1alpha3"
	openstack_machinetemplate_resources string = "openstackmachinetemplates"

	// Worker, ConfigTemplates
	openstack_configtemplate_group          string = "bootstrap.cluster.x-k8s.io"
	openstack_configtemplate_version        string = "v1alpha3"
	openstack_configtemplate_resources      string = "kubeadmconfigtemplates"
	openstack_configtemplate_k3s_version    string = "v1beta1"
	openstack_configtemplate_k3s_resources  string = "kthreesconfigtemplates"
	openstack_configtemplate_mk8s_version   string = "v1beta1"
	openstack_configtemplate_mk8s_resources string = "microk8sconfigtemplates"

	// openstack_machinesets_group     string = "machinesets.cluster.x-k8s.io"
	// openstack_machinesets_version   string = "v1alpha3"
	// openstack_machinesets_resources string = "machinesets"
)

// // findNodeSetByName - 조회된 NodeSet CR 정보에 대한 이름 기반 검색
// func findNodeSetByName(list *unstructured.UnstructuredList, clusterName, nodeSetName string) (*unstructured.Unstructured, error) {
// 	for _, item := range list.Items {
// 		if strings.Contains(item.GetName(), clusterName+"-"+name)
// 		if item.GetName() == utils.ArrayContains()
// 	}
// 	return nil, nil
// }

// checkRollingUpdateCompleted - 전달된 unstructured 에서 Status/conditions/MachinesSpecUpToDate 검증
func checkRollingUpdateCompleted(item *unstructured.Unstructured) (bool, error) {
	val, exists, err := unstructured.NestedSlice(item.Object, "status", "conditions")
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}

	// checking conditions
	for item := range val {
		logger.Info(item)
	}

	return true, nil
}

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

// patchNodeSetCount - 지정한 정보에 따른 Controlplane, MachineDeployment를 이용해 NodeCount를 변경한다.
func patchNodeSetCount(objectName, namespace, group, version, resources string, nodeCount int) error {
	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema("", group, version, resources)
	if err != nil {
		return err
	}

	// Patching replicaset to specified nodeCount
	patch := []interface{}{
		map[string]interface{}{
			"op":    "replace",
			"path":  "/spec/replicas",
			"value": nodeCount,
		},
	}

	payload, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	dynamicClient.SetNamespace(namespace)
	_, err = dynamicClient.Patch(objectName, types.JSONPatchType, bytes.NewReader(payload), metaV1.PatchOptions{})
	if err != nil {
		return err
	}

	return nil
}

// deleteGVRObject - 지정한 GVR에 해당하는 Object를 삭제한다.
func deleteGVRObject(objectName, namespace, group, version, resource string) error {
	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema("", group, version, resource)
	if err != nil {
		return err
	}

	dynamicClient.SetNamespace(namespace)
	err = dynamicClient.Delete(objectName, metaV1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// removeConfigTemplate - NodeSet에 연괸된 Config Template을 제거한다.
func removeConfigTemplate(objectName, namespace string, bootstrapProvider common.BootstrapProvider) error {
	// Get GVR
	group := openstack_configtemplate_group
	version := openstack_configtemplate_version
	resource := openstack_configtemplate_resources

	switch bootstrapProvider {
	case common.MicroK8s:
		version = openstack_configtemplate_mk8s_version
		resource = openstack_configtemplate_mk8s_resources
	case common.K3s:
		version = openstack_configtemplate_k3s_version
		resource = openstack_configtemplate_k3s_resources
	}

	return deleteGVRObject(objectName, namespace, group, version, resource)
}

// removeMachineTemplate - NedeSet에 연관된 Machine Template을 제거한다.
func removeMachineTemplate(objectName, namespace string) error {
	// Get GVR
	group := openstack_machinetemplate_group
	version := openstack_machinetemplate_version
	resource := openstack_machinetemplate_resources

	return deleteGVRObject(objectName, namespace, group, version, resource)
}

// removeMachineDeployment - NodeSet에 연관된 Machine Deployment를 제거한다.
func removeMachineDeployment(objectName, namespace string) error {
	// Get GVR
	group := openstack_machinedeploy_group
	version := openstack_machinedeploy_version
	resource := openstack_machinedeploy_resources

	return deleteGVRObject(objectName, namespace, group, version, resource)
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

// GetProvisionPhase - 지정한 클러스터에 대한 Provision Phase 검증.
func GetProvisionPhase(namespace, clusterName string) (string, error) {
	// Get kubernetes client
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
		// Checking Provision
		return checkProvisioningPhase(data)
	}
	return "", nil
}

// GetControlPlaneUpdatePhase - 지정한 클러스터에 대한 Update Phase 검증.
func GetControlPlaneUpdatePhase(namespace, masterSetName string) (bool, error) {
	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema("", openstack_controlplane_group, openstack_controlplane_version, openstack_controlplane_resources)
	if err != nil {
		return false, err
	}

	// checking the clsuter ready
	dynamicClient.SetNamespace(namespace)
	data, err := dynamicClient.Get(masterSetName, metaV1.GetOptions{})
	if err != nil {
		return false, err
	} else if data != nil {
		// Checking upgrade completed
		return checkRollingUpdateCompleted(data)
	}
	return false, nil
}

// GetPodList - description
func GetPodList(clusterName, namespace string) (*coreV1.PodList, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient(clusterName)
	if err != nil {
		return nil, err
	}

	return apiClient.CoreV1().Pods(namespace).List(context.TODO(), metaV1.ListOptions{})
}

// GetNodeList - 해당 클러스터의 노드 리스트 반환
func GetNodeList(clusterName string) ([]k8s.Node, error) {
	// Get kubernetes client
	apiClient, err := config.HostCluster.GetKubernetesClient(clusterName)
	if err != nil {
		return nil, err
	}

	nodes, err := apiClient.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return k8s.ConvertToNodeList(nodes)
}

// RemoveOpenstackProvisioned - 지정한 오픈스택 클러스터의 Provisioning 제거
func RemoveOpenstackProvisioned(clusterId, clusterName, namespace string) error {
	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema("", openstack_cluster_group, openstack_cluster_version, openstack_cluster_resources)
	if err != nil {
		return err
	}

	// checking clsuter
	dynamicClient.SetNamespace(namespace)
	data, err := dynamicClient.Get(clusterName, metaV1.GetOptions{})
	if err != nil {
		return err
	} else if data != nil {
		err = dynamicClient.Delete(clusterName, metaV1.DeleteOptions{})
	}
	return err
}

// ArrangeK8SNodesToNodeSetInfo - 지정한 클러스터의 NodeSet들에 K8S Node정보를 설정한다.
func ArrangeK8SNodesToNodeSetInfo(clusterName string, openStackNodeSetInfo model.OpenstackNodeSetInfo) bool {
	nodeList, err := GetNodeList(clusterName)
	if err != nil {
		logger.WithError(err).Warn("Provisioned, but cannot get kubernetes node info yet.")
		return true
	} else {
		// Add node info
		for _, node := range nodeList {
			find := false
			for _, nodeSet := range openStackNodeSetInfo.MasterSets {
				if strings.Contains(node.Name, "-"+nodeSet.Name+"-") {
					nodeSet.Nodes = append(nodeSet.Nodes, node)
					find = true
					break
				}
			}

			if !find {
				for _, nodeSet := range openStackNodeSetInfo.WorkerSets {
					if strings.Contains(node.Name, "-"+nodeSet.Name+"-") {
						nodeSet.Nodes = append(nodeSet.Nodes, node)
						break
					}
				}
			}
		}
	}

	return false
}

// UpdateNodeCount - 지정한 클러스터의 NodeCount 변경
func UpdateNodeCount(clusterName, nodeSetName, namespace string, bootstrapProvider common.BootstrapProvider, nodeSetType, nodeCount int) error {
	objectName := clusterName + "-" + nodeSetName

	if nodeSetType == common.NodeTypeMaster {
		// Resolve version and resource by BootstrapProvider
		resources := openstack_controlplane_resources
		version := openstack_controlplane_version
		if bootstrapProvider == common.MicroK8s {
			resources = openstack_controlplane_mk8s_resources
			version = openstack_controlplane_mk8s_version
		} else if bootstrapProvider == common.K3s {
			resources = openstack_controlplane_k3s_resources
			version = openstack_controlplane_k3s_version
		}
		return patchNodeSetCount(objectName, namespace, openstack_controlplane_group, version, resources, nodeCount)
	} else {
		return patchNodeSetCount(objectName, namespace, openstack_machinedeploy_group, openstack_machinedeploy_version, openstack_machinedeploy_resources, nodeCount)
	}
}

// RemoveNodeSet - 지정한 클러스터의 NodeSet제거
func RemoveNodeSet(clusterName, nodeSetName, namespace string, bootstrapProvider common.BootstrapProvider) error {
	objectName := clusterName + "-" + nodeSetName

	// Remove ConfigTemplate
	err := removeConfigTemplate(objectName, namespace, bootstrapProvider)
	if err != nil {
		return err
	}

	// Remove MachineTemplate
	err = removeMachineTemplate(objectName, namespace)
	if err != nil {
		return err
	}

	// Remove MachineDeployment
	err = removeMachineDeployment(objectName, namespace)
	if err != nil {
		return err
	}

	return nil
}

// // removeConfigTemplate - NodeSet에 연관된 ConfigTemplate 제거
// func removeConfigTemplate(objectName string, namespace string, BootstrapProviders) output {

// }
