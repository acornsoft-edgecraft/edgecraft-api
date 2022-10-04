// Package common - Defines common constants and variables for messages
package common

// ===== [ Constants and Variables ] =====

const (
	//############################################################
	// 일반 메시지
	//############################################################

	// CodeOK -
	CodeOK = 20000

	//############################################################
	// 일반 오류 메시지
	//############################################################

	// CodeProcessingError -
	CodeProcessingError = 21000
	// CodeInvalidData -
	CodeInvalidData = 21001
	// CodeFailedSaveData -
	CodeFailedSaveData = 21002
	// CodeDuplicatedData -
	CodeDuplicatedData = 21003

	// CodeInvalidParm -
	CodeInvalidParm = 21100

	//############################################################
	// 데이터베이스 오류 메시지 - 24000
	//############################################################

	// DB fail -
	CodeFailedDatabase = 24000

	// DatabaseFailedRollback -
	DatabaseFailedRollback = 24100

	// DatabaseEmptyData -
	DatabaseEmptyData = 24200
	// DatabaseFalseData -
	DatabaseFalseData = 24201
	// DatabaseExistData -
	DatabaseExistData = 24202

	//############################################################
	// 인증/계정 관련 메시지
	//############################################################

	// CodeDuplicateEmail -
	CodeDuplicateEmail = 22000
	// CodeInvalideEmail -
	CodeInvalideEmail = 22001
	// CodeInvalidUser -
	CodeInvalidUser = 22002
	// CodeNotFoundUser -
	CodeNotFoundUser = 22003
	// CodeNotValidUser -
	CodeNotValidUser = 22004
	// CodeInitialPassword -
	CodeInitialPassword = 22005
	// CodeInvalidPassword -
	CodeInvalidPassword = 22006
	// CodeFaildStatusUser -
	CodeFaildStatusUser = 22007

	// CodeInvalidToken -
	CodeInvalidToken = 22100
	// CodeFailedCreateToken -
	CodeFailedCreateToken = 22101
	// CodeInvalidVerification -
	CodeInvalidVerification = 22102

	//############################################################
	// 메일 관련 메시지
	//############################################################

	// CodeFailedSendMail -
	CodeFailedSendMail = 23300

	//############################################################
	// 메일 관련 메시지
	//############################################################

	//############################################################
	// Cloud 관련 메시지
	//############################################################

	// CloudCreated_CantUpdate -
	CloudCreated_CantUpdate = 25001

	//############################################################
	// K8S 관련 메시지
	//############################################################

	// DB fail -
	CodeFailedK8SAPI = 27000

	//############################################################
	// 태그 관련 메시지
	//############################################################

	// GroupName Dup fail -
	GroupNameDupTagGroup = 28000

	//############################################################
	// 세션 관련 메시지
	//############################################################

	// Session fail -
	SessionNotFound = 29000
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// ===== [ Public Functions ] =====
