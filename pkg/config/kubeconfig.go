/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package config

import (
	"context"
	"fmt"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/client"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// clusterConfigProvider - Provider for cluster config
type clusterConfigProvider struct {
	load func() (*api.Config, error)
	save func(*api.Config) error
}

// fileClusterConfigProvider - 파일 기반 cluster config
type fileClusterConfigProvider struct {
	clusterConfigProvider
	defaultFilename string
}

// parsingConfig - 바이트 배열을 kubeconfig 정보로 파싱
func parsingConfig(conf []byte) (*api.Config, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes(conf)
	if err != nil {
		return nil, err
	}

	apiConfig, err := clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	return apiConfig.DeepCopy(), nil
}

// createFileClusterConfig - 파일 기반의 cluster config 생성
func createFileClusterConfig(kubeconfig string) (*clusterConfigProvider, error) {
	conf := &fileClusterConfigProvider{}
	conf.load = func() (*api.Config, error) {
		var configLoadingRules clientcmd.ClientConfigLoader
		if kubeconfig != "" {
			configLoadingRules = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
		} else {
			configLoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
		}
		apiConfig, err := configLoadingRules.Load()
		conf.defaultFilename = configLoadingRules.GetDefaultFilename()
		if err != nil {
			return nil, err
		}
		return apiConfig.DeepCopy(), nil
	}
	conf.save = func(apiConfig *api.Config) error {
		if err := clientcmd.WriteToFile(*apiConfig, conf.defaultFilename); err != nil {
			return err
		}
		return nil
	}

	return &conf.clusterConfigProvider, nil
}

// createConfigmapClusterConfig - configmap 기반의 cluster config 생성
func createConfigmapClusterConfig(nsName, cmName, fileName string) (*clusterConfigProvider, error) {
	conf := &clusterConfigProvider{}
	conf.load = func() (*api.Config, error) {
		rest, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		clientset, _ := kubernetes.NewForConfig(rest)
		cm, err := clientset.CoreV1().ConfigMaps(nsName).Get(context.TODO(), cmName, v1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if cm.Data[fileName] == "" {
			return nil, fmt.Errorf("kubeconfig data is empty namespace=%s, configmap=%s, filename=%s, data=%v", nsName, cmName, fileName, cm.Data)
		}
		clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(cm.Data[fileName]))
		if err != nil {
			return nil, err
		}

		apiConfig, err := clientConfig.RawConfig()
		if err != nil {
			return nil, err
		}

		return apiConfig.DeepCopy(), nil
	}
	conf.save = func(apiConfig *api.Config) error {
		rest, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		clientset, _ := kubernetes.NewForConfig(rest)
		cm, err := clientset.CoreV1().ConfigMaps(nsName).Get(context.TODO(), cmName, v1.GetOptions{})
		if err != nil {
			return err
		}

		b, err := clientcmd.Write(*apiConfig)
		if err != nil {
			return err
		}
		cm.Data[fileName] = string(b)

		_, err = clientset.CoreV1().ConfigMaps(nsName).Update(context.TODO(), cm, v1.UpdateOptions{})
		if err != nil {
			return err
		}

		return nil
	}

	return conf, nil
}

// createClientSet - ClientSet 생성
func createClientSet(name string, restConfig *rest.Config) *ClientSet {
	cluster := &ClientSet{Name: name, RESTConfig: restConfig}

	// New kubernetes client
	cluster.NewKubernetesClient = func() (*kubernetes.Clientset, error) {
		return kubernetes.NewForConfig(restConfig)
	}

	// New discovery client
	cluster.NewDiscoveryClient = func() (*discovery.DiscoveryClient, error) {
		return discovery.NewDiscoveryClientForConfig(restConfig)
	}

	// New dynamic client
	cluster.NewDynamicClient = func() (*client.DynamicClient, error) {
		return client.NewDynamicClient(restConfig), nil
	}

	// ex. schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	cluster.NewDynamicClientSchema = func(group string, version string, resource string) (*client.DynamicClient, error) {
		return client.NewDynamicClientSchema(restConfig, group, version, resource), nil
	}

	return cluster
}

// newKubeCluster - Config 기반의 Cluster정보 생성
func newKubeCluster(conf *strategyInfo) (*kubeCluster, error) {
	var provider *clusterConfigProvider
	var err error

	// configmap provider 설정
	if conf.Strategy == KubeConfigStrategyConfigmap {
		provider, err = createConfigmapClusterConfig(conf.Data["namespace"], conf.Data["configmap"], conf.Data["filename"])
		if err != nil {
			logger.Warnf("can't create a 'configmap' kubeconfig-provider (cause-%s)", err.Error())
		}
	}

	// default provider
	if provider == nil {
		if provider, err = createFileClusterConfig(conf.Data["path"]); err != nil {
			logger.Warnf("can't create a 'file' kubeconfig-provider (cause-%s)", err.Error())
		}
	}

	if provider == nil {
		logger.Panicf("can't create a kubeconfig-provider")
	}

	clusters := make(map[string]*ClientSet)
	clusterNames := []string{}
	defaultContext := ""

	apiConfig, err := provider.load()
	if err != nil {
		logger.Warnf("can't load a kubeconfig (cause=%s)", err.Error())
		apiConfig = &api.Config{
			Clusters:  make(map[string]*api.Cluster),
			AuthInfos: make(map[string]*api.AuthInfo),
			Contexts:  make(map[string]*api.Context),
		}
	} else {
		if apiConfig.CurrentContext != "" {
			defaultContext = apiConfig.CurrentContext
		}

		for key := range apiConfig.Contexts {
			if defaultContext == "" {
				defaultContext = key
			}

			if restConfig, err := clientcmd.NewNonInteractiveClientConfig(*apiConfig, key, &clientcmd.ConfigOverrides{}, nil).ClientConfig(); err == nil {
				clusterNames = append(clusterNames, key)
				clusters[key] = createClientSet(key, restConfig)
			}
		}
	}

	// in-cluster
	inClusterConfig, _ := rest.InClusterConfig()

	// 로드된 컨텍스트가 없는 경우는 In-cluster 사용
	if len(clusters) == 0 && inClusterConfig != nil {
		clusters[IN_CLUSTER_NAME] = createClientSet(IN_CLUSTER_NAME, inClusterConfig)
		clusterNames = []string{IN_CLUSTER_NAME}
		defaultContext = IN_CLUSTER_NAME
	}

	// kubeCluster
	kc := &kubeCluster{
		KubeConfig:     apiConfig,
		clients:        clusters,
		ClusterNames:   clusterNames,
		DefaultContext: defaultContext,
	}

	if inClusterConfig != nil {
		kc.InCluster = createClientSet(IN_CLUSTER_NAME, inClusterConfig)
	}

	kc.IsRunningInCluster = (len(kc.ClusterNames) == 1 && kc.InCluster != nil) // running within a cluster

	// Get client by context
	kc.Client = func(context string) (*ClientSet, error) {
		val := utils.EndWithOnArray(kc.ClusterNames, context)
		if val == "" {
			//if !utils.ArrayContains(kc.ClusterNames, context) {
			if context == IN_CLUSTER_NAME && kc.InCluster != nil {
				return kc.InCluster, nil
			} else {
				return nil, fmt.Errorf("can't find a context '%s' in %v", context, kc.ClusterNames)
			}
		}
		return kc.clients[val], nil
	}

	// Save kubeconfig using provider (file or configmap)
	kc.Save = func() error {
		return provider.save(kc.KubeConfig)
	}

	// Add kubeconfig to managed cluster and save using provider
	kc.Add = func(confBytes []byte) error {
		config, err := parsingConfig(confBytes)
		if err != nil {
			return err
		}

		for k, v := range config.Clusters {
			if _, exists := kc.KubeConfig.Clusters[k]; exists {
				continue
			}

			kc.KubeConfig.Clusters[k] = v
		}

		for k, v := range config.Contexts {
			if _, exists := kc.KubeConfig.Contexts[k]; exists {
				continue
			}

			kc.KubeConfig.Contexts[k] = v
		}

		for k, v := range config.AuthInfos {
			if _, exists := kc.KubeConfig.AuthInfos[k]; exists {
				continue
			}

			kc.KubeConfig.AuthInfos[k] = v
		}

		return kc.Save()
	}

	// remove kubeconfig from managed cluster and save using provider
	kc.Remove = func(clusterName string) error {
		var modified bool = false
		for k, v := range kc.KubeConfig.Contexts {
			if v.Cluster == clusterName {
				delete(kc.KubeConfig.Clusters, v.Cluster)
				delete(kc.KubeConfig.AuthInfos, v.AuthInfo)
				delete(kc.KubeConfig.Contexts, k)

				modified = true
				break
			}
		}

		if modified {
			return kc.Save()
		}

		return nil
	}

	return kc, nil
}
