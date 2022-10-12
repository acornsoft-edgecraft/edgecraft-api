-- 오픈스택클러스터노드셋
ALTER TABLE "edgecraft"."tbl_nodeset"
	DROP CONSTRAINT IF EXISTS "PK_tbl_nodeset"; -- 오픈스택클러스터노드셋 기본키

-- 오픈스택클러스터노드셋 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_nodeset";

-- 오픈스택클러스터노드셋
DROP TABLE IF EXISTS "edgecraft"."tbl_nodeset";

-- 오픈스택클러스터노드셋
CREATE TABLE "edgecraft"."tbl_nodeset"
(
  	"cluster_uid"   CHAR(36)                            NOT NULL, -- 클러스터식별자
	"nodeset_uid"   CHAR(36)                            NOT NULL, -- 노드셋식별자
	"type"          INTEGER         DEFAULT 1           NOT NULL,     -- 노드셋유형
	"namespace"     VARCHAR(50)                         NULL,     -- 노드셋네임스페이스
	"name"          VARCHAR(50)                         NULL,     -- 노드셋이름
	"node_count"    INTEGER         DEFAULT 0           NULL,     -- 노드셋노드개수
	"flavor"        VARCHAR(50)                         NULL,     -- 노드셋FLAVOR
	"label"         TEXT                                NULL,     -- 노드셋라벨
	"creator"       VARCHAR(30)     DEFAULT 'system'    NOT NULL,     -- 생성자
	"created_at"    TIMESTAMP       DEFAULT NOW()       NOT NULL,     -- 생성일시
	"updater"       VARCHAR(30)                         NULL,     -- 수정자
	"updated_at"    TIMESTAMP                           NULL      -- 수정일시
)
WITH (
OIDS=false
);

-- 오픈스택클러스터노드셋 기본키
CREATE UNIQUE INDEX "PK_tbl_nodeset"
	ON "edgecraft"."tbl_nodeset"
	( -- 오픈스택클러스터노드셋
  		"cluster_uid" ASC, -- 클러스터식별자
		"nodeset_uid" ASC  -- 노드셋식별자

	);

-- 오픈스택클러스터노드셋
ALTER TABLE "edgecraft"."tbl_nodeset"
	ADD CONSTRAINT "PK_tbl_nodeset"
		-- 오픈스택클러스터노드셋 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_nodeset"
	NOT DEFERRABLE;