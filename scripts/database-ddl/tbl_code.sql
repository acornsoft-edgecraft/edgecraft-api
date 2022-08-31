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
	"code_uid"       SERIAL      NOT NULL, -- 코드식별자
	"code_name"      VARCHAR(50) NULL,     -- 코드이름
	"display_name"   VARCHAR(50) NULL,     -- 노출이름
	"display_order"  INTEGER     NULL,     -- 표시순서
	"use_yn"         CHAR(1)     NULL     DEFAULT 'N', -- 사용여부
	"group_code_uid" INT4        NULL,     -- 그룹코드식별자
	"creator"        INTEGER     NOT NULL, -- 생성자
	"created"        TIMESTAMP   NULL,     -- 생성일시
	"updater"        INTEGER     NULL,     -- 수정자
	"updated"        TIMESTAMP   NULL      -- 수정일시
)
WITH (
OIDS=false
);

-- 공통코드 기본키
CREATE UNIQUE INDEX "PK_tbl_code"
	ON "edgecraft"."tbl_code"
	( -- 공통코드
		"code_uid" ASC -- 코드식별자
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
			"group_code_uid" -- 그룹코드식별자
		)
		REFERENCES "edgecraft"."tbl_code_group" ( -- 코드그룹
			"group_code_uid" -- 그룹코드식별자
		)
		ON UPDATE NO ACTION ON DELETE NO ACTION
		NOT VALID;