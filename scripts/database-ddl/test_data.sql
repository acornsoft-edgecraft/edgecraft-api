-- User
INSERT INTO "edgecraft"."tbl_user"
(
    user_role,
    user_name,
    user_id,
    password,
    email,
    
)
VALUES
( 1, 'morris', 'ccambo', '662edc7bc97f9425be9b3abd61f9eea63fbff1cf144a9a572bab168145fc8838ac2f9f92e284f22a974e3820f947a909190d686e4fe9476316ca7d36fd9965d7', 'ccambo@acornsoft.io' )
-- pwd >> 1234abcd@Acorn

-- Code Group
INSERT INTO "edgecraft"."tbl_code_group"
(
	"group_id"
)
VALUES 
('CloudTypes'),
('CloudStatus'),
('K8sVersions'),
('ImageChecksumTypes'),
('ImageFormats'),
('BootMode'),
('NodeTypes'),
('UserRoles')

-- Code
INSERT INTO "edgecraft"."tbl_code"
(
	"group_id",
	"code",
	"name",
	"display_order"
)
VALUES
('CloudTypes', 			'1', 		'Baremetal', 		0),
('CloudTypes', 			'2', 		'Openstack', 		1),
('CloudStatus', 		'1', 		'Saved', 			0),
('CloudStatus', 		'2', 		'Provisioning',		1),
('CloudStatus', 		'3', 		'Provisioned', 		2),
('CloudStatus', 		'4', 		'Failed',			3),
('CloudStatus', 		'5', 		'Deleting', 		4),
('CloudStatus', 		'6', 		'Deleted', 			5),
('K8sVersions', 		'1', 		'1.22.0',			0),
('K8sVersions', 		'2', 		'1.22.3',			0),
('ImageChecksumTypes',	'1',		'md5',				0),
('ImageChecksumTypes',	'2',		'sha256',			1),
('ImageChecksumTypes',	'3',		'sha512',			2),
('ImageFormats',		'1',		'raw',				0),
('ImageFormats',		'2',		'qcow2',			1),
('ImageFormats',		'3',		'vdi',				2),
('ImageFormats',		'4',		'vmdk',				3),
('ImageFormats',		'5',		'live-iso',			4),
('BootMode',			'1',		'UEFI',				0),
('BootMode',			'2',		'legacy',			1),
('BootMode',			'3',		'UEFISecureBoot',	2),
('NodeTypes',			'1',		'MASTER',			0),
('NodeTypes',			'2',		'WORKER',			1),
('UserRoles',			'1',		'Admin',			0),
('UserRoles',			'2',		'Manager',			0),
('UserRoles',			'3',		'User',				0)

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

