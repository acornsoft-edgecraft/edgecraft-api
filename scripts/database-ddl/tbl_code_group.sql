-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	DROP CONSTRAINT IF EXISTS "PK_tbl_code_group"; -- 코드그룹 기본키

-- 코드그룹
DROP TABLE IF EXISTS "edgecraft"."tbl_code_group";

-- 코드그룹
CREATE TABLE "edgecraft"."tbl_code_group"
(
	"group_code_uid"  SERIAL        NOT NULL, -- 그룹코드식별자
	"group_code_name" VARCHAR(50)   NULL,     -- 그룹코드이름
	"use_yn"          CHAR(1)       NULL     DEFAULT 'N', -- 사용여부
	"description"     VARCHAR(3000) NULL      -- 설명
)
WITH (
OIDS=false
);

-- 코드그룹 기본키
CREATE UNIQUE INDEX "PK_tbl_code_group"
	ON "edgecraft"."tbl_code_group"
	( -- 코드그룹
		"group_code_uid" ASC -- 그룹코드식별자
	)
;
-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	ADD CONSTRAINT "PK_tbl_code_group"
		 -- 코드그룹 기본키
	PRIMARY KEY 
	USING INDEX "PK_tbl_code_group"
	NOT DEFERRABLE;