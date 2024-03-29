-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	DROP CONSTRAINT IF EXISTS "PK_tbl_user"; -- 사용자 기본키2

-- 사용자
DROP TABLE IF EXISTS "edgecraft"."tbl_user";

-- 사용자
CREATE TABLE "edgecraft"."tbl_user"
(
	"user_uid"                       UUID         NOT NULL, -- 사용자식별자
	"user_role"                      SMALLINT     NULL,     -- 사용자권한
	"user_name"                      VARCHAR(50)  NULL,     -- 사용자이름
	"user_id"                        VARCHAR(50)  NOT NULL, -- 사용자ID
	"password"                       VARCHAR(100) NULL,     -- 비밀번호
	"email"                          VARCHAR(40)  NULL,     -- 이메일
	"last_login"                     TIMESTAMP    NULL,     -- 마지막로그인
	"password_expiration_begin_time" TIMESTAMP    NOT NULL, -- 비밀번호만료시작시간
	"password_expiration_end_time"   TIMESTAMP    NOT NULL, -- 비밀번호만료끝시간
	"reset_password_yn"              CHAR(1)      NULL     DEFAULT 'N', -- 초기화비밀번호여부
	"active_datetime"                TIMESTAMP    NULL,     -- 활성화일시
	"inactive_yn"                    CHAR(1)      NULL     DEFAULT 'N', -- 비활성여부
	"user_state"                     CHAR(1)      NULL,     -- 사용자상태
	"creator"                        VARCHAR(50)  NOT NULL, -- 생성자
	"created_at"                     TIMESTAMP    NULL,     -- 생성일시
	"updater"                        VARCHAR(50)  NULL,     -- 수정자
	"updated_at"                     TIMESTAMP    NULL      -- 수정일시
)
WITH (
OIDS=false
);

-- 사용자 기본키2
CREATE UNIQUE INDEX "PK_tbl_user"
	ON "edgecraft"."tbl_user"
	( -- 사용자
		"user_uid" ASC -- 사용자식별자
	)
;
-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	ADD CONSTRAINT "PK_tbl_user"
		 -- 사용자 기본키2
	PRIMARY KEY 
	USING INDEX "PK_tbl_user"
	NOT DEFERRABLE;