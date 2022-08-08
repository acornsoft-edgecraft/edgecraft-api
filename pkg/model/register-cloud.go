package model

// used pointer
type RegisterCloud struct {
	Cloud       Cloud          `json:"cloud"`
	Cluster     CloudCluster   `json:"cluster"`
	Nodes       CloudNodes     `json:"nodes"`
	EtcdStorage EtcdStorage    `json:"etcd_storage"`
	OpenStack   CloudOpenStack `json:"openstack"`
}
