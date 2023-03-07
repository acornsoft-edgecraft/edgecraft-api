/*
Copyright 2023 Acornsoft Authors. All right reserved.
*/
package model

// ===== [ Constants and Variables ] =====

// K8sUpgradeInfo - Kubernetes Version Upgrade Info for openstack cluster
type K8sUpgradeInfo struct {
	Version int    `json:"version"`
	Image   string `json:"image"`
}
