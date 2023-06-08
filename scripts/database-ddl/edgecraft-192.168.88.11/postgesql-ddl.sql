-- DROP SCHEMA edgecraft;

CREATE SCHEMA edgecraft AUTHORIZATION edgecraft;
-- edgecraft.tbl_cloud definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud;

CREATE TABLE edgecraft.tbl_cloud (
	cloud_uid bpchar(36) NOT NULL,
	"name" varchar(30) NOT NULL,
	"type" int4 NOT NULL DEFAULT 1,
	description varchar(300) NULL,
	state int4 NOT NULL DEFAULT 1,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_cloud" PRIMARY KEY (cloud_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cloud OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cloud TO edgecraft;


-- edgecraft.tbl_cloud_cluster definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud_cluster;

CREATE TABLE edgecraft.tbl_cloud_cluster (
	cloud_uid bpchar(36) NOT NULL,
	cluster_uid bpchar(36) NOT NULL,
	k8s_version int4 NOT NULL DEFAULT 1,
	pod_cidr varchar(30) NULL,
	service_cidr varchar(30) NULL,
	bmc_credential_secret varchar(50) NULL,
	bmc_credential_user varchar(50) NULL,
	bmc_credential_password varchar(100) NULL,
	image_url varchar(255) NULL,
	image_checksum varchar(255) NULL,
	image_checksum_type int4 NULL,
	image_format int4 NULL,
	master_extra_config json NULL,
	worker_extra_config json NULL,
	loadbalancer_use_yn bool NULL DEFAULT false,
	loadbalancer_address varchar(30) NULL,
	loadbalancer_port varchar(6) NULL,
	external_etcd_use bool NULL DEFAULT false,
	external_etcd_endpoints json NULL,
	external_etcd_certificate_ca text NULL,
	external_etcd_certificate_cert text NULL,
	external_etcd_certificate_key text NULL,
	storage_class json NULL,
	state int4 NOT NULL DEFAULT 1,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	service_domain varchar(30) NULL,
	"namespace" varchar(100) NULL DEFAULT '''default'''::character varying,
	bootstrap_provider int4 NOT NULL DEFAULT 1,
	CONSTRAINT "PK_tbl_cloud_cluster" PRIMARY KEY (cluster_uid, cloud_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cloud_cluster OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cloud_cluster TO edgecraft;


-- edgecraft.tbl_cloud_node definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud_node;

CREATE TABLE edgecraft.tbl_cloud_node (
	cloud_uid bpchar(36) NOT NULL,
	cluster_uid bpchar(36) NOT NULL,
	node_uid bpchar(36) NOT NULL,
	host_name varchar(30) NULL,
	bmc_address varchar(255) NULL,
	mac_address varchar(255) NULL,
	boot_mode int4 NULL,
	online_power bool NULL,
	external_provisioning bool NULL,
	"name" varchar(30) NULL,
	ipaddress varchar(30) NULL,
	labels json NULL,
	osd_path varchar(300) NULL,
	"type" int4 NOT NULL DEFAULT 1,
	state int4 NOT NULL DEFAULT 1,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_cloud_node" PRIMARY KEY (node_uid, cloud_uid, cluster_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cloud_node OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cloud_node TO edgecraft;


-- edgecraft.tbl_cluster definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cluster;

CREATE TABLE edgecraft.tbl_cluster (
	cloud_uid bpchar(36) NOT NULL,
	cluster_uid bpchar(36) NOT NULL,
	project_name varchar(50) NULL,
	"name" varchar(50) NULL,
	"version" int4 NOT NULL DEFAULT 1,
	description varchar(100) NULL,
	credential text NULL,
	pod_cidr varchar(100) NULL,
	service_cidr varchar(100) NULL,
	openstack_info json NULL,
	loadbalancer_use_yn bool NULL DEFAULT false,
	external_etcd_use bool NULL DEFAULT false,
	external_etcd_endpoints json NULL,
	external_etcd_certificate_ca text NULL,
	external_etcd_certificate_cert text NULL,
	external_etcd_certificate_key text NULL,
	storage_class json NULL,
	state int4 NOT NULL DEFAULT 1,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	service_domain varchar(30) NULL,
	"namespace" varchar(100) NULL DEFAULT '''default'''::character varying,
	master_extra_config json NULL,
	worker_extra_config json NULL,
	bootstrap_provider int4 NOT NULL DEFAULT 1,
	CONSTRAINT "PK_tbl_cluster" PRIMARY KEY (cluster_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cluster OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cluster TO edgecraft;


-- edgecraft.tbl_cluster_backres definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cluster_backres;

CREATE TABLE edgecraft.tbl_cluster_backres (
	cloud_uid bpchar(36) NOT NULL,
	cluster_uid bpchar(36) NOT NULL,
	backres_uid bpchar(36) NOT NULL,
	"name" varchar(200) NOT NULL,
	"type" bpchar(1) NOT NULL DEFAULT 'B'::bpchar,
	status bpchar(1) NOT NULL DEFAULT 'R'::bpchar,
	reason varchar(200) NULL,
	backup_name varchar(200) NULL,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_cluster_backres" PRIMARY KEY (backres_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cluster_backres OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cluster_backres TO edgecraft;


-- edgecraft.tbl_cluster_benchmarks definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cluster_benchmarks;

CREATE TABLE edgecraft.tbl_cluster_benchmarks (
	cluster_uid bpchar(36) NOT NULL,
	benchmarks_uid bpchar(36) NOT NULL,
	cis_version varchar(30) NULL,
	detected_version varchar(30) NULL,
	state int4 NULL DEFAULT 1,
	results json NULL,
	totals json NULL,
	reason varchar(200) NULL,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_cluster_benchmarks" PRIMARY KEY (cluster_uid, benchmarks_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_cluster_benchmarks OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_cluster_benchmarks TO edgecraft;


-- edgecraft.tbl_code definition

-- Drop table

-- DROP TABLE edgecraft.tbl_code;

CREATE TABLE edgecraft.tbl_code (
	group_id varchar(30) NOT NULL,
	code int4 NOT NULL,
	"name" varchar(50) NOT NULL,
	display_order int4 NULL,
	description varchar(300) NULL,
	use_yn bool NOT NULL DEFAULT true,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_code" PRIMARY KEY (group_id, code)
);

-- Permissions

ALTER TABLE edgecraft.tbl_code OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_code TO edgecraft;


-- edgecraft.tbl_code_group definition

-- Drop table

-- DROP TABLE edgecraft.tbl_code_group;

CREATE TABLE edgecraft.tbl_code_group (
	group_id varchar(30) NOT NULL,
	description varchar(3000) NULL,
	use_yn bool NULL DEFAULT true,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_code_group" PRIMARY KEY (group_id)
);

-- Permissions

ALTER TABLE edgecraft.tbl_code_group OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_code_group TO edgecraft;


-- edgecraft.tbl_nodeset definition

-- Drop table

-- DROP TABLE edgecraft.tbl_nodeset;

CREATE TABLE edgecraft.tbl_nodeset (
	cluster_uid bpchar(36) NOT NULL,
	nodeset_uid bpchar(36) NOT NULL,
	"type" int4 NOT NULL DEFAULT 1,
	"name" varchar(50) NULL,
	node_count int4 NULL DEFAULT 0,
	flavor varchar(50) NULL,
	labels json NULL,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NOT NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_nodeset" PRIMARY KEY (cluster_uid, nodeset_uid)
);

-- Permissions

ALTER TABLE edgecraft.tbl_nodeset OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_nodeset TO edgecraft;


-- edgecraft.tbl_user definition

-- Drop table

-- DROP TABLE edgecraft.tbl_user;

CREATE TABLE edgecraft.tbl_user (
	user_uid bpchar(36) NOT NULL,
	"role" int4 NOT NULL,
	"name" varchar(50) NOT NULL,
	id varchar(30) NULL,
	"password" varchar(128) NOT NULL,
	email varchar(50) NOT NULL,
	last_login timestamp NULL,
	password_expiration_begin_time timestamp NULL,
	password_expiration_end_time timestamp NULL,
	reset_password_yn bool NULL DEFAULT false,
	active_datetime timestamp NULL,
	inactive_yn bool NULL DEFAULT false,
	state int4 NULL DEFAULT 1,
	creator varchar(30) NOT NULL DEFAULT 'system'::character varying,
	created_at timestamp NULL DEFAULT now(),
	updater varchar(30) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_user" PRIMARY KEY (user_uid),
	CONSTRAINT "UK_tbl_user" UNIQUE (id, email)
);

-- Permissions

ALTER TABLE edgecraft.tbl_user OWNER TO edgecraft;
GRANT ALL ON TABLE edgecraft.tbl_user TO edgecraft;



CREATE OR REPLACE FUNCTION edgecraft.uuid_generate_v1()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v1$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_generate_v1() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_generate_v1() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_generate_v1mc()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v1mc$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_generate_v1mc() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_generate_v1mc() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_generate_v3(namespace uuid, name text)
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v3$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_generate_v3(uuid, text) OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_generate_v3(uuid, text) TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_generate_v4()
 RETURNS uuid
 LANGUAGE c
 PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v4$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_generate_v4() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_generate_v4() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_generate_v5(namespace uuid, name text)
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v5$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_generate_v5(uuid, text) OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_generate_v5(uuid, text) TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_nil()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_nil$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_nil() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_nil() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_ns_dns()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_dns$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_ns_dns() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_ns_dns() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_ns_oid()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_oid$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_ns_oid() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_ns_oid() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_ns_url()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_url$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_ns_url() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_ns_url() TO postgres;

CREATE OR REPLACE FUNCTION edgecraft.uuid_ns_x500()
 RETURNS uuid
 LANGUAGE c
 IMMUTABLE PARALLEL SAFE STRICT
AS '$libdir/uuid-ossp', $function$uuid_ns_x500$function$
;

-- Permissions

ALTER FUNCTION edgecraft.uuid_ns_x500() OWNER TO postgres;
GRANT ALL ON FUNCTION edgecraft.uuid_ns_x500() TO postgres;


-- Permissions

GRANT ALL ON SCHEMA edgecraft TO edgecraft;


-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="Asia/Seoul";
