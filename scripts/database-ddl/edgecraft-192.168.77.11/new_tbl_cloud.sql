-- 클라우드
ALTER TABLE "edgecraft"."tbl_cloud"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cloud"; -- 클라우드 기본키

-- 클라우드 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cloud";

-- 클라우드
DROP TABLE IF EXISTS "edgecraft"."tbl_cloud";

-- 클라우드
CREATE TABLE "edgecraft"."tbl_cloud"
(
	"cloud_uid"         CHAR(36)							NOT NULL, 	-- 클라우드식별자
	"name"        		VARCHAR(30)  						NOT NULL,   -- 클라우드이름
	"type"        		INTEGER			DEFAULT 1  			NOT NULL,   -- 클라우드유형 (Code - CloudTypes)
	"description" 		VARCHAR(300) 						NULL,     	-- 클라우드설명
	"state"       		INTEGER			DEFAULT 1	  		NOT NULL,   -- 클라우드상태 (Code - CloudStatus)
	"creator"           VARCHAR(30)  	DEFAULT 'system'	NOT NULL,   -- 생성자
	"created_at"        TIMESTAMP   	DEFAULT NOW()		NOT NULL,   -- 생성일시
	"updater"           VARCHAR(30)  						NULL,     	-- 수정자
	"updated_at"        TIMESTAMP    						NULL      	-- 수정일시
)
WITH (
OIDS=false
);

-- 클라우드 기본키
CREATE UNIQUE INDEX "PK_tbl_cloud"
	ON "edgecraft"."tbl_cloud"
	( -- 클라우드
		"cloud_uid" ASC -- 클라우드식별자
	);

-- 클라우드
ALTER TABLE "edgecraft"."tbl_cloud"
	ADD CONSTRAINT "PK_tbl_cloud"
		-- 클라우드 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cloud"
	NOT DEFERRABLE;