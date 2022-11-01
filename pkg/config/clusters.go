/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package config

import (
	"strings"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/client"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	IN_CLUSTER_NAME             = "kubernetes@in-cluster"
	KubeConfigStrategyFile      = "file"
	KubeConfigStrategyConfigmap = "configmap"
)

var (
	ConfigStrategy = &strategyInfo{Data: make(map[string]string)}
	HostCluster    *kubeCluster
)

// strategyInfo - kubeconfig 관리 전략
type strategyInfo struct {
	Strategy string            // file, configmap
	Data     map[string]string // file - path, configmap - namespace, configmap name, file name
}

// kubeCluster - Kubernetes Cluster 관리 정보
type kubeCluster struct {
	KubeConfig         *api.Config                          // kubeconfig used by client-go
	IsRunningInCluster bool                                 // edgecraft default cluster 여부 (In-Cluster Mode)
	InCluster          *ClientSet                           // ClientSet on In-Cluster
	DefaultContext     string                               // kubeconfig file - default context
	ClusterNames       []string                             // Context list
	Add                func([]byte) error                   // Add cluster's kubeconfig
	Remove             func(clusterName string) error       // Remove cluster's kubeconfig
	Save               func() error                         // Save kubeconfig
	Client             func(ctx string) (*ClientSet, error) // Get cluster's client
	clients            map[string]*ClientSet                // cluster's client (using rest.Config)
}

// checkClusterName - 클러스터 명을 검증하고, 없는 경우는 DefaultContext 반환
func (kc *kubeCluster) checkClusterName(clusterName string) string {
	if clusterName == "" {
		return HostCluster.DefaultContext
	}
	return clusterName
}

// GetClientSet - 클러스터명에 해당하는 ClientSet 반환
func (kc *kubeCluster) GetClientSet(clusterName string) (*ClientSet, error) {
	clusterName = kc.checkClusterName(clusterName)
	return HostCluster.Client(clusterName)
}

// GetKubernetesClient - 지정한 클러스터에 대한 Kubernetes client 반환
func (kc *kubeCluster) GetKubernetesClient(clusterName string) (*kubernetes.Clientset, error) {
	clientSet, err := kc.GetClientSet(clusterName)
	if err != nil {
		return nil, err
	}
	return clientSet.NewKubernetesClient()
}

// GetDynamicClient - 지정한 클러스터에 대한 Kubernetes dynamic client 반환
func (kc *kubeCluster) GetDynamicClient(clusterName string) (*client.DynamicClient, error) {
	clientSet, err := kc.GetClientSet(clusterName)
	if err != nil {
		return nil, err
	}
	return clientSet.NewDynamicClient()
}

// ClientSet - Kubernetes 연계용 Client 정보
type ClientSet struct {
	Name                   string
	RESTConfig             *rest.Config
	NewKubernetesClient    func() (*kubernetes.Clientset, error)
	NewDiscoveryClient     func() (*discovery.DiscoveryClient, error)
	NewDynamicClient       func() (*client.DynamicClient, error)
	NewDynamicClientSchema func(group, version, resource string) (*client.DynamicClient, error)
}

// init - Initialize on package load
func init() {}

// SetupClusters - EdgeCraft에서 관리할 클러스터들에 대한 연계 설정
func SetupClusters(configParam string) {
	var err error

	// strategy에 따라 구성 설정
	if configParam == "" || !strings.Contains(configParam, "strategy=") {
		ConfigStrategy.Strategy = "file"
		ConfigStrategy.Data["path"] = ""
	} else {
		for _, e := range strings.Split(configParam, ",") {
			parts := strings.Split(e, "=")
			if parts[0] == "strategy" {
				ConfigStrategy.Strategy = parts[1]
			} else {
				ConfigStrategy.Data[parts[0]] = parts[1]
			}
		}
	}

	// Host Cluster 구성
	if HostCluster, err = newKubeCluster(ConfigStrategy); err != nil {
		logger.WithError(err).Error("can't setup kubernetes clusters")
	} else {
		HostCluster.IsRunningInCluster = (len(HostCluster.ClusterNames) == 1 && HostCluster.InCluster != nil)
		if len(HostCluster.ClusterNames) == 0 {
			logger.Warnf("Initialized empty clusters (kubeconfig=none, in-cluster=none running-in-cluster=none)")
		} else {
			logger.Infof("Initialzied clusters (kubeconfig-strategy=%s, in-cluster=%t, running-in-cluster=%t, contexts=%s, default-context=%s)",
				ConfigStrategy.Strategy, (HostCluster.InCluster != nil), HostCluster.IsRunningInCluster, HostCluster.ClusterNames, HostCluster.DefaultContext)
		}
	}
}
