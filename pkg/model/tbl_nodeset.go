/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

import "time"

// NodesetTable - Openstack NodeSet Table 정보
// Type: Common Code NodeTypes 참조
type NodesetTable struct {
	ClusterUid *string    `json:"cluster_uid" db:"cluster_uid"`
	NodeUid    *string    `json:"nodeset_uid" db:"nodeset_uid"`
	Type       *int       `json:"type" db:"type"`
	Namespace  *string    `json:"namespace" db:"namespace"`
	Name       *string    `json:"name" db:"name"`
	NodeCount  *int       `json:"node_count" db:"node_count"`
	Flavor     string     `json:"flavor" db:"flavor"`
	Labels     string     `json:"labels" db:"labels"`
	Creator    *string    `json:"creator" db:"creator"`
	Created    *time.Time `json:"created" db:"created_at"`
	Updater    *string    `json:"updater" db:"updater"`
	Updated    *time.Time `json:"updated" db:"updated_at"`
}
