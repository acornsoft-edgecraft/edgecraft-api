-- edgecraft.tbl_cloud definition

-- Drop table

-- DROP TABLE edgecraft.tbl_cloud;

CREATE TABLE edgecraft.tbl_cloud (
	cloud_uid uuid NOT NULL DEFAULT uuid_generate_v4(),
	cloud_name varchar(50) NULL,
	cloud_type varchar(30) NULL,
	cloud_description varchar(300) NULL,
	cloud_state varchar(30) NULL,
	creator varchar(50) NULL,
	created_at timestamp NULL,
	updater varchar(50) NULL,
	updated_at timestamp NULL,
	CONSTRAINT "PK_tbl_cloud" PRIMARY KEY (cloud_uid)
);