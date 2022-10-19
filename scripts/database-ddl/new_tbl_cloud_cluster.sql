-- 클라우드 클러스터
ALTER TABLE "edgecraft"."tbl_cloud_cluster"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cloud_cluster"; -- 클라우드 클러스터 기본키

-- 클라우드 클러스터 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cloud_cluster";

-- 클라우드 클러스터
DROP TABLE IF EXISTS "edgecraft"."tbl_cloud_cluster";

-- 클라우드 클러스터
CREATE TABLE "edgecraft"."tbl_cloud_cluster"
(
	"cloud_uid"                       	CHAR(36)         					NOT NULL, 	-- 클라우드식별자
	"cluster_uid"					  	CHAR(36)     						NOT NULL, 	-- 클라우드클러스터식별자
		
	-- K8S 정보	
	"k8s_version"						INTEGER 	DEFAULT 1				NOT NULL,	-- 클라우드쿠버네티스버전 (Code - K8sVersions)
	"pod_cidr"                			VARCHAR(30) 						NULL,     	-- 클라우드클러스터포드CIDR
	"service_cidr"            			VARCHAR(30) 						NULL,     	-- 클라우드클러스터서비스CIDR
	"service_domain"           			VARCHAR(30) 						NULL,     	-- 클라우드클러스터서비스도메인
		
	-- Baremetal 정보
	"namespace"							VARCHAR(100)	DEFAULT 'default'	NULL,		-- CR정보 생성을 위한 Namespace
	"bmc_credential_secret"				VARCHAR(50)  						NULL,     	-- 클라우드클러스터BMC자격증명SECRET
	"bmc_credential_user"     			VARCHAR(50)  						NULL,     	-- 클라우드클러스터BMC자격증명사용자
	"bmc_credential_password" 			VARCHAR(100) 						NULL,     	-- 클라우드클러스터BMC자격증명비밀번호
	"image_url"               			VARCHAR(255) 						NULL,     	-- 클라우드클러스터이미지URL
	"image_checksum"          			VARCHAR(255) 						NULL,     	-- 클라우드클러스터이미지CHECKSUM
	"image_checksum_type"     			INTEGER  							NULL,     	-- 클라우드클러스터이미지CHECKSUM종류 (Code - ImageChecksumTypes)
	"image_format"            			INTEGER  							NULL,     	-- 클라우드클러스터이미지포맷 (Code - ImageFormats)
	"master_extra_config"     			JSON         						NULL,     	-- 클라우드클러스터마스터엑스트라구성
	"worker_extra_config"     			JSON         						NULL,     	-- 클라우드클러스터워커엑스트라구성
		
	-- nodes 정보	
	"loadbalancer_use_yn"     			BOOLEAN      	DEFAULT FALSE		NULL,     	-- 클라우드클러스터로드밸런서사용여부
	"loadbalancer_address"    			VARCHAR(30)  						NULL,     	-- 클라우드클러스터로드밸런서주소
	"loadbalancer_port"       			VARCHAR(6)   						NULL,     	-- 클라우드클러스터로드밸런서포트
		
	-- ETCD 정보	
	"external_etcd_use"       			BOOLEAN      	DEFAULT FALSE		NULL, 		-- 클라우드클러스터외부ETCD사용
	"external_etcd_endpoints"       	JSON         						NULL,     	-- 외부ETCD엔드포인츠
	"external_etcd_certificate_ca"  	TEXT         						NULL,     	-- 외부ETCD인증서CA
	"external_etcd_certificate_cert"	TEXT         						NULL,     	-- 외부ETCD인증서CERT
	"external_etcd_certificate_key" 	TEXT         						NULL,     	-- 외부ETCD인증서KEY
		
	-- Storage 정보	
	"storage_class"   					JSON         						NULL,     	-- 클라우드클러스터스토리지클래스
		
	-- 기본 정보	
	"state"                   			INTEGER	  		DEFAULT 1			NOT NULL,  	-- 클라우드클러스터상태 (Code - CloudStatus)
	
	"creator"                       	VARCHAR(30)  	DEFAULT 'system'	NOT NULL,   -- 생성자
	"created_at"                    	TIMESTAMP    	DEFAULT NOW()		NOT NULL,   -- 생성일시
	"updater"                       	VARCHAR(30)  						NULL,     	-- 수정자
	"updated_at"                    	TIMESTAMP    						NULL      	-- 수정일시
)
WITH (
OIDS=false
);

-- 클라우드 클러스터 기본키
CREATE UNIQUE INDEX "PK_tbl_cloud_cluster"
	ON "edgecraft"."tbl_cloud_cluster"
	( -- 클라우드 클러스터
		"cluster_uid" 	ASC, 	-- 클라우드클러스터식별자
		"cloud_uid" 	ASC 	-- 클라우드식별자
	);

-- 클라우드 클러스터
ALTER TABLE "edgecraft"."tbl_cloud_cluster"
	ADD CONSTRAINT "PK_tbl_cloud_cluster"
		-- 클라우드 클러스터 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cloud_cluster"
	NOT DEFERRABLE;