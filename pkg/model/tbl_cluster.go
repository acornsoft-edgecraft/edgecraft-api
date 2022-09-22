/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"
)

// ClusterTable - Baremetal Cluster Table 정보
type ClusterTable struct {
	CloudUid   string `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid string `json:"cluster" db:"cluster_uid"`

	// K8S 정보
	Version string `json:"version" db:"k8s_version"`
	PodCidr string `json:"pod_cidr" db:"pod_cidr"`
	SvcCidr string `json:"svc_cidr" db:"service_cidr"`

	// Baremetal 정보
	BmcCredentialSecret   string       `json:"secret_name" db:"bmc_credential_secret"`
	BmcCredentialUser     string       `json:"user_name" db:"bmc_credential_user"`
	BmcCredentialPassword string       `json:"password" db:"bmc_credential_password"`
	ImageUrl              string       `json:"image_url" db:"image_url"`
	ImageChecksum         string       `json:"image_checksum" db:"image_checksum"`
	ImageChecksumType     string       `json:"image_checksum_type" db:"image_checksum_type"`
	ImageFormat           string       `json:"image_format" db:"image_format"`
	MasterExtraConfig     *ExtraConfig `json:"cp_kubeadm_extra_config" db:"master_extra_config"`
	WorkerExtraConfig     *ExtraConfig `json:"worker_kubeadm_extra_config" db:"worker_extra_config"`

	// nodes 정보
	LoadbalancerUse     bool   `json:"use_loadbalancer" db:"loadbalancer_use_yn, default:false"`
	LoadbalancerAddress string `json:"loadbalancer_address" db:"loadbalancer_address"`
	LoadbalancerPort    string `json:"loadbalancer_port" db:"loadbalancer_port"`

	// ETCD 정보
	ExternalEtcdUse bool `json:"use_external_etcd" db:"external_etcd_use, default:false"`
	//ExternalEtcdEndPoints       *ExternalEtcdEndPoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdEndPoints       *Endpoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdCertificateCa   string     `json:"ca_file" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert string     `json:"cert_file" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey  string     `json:"key_file" db:"external_etcd_certificate_key"`

	// Storage 정보
	StorageClass *StorageClass `json:"storage_class" db:"storage_class"`

	// 기본 정보
	Status string `json:"status" db:"state"`

	Creator string    `json:"creator" db:"creator"`
	Created time.Time `json:"created_at" db:"created_at"`
	Updater string    `json:"updater" db:"updater"`
	Updated time.Time `json:"updated_at" db:"updated_at"`
}

// MapToSet - Mapping cluster data to CloudSet
func (ct *ClusterTable) MapToSet(cloudSet *CloudSet) {
	// Cloud Info
	//cloudSet.Cluster = &ClusterInfo{}.MapToSet()
	// // Cluster Info
	// var clusterInfo *ClusterInfo = &ClusterInfo{}

	// // K8s Info
	// var k8sInfo *KubernetesInfo = &KubernetesInfo{}

	// // Bearmetal Info
	// var baremetalInfo *BaremetalInfo = &BaremetalInfo{}

	// // Etcd/Storage Info
	// var EtcdStorageInfo *EtcdStorageInfo = &EtcdStorageInfo{}

	// // Node Info
	// ClusterInfo.
	// utils.CopyTo(&cloudSet.Cluster, &ct)
}
