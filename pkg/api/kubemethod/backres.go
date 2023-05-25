/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package kubemethod

import (
	"errors"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	group             string = "velero.io"
	version           string = "v1"
	backup_resources  string = "backups"
	restore_resources string = "restores"
)

// GetBackResStatusPhase - 지정한 클러스터에 대한 Backup / Restore Status/Phase 검증
func GetBackResStatusPhase(clusterName, namespace, backresName string, isBackup bool) (string, error) {
	var resources string

	if isBackup {
		resources = backup_resources
	} else {
		resources = restore_resources
	}

	// Get kubernetes client
	dynamicClient, err := config.HostCluster.GetDynamicClientWithSchema(clusterName, group, version, resources)
	if err != nil {
		return "", err
	}

	// checking the cluster ready
	dynamicClient.SetNamespace(namespace)
	data, err := dynamicClient.Get(backresName, metaV1.GetOptions{})
	if err != nil {
		return "", err
	} else if data != nil {
		items, found, err := unstructured.NestedSlice(data.Object, "status", "phase")
		if err != nil {
			return "", err
		} else if !found {
			return "", errors.New("Not found [status/phase] fields in CR [" + data.GetName() + "] Object")
		}

		phase := items[0].(string)
		return string(phase[0]), nil
	}
	return "", errors.New("Unknown [status/phase] fields in CR [" + data.GetName() + "] Object")
}
