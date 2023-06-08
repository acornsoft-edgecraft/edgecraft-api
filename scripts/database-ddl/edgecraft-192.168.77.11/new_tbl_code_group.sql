-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	DROP CONSTRAINT IF EXISTS "PK_tbl_code_group"; -- 코드그룹 기본키

-- 코드그룹
DROP TABLE IF EXISTS "edgecraft"."tbl_code_group";

-- 코드그룹
CREATE TABLE "edgecraft"."tbl_code_group"
(
	"group_id"		VARCHAR(30)							NOT NULL,	-- 코드그룹식별자
	"description" 	VARCHAR(3000) 						NULL,     	-- 코드그룹설명
	"use_yn"        BOOLEAN 		DEFAULT TRUE		NULL,     	-- 사용여부
	"creator"       VARCHAR(30)   	DEFAULT 'system'	NOT NULL, 	-- 생성자
	"created_at"    TIMESTAMP     	DEFAULT NOW()		NULL,     	-- 생성일시
	"updater"       VARCHAR(30)   						NULL,     	-- 수정자
	"updated_at"    TIMESTAMP     						NULL      	-- 수정일시
)
WITH (
OIDS=false
);

-- 코드그룹 기본키
CREATE UNIQUE INDEX "PK_tbl_code_group"
	ON "edgecraft"."tbl_code_group"
	( -- 코드그룹
		"group_id" ASC -- 코드그룹식별자
	)
;
-- 코드그룹
ALTER TABLE "edgecraft"."tbl_code_group"
	ADD CONSTRAINT "PK_tbl_code_group"
		 -- 코드그룹 기본키
	PRIMARY KEY 
	USING INDEX "PK_tbl_code_group"
	NOT DEFERRABLE;