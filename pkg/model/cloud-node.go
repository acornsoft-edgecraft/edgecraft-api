package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type CloudNodes struct {
	CloudClusterNodeUid  *uuid.UUID `jsong:"" db:"cloud_cluster_node_uid"`
	CloudUid             *uuid.UUID `jsong:"" db:"cloud_uid"`
	CloudClusterUid      *uuid.UUID `jsong:"" db:"cloud_cluster_uid"`
	CloudName            *string    `jsong:"" db:"cloud_name"`
	CloudNodeType        *string    `jsong:"" db:"cloud_node_type"`
	CloudNodeState       *string    `jsong:"" db:"cloud_node_state"`
	CloudNodeHostName    *string    `jsong:"" db:"cloud_node_host_name"`
	CloudNodeName        *string    `jsong:"" db:"cloud_node_name"`
	CloudNodeBmcAddress  *string    `jsong:"" db:"cloud_node_bmc_address"`
	CloudNodeMacAddress  *string    `jsong:"" db:"cloud_node_mac_address"`
	CloudNodeIp          *string    `jsong:"" db:"cloud_node_ip"`
	CloudNodeServiceType *string    `jsong:"" db:"cloud_node_service_type"`
	CloudNodeLabel       *string    `jsong:"" db:"cloud_node_label"`
	OsdPath              *string    `jsong:"" db:"osd_path"`
	Creator              *string    `jsong:"" db:"creator"`
	CreatedAt            *time.Time `jsong:"" db:"created_at"`
	Updater              *string    `jsong:"" db:"updater"`
	UpdatedAt            *time.Time `jsong:"" db:"updated_at"`
}
