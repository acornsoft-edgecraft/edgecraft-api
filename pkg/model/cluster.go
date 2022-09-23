/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
	"github.com/gofrs/uuid"
)

// ClusterInfo - Data for Cluster
type ClusterInfo struct {
	ClusterUid string          `json:"cluster_uid" example:""`
	Status     string          `json:"status"`
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
func (ci *ClusterInfo) ToTable(clusterTable *ClusterTable) {
	ci.NewKey()
	utils.CopyTo(&clusterTable, ci.K8s)
	utils.CopyTo(&clusterTable, ci.Baremetal)

	clusterTable.ClusterUid = ci.ClusterUid
}

// FromTable - 테이블에서 Cluster로 정보 설정
func (ci *ClusterInfo) FromTable(clusterTable *ClusterTable) {
	ci.ClusterUid = clusterTable.ClusterUid
	ci.Status = clusterTable.Status

	ci.K8s = &KubernetesInfo{}
	ci.Baremetal = &BaremetalInfo{}

	ci.K8s.FromTable(clusterTable)
	ci.Baremetal.FromTable(clusterTable)
}

// EtcdStorageInfo - Data for ETCD and Storage
type EtcdStorageInfo struct {
	Etcd         *EtcdInfo         `json:"etcd"`
	StorageClass *StorageClassInfo `json:"storage_class" db:"cloud_cluster_storage_class"`
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

// type K8s struct {
// 	CloudUid                *uuid.UUID `json:"-" db:"cloud_uid"`
// 	CloudClusterUid         *uuid.UUID `json:"-" db:"cloud_cluster_uid"`
// 	CloudK8sVersion         *string    `json:"version" db:"cloud_k8s_version"`
// 	CloudClusterPodCidr     *string    `json:"pod_cidr" db:"cloud_cluster_pod_cidr"`
// 	CloudClusterServiceCidr *string    `json:"svc_cidr" db:"cloud_cluster_service_cidr"`
// }

// type CloudCluster struct {
// 	ClusterUid                  *uuid.UUID             `json:"cloud_cluster_uid" db:"cloud_cluster_uid, default:uuid_generate_v4()"`
// 	CloudUid                    *uuid.UUID             `json:"cloud_uid" db:"cloud_uid"`
// 	CloudK8sVersion             *string                `json:"version" db:"cloud_k8s_version"`
// 	PodCidr                     *string                `json:"pod_cidr" db:"cloud_cluster_pod_cidr"`
// 	ServiceCidr                 *string                `json:"svc_cidr" db:"cloud_cluster_service_cidr"`
// 	BmcCredentialSecret         *string                `json:"secret_name" db:"cloud_cluster_bmc_credential_secret"`
// 	BmcCredentialUser           *string                `json:"user_name" db:"cloud_cluster_bmc_credential_user"`
// 	BmcCredentialPassword       *string                `json:"password" db:"cloud_cluster_bmc_credential_password"`
// 	ImageUrl                    *string                `json:"image_url" db:"cloud_cluster_image_url"`
// 	ImageChecksum               *string                `json:"image_checksum" db:"cloud_cluster_image_checksum"`
// 	ImageChecksumType           *string                `json:"image_checksum_type" db:"cloud_cluster_image_checksum_type"`
// 	ImageFormat                 *string                `json:"image_format" db:"cloud_cluster_image_format"`
// 	MasterExtraConfig           *ExtraConfig           `json:"cp_kubeadm_extra_config" db:"cloud_cluster_master_extra_config"`
// 	WorkerExtraConfig           *ExtraConfig           `json:"worker_kubeadm_extra_config" db:"cloud_cluster_worker_extra_config"`
// 	LoadbalancerUse             *bool                  `json:"use_loadbalancer" db:"cloud_cluster_loadbalancer_use, default:false"`
// 	LoadbalancerAddress         *string                `json:"loadbalancer_address" db:"cloud_cluster_loadbalancer_address"`
// 	LoadbalancerPort            *string                `json:"loadbalancer_port" db:"ccloud_cluster_loadbalancer_port"`
// 	ExternalEtcdUse             *bool                  `json:"use_external_etcd" db:"cloud_cluster_external_etcd_use, default:false"`
// 	ExternalEtcdEndPoints       *ExternalEtcdEndPoints `json:"endpoints" db:"external_etcd_endpoints"`
// 	ExternalEtcdCertificateCa   *string                `json:"ca_file" db:"external_etcd_certificate_ca"`
// 	ExternalEtcdCertificateCert *string                `json:"cert_file" db:"external_etcd_certificate_cert"`
// 	ExternalEtcdCertificateKey  *string                `json:"key_file" db:"external_etcd_certificate_key"`
// 	StorageClass                *StorageClass          `json:"storage_class" db:"cloud_cluster_storage_class"`
// 	State                       *string                `json:"cloudClusterState" db:"cloud_cluster_state"`
// 	CUInfo                      CreateUpdateInfo
// 	// Creator                     *string    `json:"creator" db:"creator"`
// 	// CreatedAt                   *time.Time `json:"createdAt" db:"created_at"`
// 	// Updater                     *string    `json:"updater" db:"updater"`
// 	// UpdatedAt                   *time.Time `json:"updatedAt" db:"updated_at"`
// }

// type Etcd struct {
// 	CloudUid                    *uuid.UUID             `json:"-" db:"cloud_uid"`
// 	CloudClusterUid             *uuid.UUID             `json:"-" db:"cloud_cluster_uid"`
// 	CloudClusterExternalEtcdUse *bool                  `json:"use_external_etcd" db:"cloud_cluster_external_etcd_use" default:"false"`
// 	ExternalEtcdEndPoints       *ExternalEtcdEndPoints `json:"endpoints" db:"external_etcd_endpoints"`
// 	ExternalEtcdCertificateCa   *string                `json:"ca_file" db:"external_etcd_certificate_ca"`
// 	ExternalEtcdCertificateCert *string                `json:"cert_file" db:"external_etcd_certificate_cert"`
// 	ExternalEtcdCertificateKey  *string                `json:"key_file" db:"external_etcd_certificate_key"`
// }

// type StorageClass struct {
// 	Use_ceph bool    `json:"use_ceph" db:"-, default:false"`
// 	Labels   []Label `json:"labels"`
// }

// // Value Marshal
// func (a StorageClass) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Scan Unmarshal
// func (a *StorageClass) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}
// 	return json.Unmarshal(b, &a)
// }

// type ExternalEtcdEndPoints []interface{}

// // Value Marshal
// func (a ExternalEtcdEndPoints) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Scan Unmarshal
// func (a *ExternalEtcdEndPoints) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}
// 	return json.Unmarshal(b, &a)
// }

// //region [ ExtraConfig ]

// // Value Marshal
// func (a ExtraConfig) Value() (driver.Value, error) {
// 	return json.Marshal(a)
// }

// // Scan Unmarshal
// func (a *ExtraConfig) Scan(value interface{}) error {
// 	b, ok := value.([]byte)
// 	if !ok {
// 		return errors.New("type assertion to []byte failed")
// 	}
// 	return json.Unmarshal(b, &a)
// }

// //endregion [ ExtraConfig ]
