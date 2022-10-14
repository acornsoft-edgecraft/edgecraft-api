/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

// BaremetalInfo - Data for Baremetal
type BaremetalInfo struct {
	BmcCredentialSecret   string       `json:"secret_name" example:"secret1"`
	BmcCredentialUser     string       `json:"user_name" example:"asdf"`
	BmcCredentialPassword string       `json:"password" example:"asdf"`
	ImageUrl              string       `json:"image_url" example:"http://192.168.0.1/ubuntu.qcow2"`
	ImageChecksum         string       `json:"image_checksum" example:"http://192.168.0.1/ubuntu.qcow2.md5sum"`
	ImageChecksumType     int          `json:"image_checksum_type" example:"1"`
	ImageFormat           int          `json:"image_format" example:"2"`
	MasterExtraConfig     *ExtraConfig `json:"cp_kubeadm_extra_config"`
	WorkerExtraConfig     *ExtraConfig `json:"worker_kubeadm_extra_config"`
}

// ToTable - Baremetal 정보를 테이블로 설정
func (bi *BaremetalInfo) ToTable(clusterTable *ClusterTable) {
	clusterTable.BmcCredentialSecret = utils.StringPtr(bi.BmcCredentialSecret)
	clusterTable.BmcCredentialUser = utils.StringPtr(bi.BmcCredentialUser)
	clusterTable.BmcCredentialPassword = utils.StringPtr(bi.BmcCredentialPassword)
	clusterTable.ImageUrl = utils.StringPtr(bi.ImageUrl)
	clusterTable.ImageChecksum = utils.StringPtr(bi.ImageChecksum)
	clusterTable.ImageChecksumType = utils.IntPrt(bi.ImageChecksumType)
	clusterTable.ImageFormat = utils.IntPrt(bi.ImageFormat)

	clusterTable.MasterExtraConfig = &ExtraConfig{}
	clusterTable.WorkerExtraConfig = &ExtraConfig{}

	bi.MasterExtraConfig.ToTable(clusterTable.MasterExtraConfig)
	bi.WorkerExtraConfig.ToTable(clusterTable.WorkerExtraConfig)
}

// FromTable - 테이블의 정보를 Baremetal 정보로 설정
func (bi *BaremetalInfo) FromTable(clusterTable *ClusterTable) {
	bi.BmcCredentialSecret = *clusterTable.BmcCredentialSecret
	bi.BmcCredentialUser = *clusterTable.BmcCredentialUser
	bi.BmcCredentialPassword = *clusterTable.BmcCredentialPassword
	bi.ImageUrl = *clusterTable.ImageUrl
	bi.ImageChecksum = *clusterTable.ImageChecksum
	bi.ImageChecksumType = *clusterTable.ImageChecksumType
	bi.ImageFormat = *clusterTable.ImageFormat

	bi.MasterExtraConfig = &ExtraConfig{}
	bi.WorkerExtraConfig = &ExtraConfig{}

	bi.MasterExtraConfig.FromTable(clusterTable.MasterExtraConfig)
	bi.WorkerExtraConfig.FromTable(clusterTable.WorkerExtraConfig)
}

// BaremetalHostInfo - Data for Bearmetal host
type BaremetalHostInfo struct {
	HostName             string `json:"host_name" example:"sadf"`
	BmcAddress           string `json:"bmc_address" example:"98:03:9b:61:80:48"`
	BootMacAddress       string `json:"boot_mac_address" example:"00:b2:8c:ee:22:98"`
	BootMode             int    `json:"boot_mode" example:"1"`
	OnlinePower          bool   `json:"online_power" example:"false"`
	ExternalProvisioning bool   `json:"external_provisioning" example:"false"`
}

// ToTable - Baremetal Host 정보를 테이블로 설정
func (bhi *BaremetalHostInfo) ToTable(nodeTable *NodeTable) {
	nodeTable.HostName = utils.StringPtr(bhi.HostName)
	nodeTable.BmcAddress = utils.StringPtr(bhi.BmcAddress)
	nodeTable.MacAddress = utils.StringPtr(bhi.BootMacAddress)
	nodeTable.BootMode = utils.IntPrt(bhi.BootMode)
	nodeTable.OnlinePower = utils.BoolPtr(bhi.OnlinePower)
	nodeTable.ExternalProvisioning = utils.BoolPtr(bhi.ExternalProvisioning)
}

// FromTable - 테이블 정보를 Baremetal Host 정보로 설정
func (bhi *BaremetalHostInfo) FromTable(nodeTable *NodeTable) {
	bhi.HostName = *nodeTable.HostName
	bhi.BmcAddress = *nodeTable.BmcAddress
	bhi.BootMacAddress = *nodeTable.MacAddress
	bhi.BootMode = *nodeTable.BootMode
	bhi.OnlinePower = *nodeTable.OnlinePower
	bhi.ExternalProvisioning = *nodeTable.ExternalProvisioning
}

// ExtraConfig - Extra configuration for kubeadm
type ExtraConfig struct {
	PreKubeadmCommands  string `json:"pre_kubeadm_commands" example:"a"`
	PostKubeadmCommands string `json:"post_kubeadm_commands" example:"b"`
	Files               string `json:"files" example:"c"`
	Users               string `json:"users" example:"d"`
	Ntp                 string `json:"ntp" example:"e"`
	Format              string `json:"format" example:"f"`
}

// ToTable - ExtraConfig 정보를 테이블 정보로 설정
func (ec *ExtraConfig) ToTable(config *ExtraConfig) {
	config.PreKubeadmCommands = ec.PreKubeadmCommands
	config.PostKubeadmCommands = ec.PostKubeadmCommands
	config.Files = ec.Files
	config.Users = ec.Users
	config.Ntp = ec.Ntp
	config.Format = ec.Format
}

// FromTable - 테이블의 정보를 ExtraConfig 정보로 설정
func (ec *ExtraConfig) FromTable(config *ExtraConfig) {
	ec.PreKubeadmCommands = config.PreKubeadmCommands
	ec.PostKubeadmCommands = config.PostKubeadmCommands
	ec.Files = config.Files
	ec.Users = config.Users
	ec.Ntp = config.Ntp
	ec.Format = config.Format
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
