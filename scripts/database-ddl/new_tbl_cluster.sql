-- 오픈스택 클러스터
ALTER TABLE "edgecraft"."tbl_cluster"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cluster"; -- 오픈스택 클러스터 기본키

-- 오픈스택 클러스터 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cluster";

-- 오픈스택 클러스터
DROP TABLE IF EXISTS "edgecraft"."tbl_cluster";

-- 오픈스택 클러스터
CREATE TABLE "edgecraft"."tbl_cluster"
(
    "cloud_uid"             			CHAR(36)                            NOT NULL,   -- 클라우드식별자
	"cluster_uid"           			CHAR(36)                            NOT NULL,   -- 클러스터식별자
	"namespace"							VARCHAR(100)	DEFAULT 'default'	NULL,		-- CR정보 생성을 위한 Namespace
	"name"                  			VARCHAR(50)                         NULL,       -- 클러스터이름
	"description"           			VARCHAR(100)                        NULL,       -- 클러스터설명
	"credential"            			TEXT                                NULL,       -- 클러스터자격증명


	-- K8s 정보
	"version"               			INTEGER         DEFAULT 1           NOT NULL,   -- 클러스터버전 (Kubernetes)
	"pod_cidr"              			VARCHAR(30)                         NULL,       -- 클러스터포드CIDR
	"service_cidr"          			VARCHAR(30)                         NULL,       -- 클러스터서비스CIDR
	"service_domain"          			VARCHAR(30)                         NULL,       -- 클러스터서비스도메인

	-- Openstack 정보
	"openstack_info"        			JSON                                NULL,       -- 클러스터오픈스택정보

	-- nodeset 정보
	"loadbalancer_use_yn"     			BOOLEAN      	DEFAULT FALSE		NULL,     	-- 클러스터로드밸런서사용여부

	-- ETCD 정보	
	"external_etcd_use"       			BOOLEAN      	DEFAULT FALSE		NULL, 		-- 클라우드클러스터외부ETCD사용
	"external_etcd_endpoints"       	JSON         						NULL,     	-- 외부ETCD엔드포인츠
	"external_etcd_certificate_ca"  	TEXT         						NULL,     	-- 외부ETCD인증서CA
	"external_etcd_certificate_cert"	TEXT         						NULL,     	-- 외부ETCD인증서CERT
	"external_etcd_certificate_key" 	TEXT         						NULL,     	-- 외부ETCD인증서KEY

	-- Storage 정보	
	"storage_class"   					JSON         						NULL,     	-- 클라우드클러스터스토리지클래스

	-- "storage_user_id"       			VARCHAR(50)                         NULL,       -- 스토리지사용자ID
	-- "storage_user_secret"   			VARCHAR(100)                        NULL,       -- 스토리지사용자비밀
	"state"                 			INTEGER         DEFAULT 1           NOT NULL,   -- 클러스터상태 (Code - ClusterStatus)
	"creator"               			VARCHAR(30)     DEFAULT 'system'    NOT NULL,   -- 생성자
	"created_at"            			TIMESTAMP       DEFAULT NOW()       NOT NULL,   -- 생성일시
	"updater"               			VARCHAR(30)                         NULL,       -- 수정자
	"updated_at"            			TIMESTAMP                           NULL        -- 수정일시
)
WITH (
OIDS=false
);

-- 오픈스택 클러스터 기본키
CREATE UNIQUE INDEX "PK_tbl_cluster"
	ON "edgecraft"."tbl_cluster"
	( -- 오픈스택 클러스터
		"cluster_uid" ASC -- 클러스터식별자
	);

-- 오픈스택 클러스터
ALTER TABLE "edgecraft"."tbl_cluster"
	ADD CONSTRAINT "PK_tbl_cluster"
		-- 오픈스택 클러스터 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cluster"
	NOT DEFERRABLE;