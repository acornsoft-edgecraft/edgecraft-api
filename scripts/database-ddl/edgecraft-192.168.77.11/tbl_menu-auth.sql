-- 메뉴권한
CREATE TABLE tbl_menu_auth
(
	-- auth_uid
	auth_uid serial NOT NULL UNIQUE,
	-- 메뉴아이디 : 메뉴ID
	menu_id varchar(20) NOT NULL,
	-- 사용자롤id : 사용자롤id
	user_role_id varchar(8) NOT NULL,
	-- rw 여부 : 공통코드 사용 (ex. READONLY, READWRITE)
	attr_rw varchar(50),
	-- 사용여부 : 사용여부
	use_yn varchar(1) NOT NULL,
	-- 등록자uid : 등록자
	reg_user_uid int,
	-- 등록일 : 등록일시
	reg_date timestamp DEFAULT current_timestamp NOT NULL,
	-- 수정자uid : 수정자
	edt_user_uid int,
	-- 수정일 : 수정일시
	edt_date timestamp,
	PRIMARY KEY (menu_id, user_role_id)
) WITHOUT OIDS;