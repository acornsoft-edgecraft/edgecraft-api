-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	DROP CONSTRAINT IF EXISTS "UK_tbl_user"; -- 사용자 유니크 제약

-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	DROP CONSTRAINT IF EXISTS "PK_tbl_user"; -- 사용자 기본키2

-- 사용자 유니크 인덱스
DROP INDEX IF EXISTS "edgecraft"."UK_tbl_user";

-- 사용자 기본키2
DROP INDEX IF EXISTS "edgecraft"."PK_tbl_user";

-- 사용자
DROP TABLE IF EXISTS "edgecraft"."tbl_user";

-- 사용자
CREATE TABLE "edgecraft"."tbl_user"
(
	"user_uid"                       	CHAR(36)							NOT NULL, 	-- 사용자식별자
	"role"                      		INT	     							NOT NULL, 	-- 사용자권한 (Code - UserRoles)
	"name"                      		VARCHAR(50)  						NULL,     	-- 사용자이름
	"id"                        		VARCHAR(30)  						NOT NULL, 	-- 사용자ID
	"password"                       	VARCHAR(128) 						NOT NULL, 	-- 비밀번호
	"email"                          	VARCHAR(50)  						NOT NULL, 	-- 이메일
	"last_login"                     	TIMESTAMP    						NULL,     	-- 마지막로그인
	"password_expiration_begin_time" 	TIMESTAMP    						NOT NULL, 	-- 비밀번호만료시작시간
	"password_expiration_end_time"   	TIMESTAMP    						NOT NULL, 	-- 비밀번호만료끝시간
	"reset_password_yn"              	BOOLEAN      	DEFAULT FALSE		NULL, 		-- 초기화비밀번호여부
	"active_datetime"                	TIMESTAMP    						NULL,     	-- 활성화일시
	"inactive_yn"                    	BOOLEAN      	DEFAULT FALSE		NULL, 		-- 비활성여부
	"state"                     		INTEGER      	DEFAULT 1			NULL,     	-- 사용자상태 (Code - UserStatus)
	"creator"                        	VARCHAR(30)  	DEFAULT 'system'	NOT NULL, 	-- 생성자
	"created_at"                     	TIMESTAMP    	DEFAULT NOW()		NULL,     	-- 생성일시
	"updater"                        	VARCHAR(30)  						NULL,     	-- 수정자
	"updated_at"                     	TIMESTAMP    						NULL      	-- 수정일시
)
WITH (
OIDS=false
);

-- 사용자 기본키2
CREATE UNIQUE INDEX "PK_tbl_user"
	ON "edgecraft"."tbl_user"
	( -- 사용자
		"user_uid" ASC -- 사용자식별자
	);

-- 사용자 유니크 인덱스
CREATE UNIQUE INDEX "UK_tbl_user"
	ON "edgecraft"."tbl_user"
	( -- 사용자
		"id" ASC, -- 사용자ID
		"email" ASC -- 이메일
	);

-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	ADD CONSTRAINT "PK_tbl_user"
		-- 사용자 기본키2
	PRIMARY KEY
	USING INDEX "PK_tbl_user"
	NOT DEFERRABLE;

-- 사용자
ALTER TABLE "edgecraft"."tbl_user"
	ADD CONSTRAINT "UK_tbl_user" -- 사용자 유니크 제약
	UNIQUE 
	USING INDEX "UK_tbl_user"
	NOT DEFERRABLE;