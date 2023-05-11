-- 오픈스택 클러스터 벤치마크
ALTER TABLE "edgecraft"."tbl_cluster_benchmarks"
	DROP CONSTRAINT IF EXISTS "PK_tbl_cluster_benchmarks"; -- 오픈스택 클러스터 벤치마크 기본키

-- 오픈스택 클러스터 벤치마크 기본키
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_cluster_benchmarks";

-- 오픈스택 클러스터 벤치마크
DROP TABLE IF EXISTS "edgecraft"."tbl_cluster_benchmarks";

-- 오픈스택 클러스터 벤치마크
CREATE TABLE "edgecraft"."tbl_cluster_benchmarks"
(
    "cloud_uid"             			CHAR(36)                            NOT NULL,   -- 클라우드식별자
	"cluster_uid"           			CHAR(36)                            NOT NULL,   -- 클러스터식별자
	"benchmarks_uid"           			CHAR(36)                            NOT NULL,   -- 클러스터벤치마크식별자
	"cis_version"						VARCHAR(30)							NULL,		-- cis benchmarks version
	"detected_version"					VARCHAR(30)							NULL,		-- kubernetes version
	"results"     						JSON         						NULL,     	-- 실행결과
	"totals"     						JSON         						NULL,     	-- 실행결과
	"success_yn"						BOOLEAN      	DEFAULT FALSE		NULL, 		-- 성공 여부
	"reason"							VARCHAR(200)						NULL,		-- 실패 이유
	"creator"               			VARCHAR(30)     DEFAULT 'system'    NOT NULL,   -- 생성자
	"created_at"            			TIMESTAMP       DEFAULT NOW()       NOT NULL,   -- 생성일시
	"updater"               			VARCHAR(30)                         NULL,       -- 수정자
	"updated_at"            			TIMESTAMP                           NULL        -- 수정일시
)
WITH (
OIDS=false
);

-- 오픈스택 클러스터 벤치마크 기본키
CREATE UNIQUE INDEX "PK_tbl_cluster_benchmarks"
	ON "edgecraft"."tbl_cluster_benchmarks"
	( -- 오픈스택 클러스터 벤치마크
		"benchmarks_uid" ASC -- 클러스터벤치마크식별자
	);

-- 오픈스택 클러스터 벤치마크
ALTER TABLE "edgecraft"."tbl_cluster_benchmarks"
	ADD CONSTRAINT "PK_tbl_cluster_benchmarks"
		-- 오픈스택 클러스터 벤치마크 기본키
	PRIMARY KEY
	USING INDEX "PK_tbl_cluster_benchmarks"
	NOT DEFERRABLE;