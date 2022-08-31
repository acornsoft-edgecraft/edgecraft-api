package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type CloudNode struct {
	CloudNodeUid                  *uuid.UUID `json:"cloud_node_uid" db:"cloud_node_uid, default:uuid_generate_v4()"`
	CloudUid                      *uuid.UUID `json:"cloud_uid" db:"cloud_uid"`
	CloudClusterUid               *uuid.UUID `json:"cloud_cluster_uid" db:"cloud_cluster_uid"`
	CloudNodeType                 *string    `json:"cloud_node_type" db:"cloud_node_type"`
	CloudNodeState                *string    `json:"cloud_node_state" db:"cloud_node_state"`
	CloudNodeHostName             *string    `json:"host_name" db:"cloud_node_host_name"`
	CloudNodeBmcAddress           *string    `json:"bmc_address" db:"cloud_node_bmc_address"`
	CloudNodeMacAddress           *string    `json:"boot_mac_address" db:"cloud_node_mac_address"`
	CloudNodeBootMode             *string    `json:"boot_mode" db:"cloud_node_boot_mode"`
	CloudNodeOnlinePower          *bool      `json:"online_power" db:"cloud_node_online_power"`
	CloudNodeExternalProvisioning *bool      `json:"external_provisioning" db:"cloud_node_external_provisioning"`
	CloudNodeName                 *string    `json:"node_name" db:"cloud_node_name"`
	CloudNodeIp                   *string    `json:"ip_address" db:"cloud_node_ip"`
	CloudNodeLabel                *Labels    `json:"labels" db:"cloud_node_label"`
	OsdPath                       *string    `json:"osd_path" db:"osd_path"`
	Creator                       *string    `json:"creator" db:"creator"`
	CreatedAt                     *time.Time `json:"created_at" db:"created_at"`
	Updater                       *string    `json:"updater" db:"updater"`
	UpdatedAt                     *time.Time `json:"updated_at" db:"updated_at"`
}

type MasterNode struct {
	CloudNodeUid *uuid.UUID `json:"cloud_node_uid" db:"cloud_node_uid"`
	// CloudUid        *uuid.UUID    `json:"-" db:"cloud_uid"`
	// CloudClusterUid *uuid.UUID    `json:"-" db:"cloud_cluster_uid"`
	Baremetal NodeBaremetal `json:"baremetal"`
	Node      Nodes         `json:"node" db:"-"`
}
type WorkerNode struct {
	CloudNodeUid *uuid.UUID `json:"cloud_node_uid" db:"cloud_node_uid"`
	// CloudUid        *uuid.UUID    `json:"-" db:"cloud_uid"`
	// CloudClusterUid *uuid.UUID    `json:"-" db:"cloud_cluster_uid"`
	Baremetal NodeBaremetal `json:"baremetal" db:"-"`
	Node      Nodes         `json:"node" db:"-"`
}
type NodeBaremetal struct {
	CloudNodeHostName             *string `json:"host_name" db:"cloud_node_host_name"`
	CloudNodeBmcAddress           *string `json:"bmc_address" db:"cloud_node_bmc_address"`
	CloudNodeMacAddress           *string `json:"boot_mac_address" db:"cloud_node_mac_address"`
	CloudNodeBootMode             *string `json:"boot_mode" db:"cloud_node_boot_mode"`
	CloudNodeOnlinePower          *bool   `json:"online_power" db:"cloud_node_online_power"`
	CloudNodeExternalProvisioning *bool   `json:"external_provisioning" db:"cloud_node_external_provisioning"`
}

type Nodes struct {
	CloudNodeName  *string `json:"node_name" db:"cloud_node_name"`
	CloudNodeIp    *string `json:"ip_address" db:"cloud_node_ip"`
	CloudNodeLabel *Labels `json:"labels" db:"cloud_node_label"`
}

type Label struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type Labels []interface{}

// Value Marshal
func (a Labels) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Labels) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
