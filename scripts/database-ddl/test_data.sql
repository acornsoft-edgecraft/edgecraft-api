-- User
INSERT INTO "edgecraft"."tbl_user"
(
    user_uid,
    role,
    name,
    id,
    password,
    email,
    password_expiration_begin_time,
    password_expiration_end_time
)
SELECT 
	uuid_generate_v4(), 
	1,
	'morris', 
	'ccambo', 
	'662edc7bc97f9425be9b3abd61f9eea63fbff1cf144a9a572bab168145fc8838ac2f9f92e284f22a974e3820f947a909190d686e4fe9476316ca7d36fd9965d7', 
	'ccambo@acornsoft.io',
	NOW(),
	NOW()
-- pwd >> 1234abcd@Acorn

-- Code Group
INSERT INTO "edgecraft"."tbl_code_group"
(
	"group_id", 
    "description"
)
VALUES 
('CloudTypes',          'Types of Cloud'),
('CloudStatus',         'Status of cloud'),
('K8sVersions',         'Version of kubernetes for cloud'),
('ImageChecksumTypes',  'Types of checkSum for image'),
('ImageFormats',        'Types of formats for image'),
('BootModes',           'Modes of boot'),
('NodeTypes',           'Types of Node'),
('UserRoles',           'Tyeps of Role for user'),
('ImageTypes',          'Types of image'),
('ImageOsTypes',        'Types of image for OS'),
('SecurityStatus',      'Status of security'),
('SecurityItemStatus',  'Status of items for security'),
('UserStatus',          'Status of user'),
('ClusterStatus',       'Status of cluster'),
('NodeState',           'Status of node')

-- Code
INSERT INTO "edgecraft"."tbl_code"
(
	"group_id",
	"code",
	"name",
	"display_order",
    "description"
)
VALUES
('CloudTypes', 			1, 		'Baremetal', 		0,  'Baremetal Cloud'),
('CloudTypes', 			2, 		'Openstack', 		1,  'Openstsack Cloud'),
('CloudStatus', 		1, 		'Saved', 			0,  'Saved status'),
('CloudStatus', 		2, 		'Provisioning',		1,  'Provisioning status'),
('CloudStatus', 		3, 		'Provisioned', 		2,  'Provisioned status'),
('CloudStatus', 		4, 		'Failed',			3,  'Failed status'),
('CloudStatus', 		5, 		'Deleting', 		4,  'Deleting status'),
('CloudStatus', 		6, 		'Deleted', 			5,  'Deleted status'),
('K8sVersions', 		1, 		'1.22.0',			0,  'K8s version 1.22.0'),
('K8sVersions', 		2, 		'1.22.3',			0,  'K8s version 1.22.3'),
('ImageChecksumTypes',	1,		'md5',				0,  'MD5 checking for image'),
('ImageChecksumTypes',	2,		'sha256',			1,  'SHA256 checking for image'),
('ImageChecksumTypes',	3,		'sha512',			2,  'SHA512 checking for image'),
('ImageFormats',		1,		'raw',				0,  'RAW format for image'),
('ImageFormats',		2,		'qcow2',			1,  'QCOW2 format for image'),
('ImageFormats',		3,		'vdi',				2,  'VDI format for image'),
('ImageFormats',		4,		'vmdk',				3,  'VMDK format for image'),
('ImageFormats',		5,		'live-iso',			4,  'LIVE-ISO format for image'),
('BootModes',			1,		'UEFI',				0,  'UEFI boot'),
('BootModes',			2,		'legacy',			1,  'LEGACY boot'),
('BootModes',			3,		'UEFISecureBoot',	2,  'UEFI Security boot'),
('NodeTypes',			1,		'Master',			0,  'Master Node'),
('NodeTypes',			2,		'Worker',			1,  'Worker Node'),
('UserRoles',			1,		'Admin',			0,  'Admin user'),
('UserRoles',			2,		'Manager',			1,  'Manage user'),
('UserRoles',			3,		'User',				2,  'Normal user'),
('ImageTypes',          1,      'Baremetal Cloud',  0,  'Baremetal cloud image'),
('ImageTypes',          2,      'Openstack Cloud',  1,  'Openstack cloud image'),
('ImageOsTypes',        1,      'Ubuntu 20.04',     0,  'Ubuntu v20.04'),
('ImageOsTypes',        2,      'Ubuntu 18.04',     1,  'Ubuntu v18.04'),
('SecurityStatus',      1,      'In Progress',      0,  'In-Progress status for security check'),
('SecurityStatus',      2,      'Completed',        1,  'Completed status for security check'),
('SecurityStatus',      3,      'Failed',           2,  'Failed status for security check'),
('SecurityItemStatus',  1,      'Pass',             0,  'Passed status for security check item'),
('SecurityItemStatus',  2,      'Failed',           1,  'Failed status for security check item'),
('UserStatus',          1,      'Activated',        0,  'Activated status of user'),
('UserStatus',          2,      'Deactivated',    	1,  'Deactivated status of user'),
('ClusterStatus',       1,      'Activated',        0,  'Activated status of cluster'),
('NodeStatus',          1,      'Activated',        0,  'Activated status of node')



-- -- Cloud

-- INSERT INTO "edgecraft"."tbl_cloud"
-- (
--     cloud_uid, 
--     name, 
--     type,
--     description,
--     state,
--     creator
-- )
-- VALUES
-- ('0976a14f-88b7-4c0d-81ee-e856d217ede9',	'Flexible, Restful access...',  '2',	'test', 	'3',	'system'),
-- ('60ce55ff-f4f9-492d-91be-3a88388cb520',	'Collect configure and ...',    '2',	'test',	    '1',    'system'),
-- ('89c83f49-6552-4900-83fc-6a0320775e0d',	'Non-profit charity...',    	'1',	'test',	    '3',    'system'),
-- ('ce2a80b1-c5d3-4db3-8567-f36cdfb3de6f',	'Search for logos and ...',     '1',	'test',	    '1',	'system'),
-- ('e4f87fcc-1c67-423d-8c98-5a4f6d27ed85',	'Registered Domain Names...',   '2',	'test', 	'2',	'system')

-- -- Cloud - Cluster

-- INSERT INTO "edgecraft"."tbl_cloud_cluster"
-- (
--     "cloud_uid"
-- 	"k8s_version"						VARCHAR(30) 	DEFAULT '2'								NOT NULL,	-- 클라우드쿠버네티스버전
-- 	"bmc_credential_secret"				VARCHAR(50)  											NULL,     	-- 클라우드클러스터BMC자격증명SECRET
-- 	"bmc_credential_user"     			VARCHAR(50)  											NULL,     	-- 클라우드클러스터BMC자격증명사용자
-- 	"bmc_credential_password" 			VARCHAR(100) 											NULL,     	-- 클라우드클러스터BMC자격증명비밀번호
-- 	"pod_cidr"                			VARCHAR(100) 											NULL,     	-- 클라우드클러스터포드CIDR
-- 	"service_cidr"            			VARCHAR(100) 											NULL,     	-- 클라우드클러스터서비스CIDR
-- 	"loadbalancer_use_yn"     			BOOLEAN      	DEFAULT FALSE							NULL,     	-- 클라우드클러스터로드밸런서사용여부
-- 	"loadbalancer_address"    			VARCHAR(50)  											NULL,     	-- 클라우드클러스터로드밸런서주소
-- 	"loadbalancer_port"       			VARCHAR(6)   											NULL,     	-- 클라우드클러스터로드밸런서포트
-- 	"image_url"               			VARCHAR(255) 											NULL,     	-- 클라우드클러스터이미지URL
-- 	"image_checksum"          			VARCHAR(255) 											NULL,     	-- 클라우드클러스터이미지CHECKSUM
-- 	"image_checksum_type"     			VARCHAR(30)  											NULL,     	-- 클라우드클러스터이미지CHECKSUM종류
-- 	"image_format"            			VARCHAR(30)  											NULL,     	-- 클라우드클러스터이미지포맷
-- 	"master_extra_config"     			JSON         											NULL,     	-- 클라우드클러스터마스터엑스트라구성
-- 	"worker_extra_config"     			JSON         											NULL,     	-- 클라우드클러스터워커엑스트라구성
-- 	"state"                   			VARCHAR(30)  											NULL,     	-- 클라우드클러스터상태
-- 	"external_etcd_use"       			CHAR(1)      	DEFAULT 'N'								NULL, 		-- 클라우드클러스터외부ETCD사용
-- 	"external_etcd_endpoints"       	JSON         											NULL,     	-- 외부ETCD엔드포인츠
-- 	"external_etcd_certificate_ca"  	TEXT         											NULL,     	-- 외부ETCD인증서CA
-- 	"external_etcd_certificate_cert"	TEXT         											NULL,     	-- 외부ETCD인증서CERT
-- 	"external_etcd_certificate_key" 	TEXT         											NULL,     	-- 외부ETCD인증서KEY
-- 	"cloud_cluster_storage_class"   	JSON         											NULL,     	-- 클라우드클러스터스토리지클래스
-- 	"creator"                       	VARCHAR(50)  											NOT NULL,   -- 생성자
-- 	"created_at"                    	TIMESTAMP    	DEFAULT NOW()							NOT NULL,   -- 생성일시
-- 	"updater"                       	VARCHAR(50)  											NULL,     	-- 수정자
-- 	"updated_at"                    	TIMESTAMP    											NULL      	-- 수정일시
-- )
-- VALUES
-- ('2', )
-- -- Cloud - Cluster - Node

