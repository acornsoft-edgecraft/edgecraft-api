-- edgecraft.tbl_cloud_cluster definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud_cluster;

CREATE TABLE edgecraft.tbl_cloud_cluster (
	cloud_cluster_uid uuid NOT NULL DEFAULT uuid_generate_v4(),
	cloud_uid uuid NOT NULL,
	cloud_k8s_version varchar(100) NULL,
	cloud_cluster_bmc_credential_secret varchar(50) NULL,
	cloud_cluster_bmc_credential_user varchar(50) NULL,
	cloud_cluster_bmc_credential_password varchar(100) NULL,
	cloud_cluster_pod_cidr varchar(100) NULL,
	cloud_cluster_service_cidr varchar(100) NULL,
	cloud_cluster_loadbalancer_address varchar(50) NULL,
	ccloud_cluster_loadbalancer_port varchar(4) NULL,
	cloud_cluster_image_url varchar(255) NULL,
	cloud_cluster_image_checksum varchar(255) NULL,
	cloud_cluster_image_checksum_type varchar(30) NULL,
	cloud_cluster_image_format varchar(30) NULL,
	cloud_cluster_master_extra_config jsonb NULL,
	cloud_cluster_worker_extra_config jsonb NULL,
	cloud_cluster_storage_class jsonb NULL,
	cloud_cluster_state varchar(30) NULL,
	external_etcd_certificate_ca text NULL,
	external_etcd_certificate_cert text NULL,
	external_etcd_certificate_key text NULL,
	creator varchar(50) NULL,
	created_at timestamp NULL,
	updater varchar(50) NULL,
	updated_at timestamp NULL,
	cloud_cluster_loadbalancer_use bool NOT NULL,
	external_etcd_endpoints jsonb NULL,
	cloud_cluster_external_etcd_use bool NOT NULL,
	CONSTRAINT "PK_tbl_cloud_cluster" PRIMARY KEY (cloud_cluster_uid,cloud_uid)
);


-- edgecraft.tbl_cloud_cluster foreign keys

ALTER TABLE edgecraft.tbl_cloud_cluster ADD CONSTRAINT "FK_tbl_cloud_TO_tbl_cloud_cluster" FOREIGN KEY (cloud_uid) REFERENCES edgecraft.tbl_cloud(cloud_uid);