/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

// KubernetesInfo - Data for Kubernetes
type KubernetesInfo struct {
	Version string `json:"version" example:"1"`
	PodCidr string `json:"pod_cidr" example:"10.244.0.0/16"`
	SvcCidr string `json:"svc_cidr" example:"10.244.0.0/16"`
}

// ToTable - K8S 정보르 테이블로 설정
func (ki *KubernetesInfo) ToTable(clusterTable *ClusterTable) {
	*clusterTable.Version = ki.Version
	*clusterTable.PodCidr = ki.PodCidr
	*clusterTable.SvcCidr = ki.SvcCidr
}

// FromTable - 테이블 정보를 K8S로 설정
func (ki *KubernetesInfo) FromTable(clusterTable *ClusterTable) {
	ki.Version = *clusterTable.Version
	ki.PodCidr = *clusterTable.PodCidr
	ki.SvcCidr = *clusterTable.SvcCidr
}
