/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package kubemethod

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/model"
	batchV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const (
	NS = "edge-benchmarks"

	CM_Config  = "config"
	CM_Plugins = "plugins"

	role_Master = "master"
	role_Node   = "node"

	defaultBackoffLimit            = int32(0)
	defaultTTLSecondsAfterFinished = int32(0)
)

var (
	serviceName            = fmt.Sprintf("%s-service", NS)
	serviceAccountName     = fmt.Sprintf("%s-serviceaccount", NS)
	clusterRoleName        = fmt.Sprintf("%s-clusterrole", NS)
	clusterRoleBindingName = fmt.Sprintf("%s-clusterrolebinding", NS)
	configCmName           = fmt.Sprintf("%s-%s-cm", NS, CM_Config)
	pluginsCmName          = fmt.Sprintf("%s-%s-cm", NS, CM_Plugins)

	BenchmarksImage string
	SonobuoyImage   string
	SonobuoyVersion string
	Debug           string
)

func SetEdgeBenchmarks(clusterTable model.OpenstackClusterTable, benchmarksId string, conf *config.Benchmarks) error {
	BenchmarksImage = fmt.Sprintf("%s:%s", conf.Image, conf.Version)
	SonobuoyImage = conf.SonobuoyImage
	SonobuoyVersion = conf.SonobuoyVersion
	Debug = conf.Debug

	apiClient, err := config.HostCluster.GetKubernetesClient(*clusterTable.Name)
	if err != nil {
		return err
	}

	// Namespace
	err = setNs(apiClient)
	if err != nil {
		return err
	}

	// ServiceAccount
	err = setServiceAccount(apiClient)
	if err != nil {
		return err
	}

	// ClusterRole
	err = setClusterRole(apiClient)
	if err != nil {
		return err
	}

	// ClusterRoleBinding
	err = setClusterRoleBinding(apiClient)
	if err != nil {
		return err
	}

	// Service
	err = setService(apiClient)
	if err != nil {
		return err
	}

	// ConfigMap
	err = setConfigCM(apiClient, *clusterTable.BootstrapProvider)
	if err != nil {
		return err
	}
	err = setPluginsCM(apiClient, *clusterTable.BootstrapProvider)
	if err != nil {
		return err
	}

	// Job
	err = setJob(apiClient, benchmarksId)
	if err != nil {
		return err
	}

	return nil
}

func newNs() *coreV1.Namespace {
	namespace := &coreV1.Namespace{
		ObjectMeta: metaV1.ObjectMeta{
			Name: NS,
			Labels: labels.Set{
				"kubernetes.io/metadata.name": NS,
			},
		},
	}
	return namespace
}

func newServiceAccount() *coreV1.ServiceAccount {
	sa := &coreV1.ServiceAccount{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: NS,
			Name:      serviceAccountName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
			},
		},
	}
	return sa
}

func newClusterRole() *rbacV1.ClusterRole {
	cr := &rbacV1.ClusterRole{
		ObjectMeta: metaV1.ObjectMeta{
			Name: clusterRoleName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
			},
		},
		Rules: []rbacV1.PolicyRule{{
			APIGroups: []string{""},
			Resources: []string{"namespaces", "nodes", "pods", "serviceaccounts", "services", "configmaps", "secrets"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"secrets", "pods"},
			Verbs:     []string{"create", "update", "patch"},
		}, {
			APIGroups: []string{"apps"},
			Resources: []string{"daemonsets"},
			Verbs:     []string{"get", "list", "create", "update", "patch", "delete", "deletecollection"},
		}, {
			APIGroups: []string{"extensions"},
			Resources: []string{"daemonsets"},
			Verbs:     []string{"create", "update", "patch", "delete", "deletecollection"},
		}, {
			APIGroups: []string{"rbac.authorization.k8s.io"},
			Resources: []string{"rolebindings", "clusterrolebindings", "clusterroles"},
			Verbs:     []string{"get", "list"},
		}, {
			APIGroups: []string{"batch"},
			Resources: []string{"jobs"},
			Verbs:     []string{"list", "create", "patch", "update", "watch"},
		}, {
			NonResourceURLs: []string{"/metrics", "/logs", "/logs/*"},
			Verbs:           []string{"get"},
		}},
	}
	return cr
}

func newClusterRoleBinding() *rbacV1.ClusterRoleBinding {
	crb := &rbacV1.ClusterRoleBinding{
		ObjectMeta: metaV1.ObjectMeta{
			Name: clusterRoleBindingName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
			},
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
		},
		Subjects: []rbacV1.Subject{{
			Kind:      "ServiceAccount",
			Name:      serviceAccountName,
			Namespace: NS,
		}},
	}
	return crb
}

func newService() *coreV1.Service {
	service := &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: NS,
			Name:      serviceName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
				"sonobuoy-component":     "aggregator",
			},
		},
		Spec: coreV1.ServiceSpec{
			Type: coreV1.ServiceTypeClusterIP,
			Ports: []coreV1.ServicePort{{
				Port: 443,
				TargetPort: intstr.IntOrString{
					IntVal: 443,
				},
				Protocol: coreV1.ProtocolTCP,
			}},
			Selector: labels.Set{
				"app.kubernetes.io/name": NS,
				"sonobuoy-component":     "aggregator",
			},
		},
	}
	return service
}

func newConfigCM(data string) *coreV1.ConfigMap {
	cm := &coreV1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: NS,
			Name:      configCmName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
			},
		},
		Data: map[string]string{
			"config.json": data,
		},
	}
	return cm
}

func newPluginsCM(masterData, nodeData string) *coreV1.ConfigMap {
	cm := &coreV1.ConfigMap{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: NS,
			Name:      pluginsCmName,
			Labels: labels.Set{
				"app.kubernetes.io/name": NS,
			},
		},
		Data: map[string]string{
			"edge-kube-bench-master.yaml": masterData,
			"edge-kube-bench-node.yaml":   nodeData,
		},
	}
	return cm
}

func newJob(benchmarksId, name string) *batchV1.Job {
	var backoffLimit int32 = defaultBackoffLimit
	var ttlSecondsAfterFinished int32 = defaultTTLSecondsAfterFinished

	job := &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: NS,
			Name:      name,
		},
		Spec: batchV1.JobSpec{
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: labels.Set{
						"app.kubernetes.io/name": NS,
						"sonobuoy-component":     "aggregator",
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{{
						Name:            NS,
						Image:           BenchmarksImage,
						ImagePullPolicy: coreV1.PullIfNotPresent,
						Env: []coreV1.EnvVar{{
							Name:  "BENCHMARKS_ID",
							Value: benchmarksId,
						}, {
							Name:  "DEBUG",
							Value: Debug,
						}},
						VolumeMounts: []coreV1.VolumeMount{{
							Name:      "s-config-volume",
							MountPath: "/etc/sonobuoy",
						}, {
							Name:      "s-plugins-volume",
							MountPath: "/plugins.d",
						}, {
							Name:      "output-volume",
							MountPath: "/tmp/sonobuoy",
						}},
					}},
					ServiceAccountName: serviceAccountName,
					HostPID:            true,
					HostIPC:            true,
					RestartPolicy:      coreV1.RestartPolicyNever,
					Volumes: []coreV1.Volume{{
						Name: "s-config-volume",
						VolumeSource: coreV1.VolumeSource{
							ConfigMap: &coreV1.ConfigMapVolumeSource{
								LocalObjectReference: coreV1.LocalObjectReference{
									Name: configCmName,
								},
							},
						},
					}, {
						Name: "s-plugins-volume",
						VolumeSource: coreV1.VolumeSource{
							ConfigMap: &coreV1.ConfigMapVolumeSource{
								LocalObjectReference: coreV1.LocalObjectReference{
									Name: pluginsCmName,
								},
							},
						},
					}, {
						Name: "output-volume",
						VolumeSource: coreV1.VolumeSource{
							EmptyDir: &coreV1.EmptyDirVolumeSource{},
						},
					}},
				},
			},
		},
	}
	return job
}

func setNs(clientSet *kubernetes.Clientset) error {
	_, err := clientSet.CoreV1().Namespaces().Get(context.TODO(), NS, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		logger.Infof("Namespace \"%s\" not found", NS)

		namespace := newNs()
		_, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Infof("Namespace \"%s\" has created", NS)
		return nil
	} else if err != nil {
		return err
	}

	logger.Infof("Namespace \"%s\" already exists", NS)
	return nil
}

func setServiceAccount(clientSet *kubernetes.Clientset) error {
	_, err := clientSet.CoreV1().ServiceAccounts(NS).Get(context.TODO(), serviceAccountName, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		logger.Infof("ServiceAccount \"%s\" not found", serviceAccountName)

		sa := newServiceAccount()
		_, err := clientSet.CoreV1().ServiceAccounts(NS).Create(context.TODO(), sa, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Infof("ServiceAccount \"%s\" has created", serviceAccountName)
		return nil
	} else if err != nil {
		return err
	}

	logger.Infof("ServiceAccount \"%s\" already exists", serviceAccountName)
	return nil
}

func setService(clientSet *kubernetes.Clientset) error {
	_, err := clientSet.CoreV1().Services(NS).Get(context.TODO(), serviceName, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		logger.Infof("Service \"%s\" not found", serviceName)

		service := newService()
		_, err := clientSet.CoreV1().Services(NS).Create(context.TODO(), service, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Infof("Service \"%s\" has created", serviceName)
		return nil
	} else if err != nil {
		return err
	}

	logger.Infof("Service \"%s\" already exists", serviceName)
	return nil
}

func setClusterRole(clientSet *kubernetes.Clientset) error {
	_, err := clientSet.RbacV1().ClusterRoles().Get(context.TODO(), clusterRoleName, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		logger.Infof("ClusterRole \"%s\" not found", clusterRoleName)

		cr := newClusterRole()
		_, err := clientSet.RbacV1().ClusterRoles().Create(context.TODO(), cr, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Infof("ClusterRole \"%s\" has created", clusterRoleName)
		return nil
	} else if err != nil {
		return err
	}

	logger.Infof("ClusterRole \"%s\" already exists", clusterRoleName)
	return nil
}

func setClusterRoleBinding(clientSet *kubernetes.Clientset) error {
	_, err := clientSet.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metaV1.GetOptions{})
	if errors.IsNotFound(err) {
		logger.Infof("ClusterRoleBinding \"%s\" not found", clusterRoleBindingName)

		crb := newClusterRoleBinding()
		_, err := clientSet.RbacV1().ClusterRoleBindings().Create(context.TODO(), crb, metaV1.CreateOptions{})
		if err != nil {
			return err
		}
		logger.Infof("ClusterRoleBinding \"%s\" has created", clusterRoleBindingName)
		return nil
	} else if err != nil {
		return err
	}

	logger.Infof("ClusterRoleBinding \"%s\" already exists", clusterRoleBindingName)
	return nil
}

func setConfigCM(clientset *kubernetes.Clientset, bootstrapProvider common.BootstrapProvider) error {
	_, err := clientset.CoreV1().ConfigMaps(NS).Get(context.TODO(), configCmName, metaV1.GetOptions{})
	if err == nil {
		logger.Infof("ConfigMap \"%s\" already exists, delete ConfigMap", configCmName)

		e := clientset.CoreV1().ConfigMaps(NS).Delete(context.TODO(), configCmName, metaV1.DeleteOptions{})
		if e != nil {
			logger.Errorf("Configmap \"%s\" deletion failed. %v", configCmName, e)
		}
	}
	if errors.IsNotFound(err) {
		logger.Infof("ConfigMap \"%s\" not found. create ConfigMap", configCmName)
	} else if err != nil {
		return err
	}

	data, err := getTemplateParsing(fmt.Sprintf("./conf/templates/cis/%s-cm.data", CM_Config), getConfigCMData())
	if err != nil {
		return err
	}
	cm := newConfigCM(data)
	_, err = clientset.CoreV1().ConfigMaps(NS).Create(context.TODO(), cm, metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	logger.Infof("ConfigMap \"%s\" has created", configCmName)
	return nil
}

func setPluginsCM(clientset *kubernetes.Clientset, bootstrapProvider common.BootstrapProvider) error {
	_, err := clientset.CoreV1().ConfigMaps(NS).Get(context.TODO(), pluginsCmName, metaV1.GetOptions{})
	if err == nil {
		logger.Infof("ConfigMap \"%s\" already exists, delete ConfigMap", pluginsCmName)

		e := clientset.CoreV1().ConfigMaps(NS).Delete(context.TODO(), pluginsCmName, metaV1.DeleteOptions{})
		if e != nil {
			logger.Errorf("Configmap \"%s\" deletion failed. %v", pluginsCmName, e)
		}
	}

	if errors.IsNotFound(err) {
		logger.Infof("ConfigMap \"%s\" not found. create ConfigMap", pluginsCmName)
	} else if err != nil {
		return err
	}

	templatePath := fmt.Sprintf("./conf/templates/cis/%s-cm.data", CM_Plugins)
	masterData, err := getTemplateParsing(templatePath, getConfigPluginsData(role_Master, bootstrapProvider))
	if err != nil {
		return err
	}
	nodeData, err := getTemplateParsing(templatePath, getConfigPluginsData(role_Node, bootstrapProvider))
	if err != nil {
		return err
	}
	cm := newPluginsCM(masterData, nodeData)
	_, err = clientset.CoreV1().ConfigMaps(NS).Create(context.TODO(), cm, metaV1.CreateOptions{})
	if err != nil {
		return err
	}
	logger.Infof("ConfigMap \"%s\" has created", pluginsCmName)
	return nil
}

func setJob(clientSet *kubernetes.Clientset, benchmarksId string) error {
	t := time.Now()
	name := fmt.Sprintf("%s-%s", NS, t.Format("20060102150405"))
	job := newJob(benchmarksId, name)

	_, err := clientSet.BatchV1().Jobs(NS).Create(context.TODO(), job, metaV1.CreateOptions{})
	if err != nil {
		return err
	}

	logger.Infof("Job \"%s\" has created", name)
	return nil
}

func getConfigCMData() map[string]interface{} {
	cmd := map[string]interface{}{
		"namespace":          NS,
		"serviceAccountName": serviceAccountName,
		"advertiseAddress":   serviceName,
		"sonobuoyImage":      SonobuoyImage,
		"sonobuoyVersion":    SonobuoyVersion,
	}
	return cmd
}

func getConfigPluginsData(role string, bootstrapProvider common.BootstrapProvider) map[string]interface{} {
	cmd := map[string]interface{}{
		"pluginName":         fmt.Sprintf("edge-kube-bench-%s", role),
		"isMaster":           "true",
		"matchCpKey":         "node-role.kubernetes.io/master",
		"matchCpOperator":    "Exists",
		"serviceAccountName": serviceAccountName,
		"benchmarksImage":    BenchmarksImage,
	}

	if role == role_Node {
		cmd["isMaster"] = "false"
		cmd["matchCpOperator"] = "DoesNotExist"
	}
	if bootstrapProvider == common.K3s {
		cmd["benchmarkVersion"] = "k3s-1.23"
	} else if bootstrapProvider == common.MicroK8s {
		cmd["matchCpKey"] = "node.kubernetes.io/microk8s-controlplane"
	}

	return cmd
}

func getTemplateParsing(path string, data map[string]interface{}) (string, error) {
	temp, err := template.ParseFiles(path)
	if err != nil {
		logger.Errorf("Template has errors. cause(%s)", err.Error())
		return "", err
	}

	var buff bytes.Buffer
	err = temp.Execute(&buff, data)
	if err != nil {
		logger.Errorf("Template execution failed. cause(%s)", err.Error())
		return "", err
	}

	logger.Infof("processed path - %s templating data (%s)", path, buff.String())
	return buff.String(), nil
}
