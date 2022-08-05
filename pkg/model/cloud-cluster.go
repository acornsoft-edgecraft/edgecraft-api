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
	Cloudclusteruid                   *uuid.UUID `json:"cloudClusterUid" db:"cloud_cluster_uid, default:uuid_generate_v4()"`
	Clouduid                          *uuid.UUID `json:"cloudUid" db:"cloud_uid"`
	Cloudk8sversion                   *string    `json:"cloudK8sVersion" db:"cloud_k8s_version"`
	Cloudclusterbmccredentialsecret   *string    `json:"cloudClusterBmcCredentialSecret" db:"cloud_cluster_bmc_credential_secret"`
	Cloudclusterbmccredentialuser     *string    `json:"cloudClusterBmcCredentialUser" db:"cloud_cluster_bmc_credential_user"`
	Cloudclusterbmccredentialpassword *string    `json:"cloudClusterBmcCredentialPassword" db:"cloud_cluster_bmc_credential_password"`
	Cloudclusterpodcidr               *string    `json:"cloudClusterPodCidr" db:"cloud_cluster_pod_cidr"`
	Cloudclusterendpoint              *string    `json:"cloudClusterEndpoint" db:"cloud_cluster_service_cidr"`
	Cloudclusterendpointport          *int8      `json:"cloudClusterEndpointPort" db:"cloud_cluster_endpoint"`
	Cloudclusterimageurl              *string    `json:"cloudClusterImageUrl" db:"cloud_cluster_endpoint_port"`
	Cloudclusterimagechecksum         *string    `json:"cloudClusterImageChecksum" db:"cloud_cluster_image_url"`
	Cloudclusterimagechecksumtype     *string    `json:"cloudClusterImageChecksumType" db:"cloud_cluster_image_checksum"`
	Cloudclusterimageformat           *string    `json:"cloudClusterImageFormat" db:"cloud_cluster_image_checksum_type"`
	Cloudclustermasterextraconfig     *JSONB     `json:"cloudClusterMasterExtraConfig" db:"cloud_cluster_image_format"`
	Cloudclusterworkerextraconfig     *JSONB     `json:"cloudClusterWorkerExtraConfig" db:"cloud_cluster_master_extra_config"`
	Cloudclusterexternaletcduse       *string    `json:"cloudClusterExternalEtcdUse" db:"cloud_cluster_worker_extra_config"`
	Cloudclusterstorageclasstype      *string    `json:"cloudClusterStorageClassType" db:"cloud_cluster_external_etcd_use"`
	Cloudclusterstate                 *string    `json:"cloudClusterState" db:"cloud_cluster_storage_class_type"`
	Externaletcdcertificateca         *string    `json:"externalEtcdCertificateCa" db:"cloud_cluster_state"`
	Externaletcdcertificatecert       *string    `json:"externalEtcdCertificateCert" db:"external_etcd_certificate_ca"`
	Externaletcdcertificatekey        *string    `json:"externalEtcdCertificateKey" db:"external_etcd_certificate_cert"`
	Creator                           *string    `json:"creator" db:"creator"`
	CreatedAt                         *time.Time `json:"createdAt" db:"created_at"`
	Updater                           *string    `json:"updater" db:"updater"`
	UpdatedAt                         *time.Time `json:"updatedAt" db:"updated_at"`
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
