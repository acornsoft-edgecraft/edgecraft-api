package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

// used pointer
type CloudCluster struct {
	CloudClusterUid                   *uuid.UUID             `json:"cloudClusterUid" db:"cloud_cluster_uid, default:uuid_generate_v4()"`
	CloudUid                          *uuid.UUID             `json:"cloudUid" db:"cloud_uid"`
	CloudK8sVersion                   *string                `json:"cloudK8sVersion" db:"cloud_k8s_version"`
	CloudClusterPodCidr               *string                `json:"cloudClusterPodCidr" db:"cloud_cluster_pod_cidr"`
	CloudClusterServiceCidr           *string                `json:"cloudClusterServiceCidr" db:"cloud_cluster_service_cidr"`
	CloudClusterBmcCredentialSecret   *string                `json:"secret_name" db:"cloud_cluster_bmc_credential_secret"`
	CloudClusterBmcCredentialUser     *string                `json:"user_name" db:"cloud_cluster_bmc_credential_user"`
	CloudClusterBmcCredentialPassword *string                `json:"password" db:"cloud_cluster_bmc_credential_password"`
	CloudClusterImageUrl              *string                `json:"image_url" db:"cloud_cluster_image_url"`
	CloudClusterImageChecksum         *string                `json:"image_checksum" db:"cloud_cluster_image_checksum"`
	CloudClusterImageChecksumType     *string                `json:"image_checksum_type" db:"cloud_cluster_image_checksum_type"`
	CloudClusterImageFormat           *string                `json:"image_format" db:"cloud_cluster_image_format"`
	CloudClusterMasterExtraConfig     *ExtraConfig           `json:"cp_kubeadm_extra_config" db:"cloud_cluster_master_extra_config"`
	CloudClusterWorkerExtraConfig     *ExtraConfig           `json:"worker_kubeadm_extra_config" db:"cloud_cluster_worker_extra_config"`
	CloudClusterLoadbalancerUse       *bool                  `json:"use_loadbalancer" db:"cloud_cluster_loadbalancer_use, default:false"`
	CloudClusterLoadbalancerAddress   *string                `json:"loadbalancer_address" db:"cloud_cluster_loadbalancer_address"`
	CloudClusterLoadbalancerPort      *string                `json:"loadbalancer_port" db:"ccloud_cluster_loadbalancer_port"`
	CloudClusterExternalEtcdUse       *bool                  `json:"cloudClusterExternalEtcdUse" db:"cloud_cluster_external_etcd_use"`
	ExternalEtcdEndPoints             *ExternalEtcdEndPoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdCertificateCa         *string                `json:"externalEtcdCertificateCa" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert       *string                `json:"externalEtcdCertificateCert" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey        *string                `json:"externalEtcdCertificateKey" db:"external_etcd_certificate_key"`
	CloudClusterStorageClass          *StorageClass          `json:"storage_class" db:"cloud_cluster_storage_class"`
	CloudClusterState                 *string                `json:"cloudClusterState" db:"cloud_cluster_state"`
	Creator                           *string                `json:"creator" db:"creator"`
	CreatedAt                         *time.Time             `json:"createdAt" db:"created_at"`
	Updater                           *string                `json:"updater" db:"updater"`
	UpdatedAt                         *time.Time             `json:"updatedAt" db:"updated_at"`
}

// type CloudCluster struct {
// 	CloudClusterUid          *uuid.UUID   `json:"cloudClusterUid" db:"cloud_cluster_uid, default:uuid_generate_v4()"`
// 	CloudUid                 *uuid.UUID   `json:"cloudUid" db:"cloud_uid"`
// 	K8s                      K8s          `json:"k8s"`
// 	Baremetal                Baremetal    `json:"baremetal" db:"-"`
// 	Nodes                    Nodes        `json:"nodes" db:"-"`
// 	Etcd                     Etcd         `json:"etcd" db:"-"`
// 	CloudClusterStorageClass StorageClass `json:"storage_class" db:"cloud_cluster_storage_class"`
// 	CloudClusterState        *string      `json:"cloudClusterState" db:"cloud_cluster_state"`
// 	Creator                  *string      `json:"creator" db:"creator"`
// 	CreatedAt                *time.Time   `json:"createdAt" db:"created_at"`
// 	Updater                  *string      `json:"updater" db:"updater"`
// 	UpdatedAt                *time.Time   `json:"updatedAt" db:"updated_at"`
// }

type K8s struct {
	CloudUid                *uuid.UUID `json:"-" db:"cloud_uid"`
	CloudClusterUid         *uuid.UUID `json:"-" db:"cloud_cluster_uid"`
	CloudK8sVersion         *string    `json:"version" db:"cloud_k8s_version"`
	CloudClusterPodCidr     *string    `json:"pod_cidr" db:"cloud_cluster_pod_cidr"`
	CloudClusterServiceCidr *string    `json:"svc_cidr" db:"cloud_cluster_service_cidr"`
}

type ClusterBaremetal struct {
	CloudUid                          *uuid.UUID   `json:"-" db:"cloud_uid"`
	CloudClusterUid                   *uuid.UUID   `json:"-" db:"cloud_cluster_uid"`
	CloudClusterBmcCredentialSecret   *string      `json:"secret_name" db:"cloud_cluster_bmc_credential_secret"`
	CloudClusterBmcCredentialUser     *string      `json:"user_name" db:"cloud_cluster_bmc_credential_user"`
	CloudClusterBmcCredentialPassword *string      `json:"password" db:"cloud_cluster_bmc_credential_password"`
	CloudClusterImageUrl              *string      `json:"image_url" db:"cloud_cluster_image_url"`
	CloudClusterImageChecksum         *string      `json:"image_checksum" db:"cloud_cluster_image_checksum"`
	CloudClusterImageChecksumType     *string      `json:"image_checksum_type" db:"cloud_cluster_image_checksum_type"`
	CloudClusterImageFormat           *string      `json:"image_format" db:"cloud_cluster_image_format"`
	CloudClusterMasterExtraConfig     *ExtraConfig `json:"cp_kubeadm_extra_config" db:"cloud_cluster_master_extra_config"`
	CloudClusterWorkerExtraConfig     *ExtraConfig `json:"worker_kubeadm_extra_config" db:"cloud_cluster_worker_extra_config"`
}

type ClusterNodes struct {
	CloudUid                        *uuid.UUID `json:"-" db:"cloud_uid"`
	CloudClusterUid                 *uuid.UUID `json:"-" db:"cloud_cluster_uid"`
	CloudClusterLoadbalancerUse     *bool      `json:"use_loadbalancer" db:"cloud_cluster_loadbalancer_use"`
	CloudClusterLoadbalancerAddress *string    `json:"loadbalancer_address" db:"cloud_cluster_loadbalancer_address"`
	CloudClusterLoadbalancerPort    *string    `json:"loadbalancer_port" db:"ccloud_cluster_loadbalancer_port"`
}

type EtcdStorage struct {
	Etcd                     Etcd          `json:"etcd"`
	CloudClusterStorageClass *StorageClass `json:"storage_class" db:"cloud_cluster_storage_class"`
}
type Etcd struct {
	CloudUid                    *uuid.UUID             `json:"-" db:"cloud_uid"`
	CloudClusterUid             *uuid.UUID             `json:"-" db:"cloud_cluster_uid"`
	CloudClusterExternalEtcdUse *bool                  `json:"use_external_etcd" db:"cloud_cluster_external_etcd_use"`
	ExternalEtcdEndPoints       *ExternalEtcdEndPoints `json:"endpoints" db:"external_etcd_endpoints"`
	ExternalEtcdCertificateCa   *string                `json:"ca_file" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert *string                `json:"cert_file" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey  *string                `json:"key_file" db:"external_etcd_certificate_key"`
}

type StorageClass struct {
	Use_ceph bool    `json:"use_ceph" db:"-, default:false"`
	Labels   []Label `json:"label"`
}

// Value Marshal
func (a StorageClass) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *StorageClass) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type ExternalEtcdEndPoints []interface{}

// Value Marshal
func (a ExternalEtcdEndPoints) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *ExternalEtcdEndPoints) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

//- Start - JSONB Interface for JSONB Field of yourTableName Table
type ExtraConfig struct {
	PreKubeadmCommands  interface{} `json:"pre_kubeadm_commands"`
	PostKubeadmCommands interface{} `json:"post_kubeadm_commands"`
	Files               interface{} `json:"files"`
	Users               interface{} `json:"users"`
	Ntp                 interface{} `json:"ntp"`
	Format              interface{} `json:"format"`
}

// Value Marshal
func (a ExtraConfig) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *ExtraConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
