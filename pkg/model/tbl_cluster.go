/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
)

// ClusterTable - Baremetal Cluster Table 정보
// Status: Common Code Status 참조
type ClusterTable struct {
	CloudUid   *string `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid *string `json:"cluster" db:"cluster_uid"`

	// K8S 정보
	BootstrapProvider *common.BootstrapProvider `json:"bootstrap_provider" db:"bootstrap_provider"`
	Version           *int                      `json:"version" db:"k8s_version"`
	PodCidr           *string                   `json:"pod_cidr" db:"pod_cidr"`
	SvcCidr           *string                   `json:"svc_cidr" db:"service_cidr"`
	SvcDomain         *string                   `json:"svc_domain" db:"service_domain"`

	// Baremetal 정보
	Namespace             *string      `json:"namespace" db:"namespace"`
	BmcCredentialSecret   *string      `json:"secret_name" db:"bmc_credential_secret"`
	BmcCredentialUser     *string      `json:"user_name" db:"bmc_credential_user"`
	BmcCredentialPassword *string      `json:"password" db:"bmc_credential_password"`
	ImageUrl              *string      `json:"image_url" db:"image_url"`
	ImageChecksum         *string      `json:"image_checksum" db:"image_checksum"`
	ImageChecksumType     *int         `json:"image_checksum_type" db:"image_checksum_type"`
	ImageFormat           *int         `json:"image_format" db:"image_format"`
	MasterExtraConfig     *ExtraConfig `json:"cp_kubeadm_extra_config" db:"master_extra_config"`
	WorkerExtraConfig     *ExtraConfig `json:"worker_kubeadm_extra_config" db:"worker_extra_config"`

	// nodes 정보
	LoadbalancerUse     *bool   `json:"use_loadbalancer" db:"loadbalancer_use_yn, default:false"`
	LoadbalancerAddress *string `json:"loadbalancer_address" db:"loadbalancer_address"`
	LoadbalancerPort    *string `json:"loadbalancer_port" db:"loadbalancer_port"`

	// ETCD 정보
	ExternalEtcdUse             *bool      `json:"use_external_etcd" db:"external_etcd_use, default:false"`
	ExternalEtcdEndPoints       *Endpoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdCertificateCa   *string    `json:"ca_file" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert *string    `json:"cert_file" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey  *string    `json:"key_file" db:"external_etcd_certificate_key"`

	// Storage 정보
	StorageClass *StorageClass `json:"storage_class" db:"storage_class"`

	// 기본 정보
	Status  *int       `json:"status" db:"state"`
	Creator *string    `json:"creator" db:"creator"`
	Created *time.Time `json:"created_at" db:"created_at"`
	Updater *string    `json:"updater" db:"updater"`
	Updated *time.Time `json:"updated_at" db:"updated_at"`
}

// ToSet - Cluster 테이블 정보를 CloudSet 정보로 설정
func (ct *ClusterTable) ToSet(cloudSet *CloudSet) {
	// Cluster 정보 설정
	var cluster *ClusterInfo = &ClusterInfo{}
	cluster.FromTable(ct)

	// Etcd/Storage 정보 설정
	var etcdStorage *EtcdStorageInfo = &EtcdStorageInfo{}
	etcdStorage.FromTable(ct)

	cloudSet.Cluster = cluster
	cloudSet.EtcdStorage = etcdStorage
	cloudSet.OpenStack = &OpenstackInfo{}
}
