-- 클라우드 노드
ALTER TABLE "edgecraft"."tbl_cloud_node"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cloud_node"; -- 클라우드 노드 기본키

-- 클라우드 노드 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cloud_node";

-- 클라우드 노드
DROP TABLE IF EXISTS "edgecraft"."tbl_cloud_node";

-- 클라우드 노드
CREATE TABLE "edgecraft"."tbl_cloud_node"
(
	"cloud_uid"                 CHAR(36)         					NOT NULL, -- 클라우드식별자
	"cluster_uid"               CHAR(36)         					NOT NULL, -- 클라우드클러스터식별자
	"node_uid"                  CHAR(36)		        			NOT NULL, -- 클라우드노드식별자
			
	-- Baremetal Host 정보		
	"host_name"            		VARCHAR(30)  						NULL,     -- 클라우드노드호스트명
	"bmc_address"          		VARCHAR(255) 						NULL,     -- 클라우드노드BMC주소
	"mac_address"          		VARCHAR(255) 						NULL,     -- 클라우드노드MAC주소
	"boot_mode"            		INTEGER  							NULL,     -- 클라우드노드부트모드 (Code - BootModes)
	"online_power"         		BOOLEAN      						NULL,     -- 클라우드노드온라인파워
	"external_provisioning"		BOOLEAN      						NULL,     -- 클라우드노드외부프로비저닝
			
	-- Node 정보		
	"name"            			VARCHAR(30)  						NULL,     -- 클라우드노드이름
	"ipaddress"            		VARCHAR(30)  						NULL,     -- 클라우드노드IP
	"labels"               		JSON         						NULL,     -- 클라우드노드라벨
			
	-- Openstack Ceph Path (? - 화면에 없음, 향후 조정 		
	"osd_path"                  VARCHAR(300) 						NULL,     -- OSD경로
		
	-- 기본 정보		
	"type"                 		INTEGER  		DEFAULT 1			NOT NULL, -- 클라우드노드종류 (Code - NodeTypes)
	"state"                		INTEGER  		DEFAULT 1			NOT NULL, -- 클라우드노드상태 (Code - NodeStatus)

	"creator"                   VARCHAR(30)  	DEFAULT 'system'	NOT NULL,     -- 생성자
	"created_at"                TIMESTAMP    	DEFAULT NOW()		NOT NULL,     -- 생성일시
	"updater"                   VARCHAR(30)  						NULL,     -- 수정자
	"updated_at"                TIMESTAMP    						NULL      -- 수정일시
)
WITH (
OIDS=false
);

-- 클라우드 노드 기본키
CREATE UNIQUE INDEX "PK_tbl_cloud_node"
	ON "edgecraft"."tbl_cloud_node"
	( -- 클라우드 노드
		"node_uid" ASC, -- 클라우드노드식별자
		"cloud_uid" ASC, -- 클라우드식별자
		"cluster_uid" ASC -- 클라우드클러스터식별자
	);

-- 클라우드 노드
ALTER TABLE "edgecraft"."tbl_cloud_node"
	ADD CONSTRAINT "PK_tbl_cloud_node"
		-- 클라우드 노드 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cloud_node"
	NOT DEFERRABLE;