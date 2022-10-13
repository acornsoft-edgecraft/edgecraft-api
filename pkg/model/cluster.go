/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// ClusterInfo - Data for Cluster
// Status: Common Code Status 참조
type ClusterInfo struct {
	ClusterUid string          `json:"cluster_uid" example:""`
	Status     int             `json:"status"`
	K8s        *KubernetesInfo `json:"k8s"`
	Baremetal  *BaremetalInfo  `json:"baremetal"`
}

// NewKey - Make new UUID V4
func (ci *ClusterInfo) NewKey() {
	if ci.ClusterUid == "" {
		ci.ClusterUid = uuid.Must(uuid.NewV4()).String()
	}
}

// ToTable - Cluster 정보를 테이블로 설정
func (ci *ClusterInfo) ToTable(clusterTable *ClusterTable, isUpdate bool, user string, at time.Time) {
	if isUpdate {
		if clusterTable.ClusterUid == nil {
			clusterTable.ClusterUid = utils.StringPtr(ci.ClusterUid)
		}
		clusterTable.Updater = utils.StringPtr(user)
		clusterTable.Updated = utils.TimePtr(at)
	} else {
		ci.NewKey()
		clusterTable.ClusterUid = utils.StringPtr(ci.ClusterUid)
		clusterTable.Creator = utils.StringPtr(user)
		clusterTable.Created = utils.TimePtr(at)
	}

	clusterTable.Status = utils.IntPrt(1)
	ci.K8s.ToTable(clusterTable)
	ci.Baremetal.ToTable(clusterTable)

}

// FromTable - 테이블에서 Cluster로 정보 설정
func (ci *ClusterInfo) FromTable(clusterTable *ClusterTable) {
	ci.ClusterUid = *clusterTable.ClusterUid
	ci.Status = *clusterTable.Status

	ci.K8s = &KubernetesInfo{}
	ci.Baremetal = &BaremetalInfo{}

	ci.K8s.FromTable(clusterTable)
	ci.Baremetal.FromTable(clusterTable)
}

// EtcdStorageInfo - Data for ETCD and Storage
type EtcdStorageInfo struct {
	Etcd         *EtcdInfo         `json:"etcd"`
	StorageClass *StorageClassInfo `json:"storage_class"`
}

// ToTable - ETCD/Storage 정보를 테이블로 설정
func (esi *EtcdStorageInfo) ToTable(clusterTable *ClusterTable) {
	esi.Etcd.ToTable(clusterTable)
	esi.StorageClass.ToTable(clusterTable)
}

// FromTable - 테이블 정보를 ETCD/Storage 정보 설정
func (esi *EtcdStorageInfo) FromTable(clusterTable *ClusterTable) {
	esi.Etcd = &EtcdInfo{}
	esi.StorageClass = &StorageClassInfo{}

	esi.Etcd.FromTable(clusterTable)
	esi.StorageClass.FromTable(clusterTable)
}

// ToTable - ETCD/Storage 정보를 Openstack 테이블로 설정
func (esi *EtcdStorageInfo) ToOpenstackTable(clusterTable *OpenstackClusterTable) {
	esi.Etcd.ToOpenstackTable(clusterTable)
	esi.StorageClass.ToOpenstackTable(clusterTable)
}

// FromTable - Openstack 테이블 정보를 ETCD/Storage 정보 설정
func (esi *EtcdStorageInfo) FromOpenstackTable(clusterTable *OpenstackClusterTable) {
	esi.Etcd = &EtcdInfo{}
	esi.StorageClass = &StorageClassInfo{}

	esi.Etcd.FromOpenstackTable(clusterTable)
	esi.StorageClass.FromOpenstackTable(clusterTable)
}
