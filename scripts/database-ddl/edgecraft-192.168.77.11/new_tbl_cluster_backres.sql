-- 오픈스택 클러스터 백업/복원
ALTER TABLE "edgecraft"."tbl_cluster_backres"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cluster_backres"; -- 오픈스택 클러스터 백업/복원 기본키

-- 오픈스택 클러스터 백업/복원 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cluster_backres";

-- 오픈스택 클러스터 백업/복원
DROP TABLE IF EXISTS "edgecraft"."tbl_cluster_backres";

-- 오픈스택 클러스터 백업/복원
CREATE TABLE "edgecraft"."tbl_cluster_backres"
(
    "cloud_uid"			CHAR(36)                            NOT NULL,   -- 클라우드식별자
	"cluster_uid"       CHAR(36)                            NOT NULL,   -- 클러스터식별자
	"backres_uid"       CHAR(36)                            NOT NULL,   -- 클러스터백업/복원식별자
	"name"				VARCHAR(200)						NOT NULL,	-- 백업/복원 이름
	"type"				CHAR(1)			DEFAULT 'B'			NOT NULL,	-- Backup ('B') / Restore ('R') 구분자
	"status"			CHAR(1)			DEFAULT 'R'			NOT NULL,	-- Running ('R') / Completed ('C') / Failed ('F') / PartiallyFailed ('P') 상태 구분자 
	"reason"			VARCHAR(200)						NULL,		-- 실패 이유
	"backup_name"		VARCHAR(200)						NULL,		-- 복원 시 사용할 백업 명
	"creator"           VARCHAR(30)     DEFAULT 'system'    NOT NULL,   -- 생성자
	"created_at"        TIMESTAMP       DEFAULT NOW()       NOT NULL,   -- 생성일시
	"updater"           VARCHAR(30)                         NULL,       -- 수정자
	"updated_at"        TIMESTAMP                           NULL        -- 수정일시
)
WITH (
OIDS=false
);

-- 오픈스택 클러스터 백업/복원 기본키
CREATE UNIQUE INDEX "PK_tbl_cluster_backres"
	ON "edgecraft"."tbl_cluster_backres"
	( -- 오픈스택 클러스터 백업/복원
		"backres_uid" ASC -- 클러스터백업/복원식별자
	);

-- 오픈스택 클러스터 백업/복원
ALTER TABLE "edgecraft"."tbl_cluster_backres"
	ADD CONSTRAINT "PK_tbl_cluster_backres"
		-- 오픈스택 클러스터 백업/복원 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cluster_backres"
	NOT DEFERRABLE;