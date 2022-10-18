/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import "github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"

// KubernetesInfo - Data for Kubernetes
type KubernetesInfo struct {
	Version   int    `json:"version" example:"1"`
	PodCidr   string `json:"pod_cidr" example:"10.96.0.1/12"`
	SvcCidr   string `json:"svc_cidr" example:"10.96.0.0/12"`
	SvcDomain string `json:"svc_domain" example:"cluster.local"`
}

// ToTable - K8S 정보르 테이블로 설정
func (ki *KubernetesInfo) ToTable(clusterTable *ClusterTable) {
	clusterTable.Version = utils.IntPrt(ki.Version)
	clusterTable.PodCidr = utils.StringPtr(ki.PodCidr)
	clusterTable.SvcCidr = utils.StringPtr(ki.SvcCidr)
	clusterTable.SvcDomain = utils.StringPtr(ki.SvcDomain)
}

// FromTable - 테이블 정보를 K8S로 설정
func (ki *KubernetesInfo) FromTable(clusterTable *ClusterTable) {
	ki.Version = *clusterTable.Version
	ki.PodCidr = *clusterTable.PodCidr
	ki.SvcCidr = *clusterTable.SvcCidr
	ki.SvcDomain = *clusterTable.SvcDomain
}

// ToOpenstackTable - K8S 정보 Openstack 테이블로 설정
func (ki *KubernetesInfo) ToOpenstackTable(clusterTable *OpenstackClusterTable) {
	clusterTable.Version = utils.IntPrt(ki.Version)
	clusterTable.PodCidr = utils.StringPtr(ki.PodCidr)
	clusterTable.SvcCidr = utils.StringPtr(ki.SvcCidr)
	clusterTable.SvcDomain = utils.StringPtr(ki.SvcDomain)
}

// FromOpenstackTable - Openstack 테이블 정보를 K8S로 설정
func (ki *KubernetesInfo) FromOpenstackTable(clusterTable *OpenstackClusterTable) {
	ki.Version = *clusterTable.Version
	ki.PodCidr = *clusterTable.PodCidr
	ki.SvcCidr = *clusterTable.SvcCidr
	ki.SvcDomain = *clusterTable.SvcDomain
}
