/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import (
	"time"
)

// NodeTable - Baremetal Node Table 정보
type NodeTable struct {
	CloudUid   string `json:"cloud_uid" db:"cloud_uid"`
	ClusterUid string `json:"cluster_uid" db:"cluster_uid"`
	NodeUid    string `json:"node_uid" db:"node_uid, primarykey"`

	// Baremetal Host 정보
	HostName             string `json:"host_name" db:"host_name"`
	BmcAddress           string `json:"bmc_address" db:"bmc_address"`
	MacAddress           string `json:"boot_mac_address" db:"mac_address"`
	BootMode             string `json:"boot_mode" db:"boot_mode"`
	OnlinePower          bool   `json:"online_power" db:"online_power"`
	ExternalProvisioning bool   `json:"external_provisioning" db:"external_provisioning"`

	// Node 정보
	Name      string  `json:"node_name" db:"name"`
	Ipaddress string  `json:"ip_address" db:"ipaddress"`
	Labels    *Labels `json:"labels" db:"label"`

	// Openstack Ceph Path
	// TODO: (? - 화면에 없음, 향후 조정 필요)
	OsdPath string `json:"osd_path" db:"osd_path"`

	Type    string    `json:"type" db:"type"`
	Status  string    `json:"status" db:"state"`
	Creator string    `json:"creator" db:"creator"`
	Created time.Time `json:"created" db:"created_at"`
	Updater string    `json:"updater" db:"updater"`
	Updated time.Time `json:"updated" db:"updated_at"`
}

// MapToSet - Mapping node data to CloudSet
func (nt *NodeTable) MapToSet(cloudSet *CloudSet) {

}
