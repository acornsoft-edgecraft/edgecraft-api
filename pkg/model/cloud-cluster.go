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
	CloudClusterUid *uuid.UUID  `json:"cloudClusterUid" db:"cloud_cluster_uid, default:uuid_generate_v4()"`
	CloudUid        *uuid.UUID  `json:"cloudUid" db:"cloud_uid"`
	K8s             *K8s        `json:"k8s"`
	Baremetal       *Baremetal  `json:"baremetal"`
	CloudNodes      *CloudNodes `json:"nodes"`

	CloudClusterExternalEtcdUse     *string    `json:"cloudClusterExternalEtcdUse" db:"cloud_cluster_external_etcd_use"`
	CloudClusterStorageClassType    *string    `json:"cloudClusterStorageClassType" db:"cloud_cluster_storage_class_type"`
	CloudClusterState               *string    `json:"cloudClusterState" db:"cloud_cluster_state"`
	CloudClusterLoadbalancerAddress *string    `json:"loadbalancer_address" db:"cloud_cluster_loadbalancer_address"`
	CloudClusterLoadbalancerPort    *int8      `json:"loadbalancer_port" db:"ccloud_cluster_loadbalancer_port"`
	ExternalEtcdCertificateCa       *string    `json:"externalEtcdCertificateCa" db:"external_etcd_certificate_ca"`
	ExternalEtcdCertificateCert     *string    `json:"externalEtcdCertificateCert" db:"external_etcd_certificate_cert"`
	ExternalEtcdCertificateKey      *string    `json:"externalEtcdCertificateKey" db:"external_etcd_certificate_key"`
	Creator                         *string    `json:"creator" db:"creator"`
	CreatedAt                       *time.Time `json:"createdAt" db:"created_at"`
	Updater                         *string    `json:"updater" db:"updater"`
	UpdatedAt                       *time.Time `json:"updatedAt" db:"updated_at"`
}

type K8s struct {
	CloudK8sVersion         *string `json:"cloudK8sVersion" db:"cloud_k8s_version"`
	CloudClusterPodCidr     *string `json:"cloudClusterPodCidr" db:"cloud_cluster_pod_cidr"`
	CloudClusterServiceCidr *string `json:"cloudClusterServiceCidr" db:"cloud_cluster_service_cidr"`
}

type Baremetal struct {
	CloudClusterBmcCredentialSecret   *string `json:"secret_name" db:"cloud_cluster_bmc_credential_secret"`
	CloudClusterBmcCredentialUser     *string `json:"user_name" db:"cloud_cluster_bmc_credential_user"`
	CloudClusterBmcCredentialPassword *string `json:"password" db:"cloud_cluster_bmc_credential_password"`
	CloudClusterImageUrl              *string `json:"image_url" db:"cloud_cluster_image_url"`
	CloudClusterImageChecksum         *string `json:"image_checksum" db:"cloud_cluster_image_checksum"`
	CloudClusterImageChecksumType     *string `json:"image_checksum_type" db:"cloud_cluster_image_checksum_type"`
	CloudClusterImageFormat           *string `json:"image_format" db:"cloud_cluster_image_format"`
	CloudClusterMasterExtraConfig     *JSONB  `json:"cp_kubeadm_extra_config" db:"cloud_cluster_master_extra_config"`
	CloudClusterWorkerExtraConfig     *JSONB  `json:"worker_kubeadm_extra_config" db:"cloud_cluster_worker_extra_config"`
}

type EtcdStorage struct {
	Etcd         Etcd         `json:"etcd"`
	StorageClass StorageClass `json:"storage_class"`
}

//- Start - JSONB Interface for JSONB Field of yourTableName Table
type JSONB []interface{}

// Value Marshal
func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
