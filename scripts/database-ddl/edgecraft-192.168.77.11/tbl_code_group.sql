-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	DROP CONSTRAINT IF EXISTS "PK_tbl_code_group"; -- 코드그룹 기본키

-- 코드그룹
DROP TABLE IF EXISTS "edgecraft"."tbl_code_group";

-- 코드그룹
CREATE TABLE "edgecraft"."tbl_code_group"
(
	"code_group_uid"         uuid          NOT NULL, -- 코드그룹식별자
	"code_group_name"        VARCHAR(50)   NULL,     -- 코드그룹이름
	"code_group_description" VARCHAR(3000) NULL,     -- 코드그룹설명
	"use_yn"                 BOOLEAN       NULL,     -- 사용여부
	"creator"                VARCHAR(50)   NOT NULL, -- 생성자
	"created_at"             TIMESTAMP     NULL,     -- 생성일시
	"updater"                VARCHAR(50)   NULL,     -- 수정자
	"updated_at"             TIMESTAMP     NULL      -- 수정일시
)
WITH (
OIDS=false
);

-- 코드그룹 기본키
CREATE UNIQUE INDEX "PK_tbl_code_group"
	ON "edgecraft"."tbl_code_group"
	( -- 코드그룹
		"code_group_uid" ASC -- 코드그룹식별자
	)
;
-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	ADD CONSTRAINT "PK_tbl_code_group"
		 -- 코드그룹 기본키
	PRIMARY KEY 
	USING INDEX "PK_tbl_code_group"
	NOT DEFERRABLE;