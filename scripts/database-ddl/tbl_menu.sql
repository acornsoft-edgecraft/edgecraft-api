CREATE TABLE tbl_menu
(
	-- 메뉴uid : 메뉴순번
	menu_uid serial NOT NULL UNIQUE,
	-- 메뉴아이디 : 메뉴ID
	menu_id varchar(20) NOT NULL,
	-- 이름 : 한글메뉴명
	name varchar(100) NOT NULL,
	-- 메뉴타입 : 메뉴타입: BACKEND-운영자사이트, FRONTEND-사용자사이트
	menu_type varchar(20) NOT NULL,
	-- 상위메뉴 ID
	up_menu_cd varchar(20) NOT NULL,
	-- 하위메뉴 존재여부 : Y(존재함), N(존재안함)
	sub_menu_yn varchar(1) NOT NULL,
	-- 페이지url
	page_url varchar(100),
	-- 메뉴출력순서 : 매뉴 출력 순서, 각 단계별 3자리 순자로 구분, 예) 100(게시판관리), 100001(공지사항관리), 100001001(공지사항 추가)
	menu_order int NOT NULL,
	-- 메뉴 깊이 : 0(HOME), 1(1 depth), 2(2 depth)
	menu_depth int NOT NULL,
	-- 로그인체크여부 : 메뉴 선택시 로그인의 필요여부: Y(로그인필요), N(로그인 필요없음)
	login_chk_yn varchar(1) DEFAULT 'N' NOT NULL,
	-- 사용여부 : 사용여부
	use_yn varchar(1) DEFAULT 'Y' NOT NULL,
	-- 메뉴표시여부 : 메뉴표시여부: Y(표시), N(표시안함)
	menu_disp_yn varchar(1) NOT NULL,
	-- 메뉴노출타입 : 메뉴노출위치
	menu_sub_type varchar(3) DEFAULT 'N' NOT NULL,
	-- 메뉴아이콘 : 메뉴 좌측 아이콘
	menu_icon varchar(100),
	-- 등록자uid : 등록자
	reg_user_uid int,
	-- 등록일 : 등록일시
	reg_date timestamp DEFAULT current_timestamp NOT NULL,
	-- 수정자uid : 수정자
	edt_user_uid int,
	-- 수정일 : 수정일시
	edt_date timestamp,
	PRIMARY KEY (menu_id)
) WITHOUT OIDS;