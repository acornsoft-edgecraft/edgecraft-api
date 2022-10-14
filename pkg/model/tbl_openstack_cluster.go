/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import "time"

// OpenstackClusterTable - 클러스터 테이블 정보 (Openstack)
type OpenstackClusterTable struct {
	CloudUid      *string        `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid    *string        `json:"cluster" db:"cluster_uid"`
	ProjectName   *string        `json:"project_name" db:"project_name"`
	Name          *string        `json:"name" db:"name"`
	Version       *int           `json:"version" db:"version"`
	Desc          *string        `json:"desc" db:"description"`
	Credential    *string        `json:"credential" db:"credential"`
	PodCidr       *string        `json:"pod_cidr" db:"pod_cidr"`
	SvcCidr       *string        `json:"svc_cidr" db:"service_cidr"`
	OpenstackInfo *OpenstackInfo `json:"openstack_info" db:"openstack_info"`

	// nodes 정보
	LoadbalancerUse *bool `json:"use_loadbalancer" db:"loadbalancer_use_yn, default:false"`

	// ETCD 정보
	ExternalEtcdUse             *bool      `json:"use_external_etcd" db:"external_etcd_use, default:false"`
	ExternalEtcdEndPoints       *Endpoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdCertificateCa   *string    `json:"ca_file" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert *string    `json:"cert_file" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey  *string    `json:"key_file" db:"external_etcd_certificate_key"`

	// Storage 정보
	StorageClass *StorageClass `json:"storage_class" db:"storage_class"`
	// StorageUserId     *string        `json:"storage_user_id" db:"storage_user_id"`
	// StorageUserSecret *string        `json:"storage_user_secret" db:"storage_user_secret"`

	Status  *int       `json:"status" db:"state"`
	Creator *string    `json:"creator" db:"creator"`
	Created *time.Time `json:"created_at" db:"created_at"`
	Updater *string    `json:"updater" db:"updater"`
	Updated *time.Time `json:"updated_at" db:"updated_at"`
}
