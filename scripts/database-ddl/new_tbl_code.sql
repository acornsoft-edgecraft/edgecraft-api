-- 공통코드
ALTER TABLE "edgecraft"."tbl_code"
	DROP CONSTRAINT IF EXISTS "FK_tbl_code_group_TO_tbl_code"; -- 코드그룹 -> 공통코드

-- 공통코드
ALTER TABLE "edgecraft"."tbl_code"
	DROP CONSTRAINT IF EXISTS "PK_tbl_code"; -- 공통코드 기본키

-- 공통코드
DROP TABLE IF EXISTS "edgecraft"."tbl_code";

-- 공통코드
CREATE TABLE "edgecraft"."tbl_code"
(
	"group_id"			VARCHAR(30)							NOT NULL,	-- 코드그룹식별자
	"code"				VARCHAR(30)							NOT NULL, 	-- 코드식별자
	"name"          	VARCHAR(50)   						NOT NULL,   -- 코드이름
	"display_order" 	INTEGER       						NULL,     	-- 코드표시순서
	"description"   	VARCHAR(3000) 						NULL,     	-- 코드설명
	"use_yn"            BOOLEAN       	DEFAULT 'Y'			NOT NULL,   -- 사용여부
	"creator"           VARCHAR(50)   	DEFAULT 'system'	NOT NULL, 	-- 생성자
	"created_at"        TIMESTAMP     	DEFAULT NOW()		NOT NULL,   -- 생성일시
	"updater"           VARCHAR(50)   						NULL,     	-- 수정자
	"updated_at"        TIMESTAMP     						NULL      	-- 수정일시
)
WITH (
OIDS=false
);

-- 공통코드 기본키
CREATE UNIQUE INDEX "PK_tbl_code"
	ON "edgecraft"."tbl_code"
	(
		"group_id" ASC,	-- 코드그룹식별자
		"code" ASC 		-- 코드식별자
	)
;
-- 공통코드
ALTER TABLE "edgecraft"."tbl_code"
	ADD CONSTRAINT "PK_tbl_code"
		 -- 공통코드 기본키
	PRIMARY KEY 
	USING INDEX "PK_tbl_code"
	NOT DEFERRABLE;

-- 공통코드
ALTER TABLE "edgecraft"."tbl_code"
	ADD CONSTRAINT "FK_tbl_code_group_TO_tbl_code"
	 -- 코드그룹 -> 공통코드
		FOREIGN KEY (
			"group_id" -- 코드그룹식별자
		)
		REFERENCES "edgecraft"."tbl_code_group" (
			"group_id" -- 코드그룹식별자
		)
		ON UPDATE NO ACTION ON DELETE NO ACTION
		NOT VALID;