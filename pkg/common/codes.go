// Package common - Defines common constants and variables for messages
package common

// ===== [ Constants and Variables ] =====

const (
	//############################################################
	// 상태코드 (공통코드 연계용)
	//############################################################
	StatusSaved        = 1
	StatusProvisioning = 2
	StatusProvisioned  = 3
	StatusFailed       = 4
	StatusDeleting     = 5
	StatusDeleted      = 6

	//############################################################
	// NodeType (공통코드 연계용)
	//############################################################
	NodeTypeMaster = 1
	NodeTypeWorker = 2

	//############################################################
	// UserStatus (공통코드 연계용)
	//############################################################
	UserStatusActivated = 1

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

	// DatabaseCodeFalseData
	DatabaseCodeFalseData = 24100

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
	// Cloud/Cluster 관련 메시지
	//############################################################

	// CloudNotFound -
	CloudNotFound = 25000
	// ClusterNotFound -
	ClusterNotFound = 25001
	// NodeNotFound -
	NodeNotFound = 25002
	// NodeSetNotFound -
	NodeSetNotFound = 25003

	// CreatedCloudNoUpdatable -
	CreatedCloudNoUpdatable = 25100
	// CreatedClusterNoUpdatable -
	CreatedClusterNoUpdatable = 25100

	// ProvisioningOnlySavedOrDeleted
	ProvisioningOnlySavedOrDeleted = 26000
	// OpenstackClusterRegistered
	OpenstackClusterRegistered = 26001
	// OpenstackClusterProvisioning
	OpenstackClusterProvisioning = 26002
	// OpenstackClusterDeleting
	OpenstackClusterDeleting = 26003
	// OpenstackClusterInfoDeleted
	OpenstackClusterInfoDeleted = 26004
	// OpenstackProvisionDeleted
	OpenstackProvisionDeleted = 26005
	// OpenstsackClusterAlreadyDeleting
	OpenstsackClusterAlreadyDeleting = 26006

	// ProvisioningCheckJobFailed -
	ProvisioningCheckJobFailed = 26100
	// DeleteProvisionedClusterJobFailed -
	DeleteProvisionedClusterJobFailed = 26101

	//############################################################
	// K8S 관련 메시지
	//############################################################

	// CodeFailedK8SAPI -
	CodeFailedK8SAPI = 27000
	// KubernetesNotYet -
	KubernetesNotYet = 27001

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
