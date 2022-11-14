/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/utils"
)

// BaremetalInfo - Data for Baremetal
type BaremetalInfo struct {
	Namespace             string `json:"namespace" example:"default"`
	BmcCredentialSecret   string `json:"secret_name" example:"secret1"`
	BmcCredentialUser     string `json:"user_name" example:"asdf"`
	BmcCredentialPassword string `json:"password" example:"asdf"`
	ImageUrl              string `json:"image_url" example:"http://192.168.0.1/ubuntu.qcow2"`
	ImageChecksum         string `json:"image_checksum" example:"http://192.168.0.1/ubuntu.qcow2.md5sum"`
	ImageChecksumType     int    `json:"image_checksum_type" example:"1"`
	ImageFormat           int    `json:"image_format" example:"2"`
}

// ToTable - Baremetal 정보를 테이블로 설정
func (bi *BaremetalInfo) ToTable(clusterTable *ClusterTable) {
	clusterTable.Namespace = utils.StringPtr(bi.Namespace)
	clusterTable.BmcCredentialSecret = utils.StringPtr(bi.BmcCredentialSecret)
	clusterTable.BmcCredentialUser = utils.StringPtr(bi.BmcCredentialUser)
	clusterTable.BmcCredentialPassword = utils.StringPtr(bi.BmcCredentialPassword)
	clusterTable.ImageUrl = utils.StringPtr(bi.ImageUrl)
	clusterTable.ImageChecksum = utils.StringPtr(bi.ImageChecksum)
	clusterTable.ImageChecksumType = utils.IntPrt(bi.ImageChecksumType)
	clusterTable.ImageFormat = utils.IntPrt(bi.ImageFormat)
}

// FromTable - 테이블의 정보를 Baremetal 정보로 설정
func (bi *BaremetalInfo) FromTable(clusterTable *ClusterTable) {
	bi.Namespace = *clusterTable.Namespace
	bi.BmcCredentialSecret = *clusterTable.BmcCredentialSecret
	bi.BmcCredentialUser = *clusterTable.BmcCredentialUser
	bi.BmcCredentialPassword = *clusterTable.BmcCredentialPassword
	bi.ImageUrl = *clusterTable.ImageUrl
	bi.ImageChecksum = *clusterTable.ImageChecksum
	bi.ImageChecksumType = *clusterTable.ImageChecksumType
	bi.ImageFormat = *clusterTable.ImageFormat
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
