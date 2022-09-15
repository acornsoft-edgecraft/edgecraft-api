package postgresdb

// 2	user_role	int2	YES	NULL	NULL		NULL
// 3	user_name	varchar(50)	YES	NULL	NULL		NULL
// 4	user_id	varchar(50)	NO	NULL	NULL		NULL
// 5	password	varchar(100)	YES	NULL	NULL		NULL
// 6	email	varchar(40)	YES	NULL	NULL		NULL
// 7	last_login	timestamp	YES	NULL	NULL		NULL
// 8	password_expiration_begin_time	timestamp	NO	NULL	NULL		NULL
// 9	password_expiration_end_time	timestamp	NO	NULL	NULL		NULL
// 10	reset_password_yn	bpchar(1)	YES	NULL	"'N'::bpchar"		NULL
// 11	active_datetime	timestamp	YES	NULL	NULL		NULL
// 12	inactive_yn	bpchar(1)	YES	NULL	"'N'::bpchar"		NULL
// 13	user_state	bpchar(1)	YES	NULL	NULL		NULL
// 14	creator	varchar(50)	NO	NULL	NULL		NULL
// 15	created	timestamp	YES	NULL	NULL		NULL
// 16	updater	varchar(50)	YES	NULL	NULL		NULL
// 17	updated	timestamp	YES	NULL	NULL		NULL
// 18	user_uid	uuid	NO	NULL	uuid_generate_v4()		NULL
