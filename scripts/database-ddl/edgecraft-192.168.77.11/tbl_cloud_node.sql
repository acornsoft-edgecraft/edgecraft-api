-- edgecraft.tbl_cloud_node definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud_node;

CREATE TABLE edgecraft.tbl_cloud_node (
	cloud_node_uid uuid NOT NULL DEFAULT uuid_generate_v4(),
	cloud_uid uuid NOT NULL,
	cloud_cluster_uid uuid NOT NULL,
	cloud_node_type varchar(30) NULL,
	cloud_node_state varchar(30) NULL,
	cloud_node_host_name varchar(50) NULL,
	cloud_node_name varchar(50) NULL,
	cloud_node_bmc_address varchar(255) NULL,
	cloud_node_mac_address varchar(255) NULL,
	cloud_node_ip varchar(50) NULL,
	cloud_node_label jsonb NULL,
	osd_path varchar(300) NULL,
	creator varchar(50) NULL,
	created_at timestamp NULL,
	updater varchar(50) NULL,
	updated_at timestamp NULL,
	cloud_node_boot_mode varchar(10) NULL,
	cloud_node_online_power bool NULL,
	cloud_node_external_provisioning bool NULL,
	CONSTRAINT "PK_tbl_cloud_node" PRIMARY KEY (cloud_node_uid, cloud_uid, cloud_cluster_uid)
);


-- edgecraft.tbl_cloud_node foreign keys

ALTER TABLE edgecraft.tbl_cloud_node ADD CONSTRAINT "FK_tbl_cloud_cluster_TO_tbl_cloud_node" FOREIGN KEY (cloud_cluster_uid,cloud_uid) REFERENCES edgecraft.tbl_cloud_cluster(cloud_cluster_uid,cloud_uid);