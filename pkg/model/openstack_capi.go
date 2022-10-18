/*
Copyright 2022 Acornsoft Authors. All right reserved.
*/
package model

// OpenstackCAPI - CAPI Provider for openstack
type OpenstackCAPI struct {
	OpenstackCloud                     string   `json:"openstack_cloud"`
	ClusterName                        string   `json:"name"`           // 화면 연계? (Cluster Name)
	Namespace                          string   `json:"namespace"`      // 화면 연계? (NodeSet 정보, Master/Worker 다른 경우는?) "default"
	LocalHostName                      string   `json:"local_hostname"` // "{{ local_hostname }}" # self binding
	KubernetesVersion                  string   `json:"version"`
	PodCidr                            string   `json:"pod_cidr"`
	NodeCidr                           string   `json:"node_cidr"` // 화면 연계 필요
	ServiceCidr                        string   `json:"svc_cidr"`
	ServiceDomain                      string   `json:"svc_domain"` // 화면 연계 필요, "cluster.local"
	OpenstackCloudCACertB64            string   `json:"openstack_cloud_cacert_b64"`
	OpenstackCloudYamlB64              string   `json:"openstack_cloud_yaml_b64"`
	OpenstackClusterProviderConfB64    string   `json:"openstack_cloud_provider_conf_b64"`
	OpenstackFailureDomain             string   `json:"failure_domain"` // "nova"
	OpenstackDNSNameServers            string   `json:"dns_nameservers"`
	OpenstackExternalNetworkID         string   `json:"external_network_id"` // "public"
	OpenstackControlPlaneMachineFlavor string   `json:"master_flavor"`       // 화면 연계? (MasterSet)
	OpenstackNodeMahineFlavor          string   `json:"worker_flavor"`       // 화면 연계? (WorkerSet)
	OpenstackImageName                 string   `json:"image_name"`
	OpenstackSSHKeyName                string   `json:"ssh_key_name"`
	APIServerFloatingIP                string   `json:"api_server_floating_ip"`
	UseExternalETCD                    bool     `json:"use_external_etcd"`
	ETCDEndpoints                      []string `json:"endpoints"`
	ETCDCAFile                         string   `json:"ca_file"`
	ETCDCertFile                       string   `json:"cert_file"`
	ETCDKeyFile                        string   `json:"key_file"`
	UseBastionHost                     bool     `json:"use_bastion_host"`
	BastionFlavor                      string   `json:"bastion_flavor"`
	BastionImageName                   string   `json:"bastion_image_name"`
	BastionSSHKeyName                  string   `json:"bastion_ssh_key_name"`
	BastionFloatingIP                  string   `json:"bastion_floating_ip"`
	ControlPlaneMachineCount           int      `json:"master_node_count"` // 화면 연계? (MasterSet)
	WorkerMachineCount                 int      `json:"worker_node_count"` // 화면 연계? (WorkerSet)
}

// UI 연계 누락 정보
// NodeCidr

// 사용여부
// StorageClass
