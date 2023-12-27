package route

import (
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/middleware"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/server"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	echo_swagger "github.com/swaggo/echo-swagger"
)

// SetRoutes - api routes setting
func SetRoutes(api *api.API, server *server.Instance) {
	setSwagger(api, server)
	setHTTPRoutes(api, server)
}

func setSwagger(api *api.API, server *server.Instance) {
	server.HTTPServer.GET("/swagger/*", echo_swagger.EchoWrapHandler(echo_swagger.URL("doc.json")))
}

func setHTTPRoutes(api *api.API, server *server.Instance) {
	server.HTTPServer.Use(middleware.CustomCORS())

	server.HTTPServer.Use(middleware.CustomLogger())

	// create a sub route group
	v1 := server.HTTPServer.Group("/api/v1")

	// Middleware, default
	v1.Use(echo_middleware.Recover())

	// Session Interceptor
	// v1.Use(middleware.SessionInterceptor)

	v1.OPTIONS("", middleware.Preflight)

	// format-swagger:route [method] [path pattern] [?tag1 tag2 tag3] [operation id]
	v1.GET("/health", middleware.HealthCheck)

	// Common - CodeGroup
	v1.GET("/codes/groups", api.GetCodeGroupListHandler)
	v1.GET("/codes/groups/:groupId", api.GetCodeGroupHandler)
	v1.POST("/codes/groups", api.SetCodeGroupHandler)
	v1.PUT("/codes/groups", api.UpdateCodeGroupHandler)
	v1.DELETE("/codes/groups/:groupId", api.DeleteCodeGroupHandler)

	// Common - Code
	v1.GET("/codes", api.GetCodeListHandler)
	v1.GET("/codes/:groupId", api.GetCodeListByGroupHandler)
	v1.GET("/codes/:groupId/:code", api.GetCodeHandler)
	v1.POST("/codes", api.SetCodeHandler)
	v1.PUT("/codes", api.UpdateCodeHandler)
	v1.DELETE("/codes/:groupId/:code", api.DeleteCodeHandler)

	// Cloud
	v1.GET("/clouds", api.GetCloudListHandler)
	v1.GET("/clouds/:cloudId", api.GetCloudHandler)
	v1.POST("/clouds", api.SetCloudHandler)
	v1.PUT("/clouds/:cloudId", api.UpdateCloudHandler)
	v1.DELETE("/clouds/:cloudId", api.DeleteCloudHandler)

	// Cloud - Node
	v1.GET("/clouds/:cloudId/nodes", api.GetCloudNodeListHandler)
	v1.GET("/clouds/:cloudId/nodes/:nodeId", api.GetCloudNodeHandler)
	v1.POST("/clouds/:cloudId/nodes", api.SetCloudNodeHandler)
	v1.PUT("/clouds/:cloudId/nodes/:nodeId", api.UpdateCloudNodeHandler)
	v1.DELETE("/clouds/:cloudId/nodes/:nodeId", api.DeleteCloudNodeHandler)

	// Cloud - App
	// v1.GET("/clouds/:cloudID/apps", api.GetCloudAppHandler)
	// v1.GET("/clouds/:cloudID/apps/:appID", api.GetCloudAppHandler)
	// v1.POST("/clouds/:cloudID/apps", api.AddCloudAppHandler)
	// v1.PUT("/clouds/:cloudID/apps/:appID", api.UpdateCloudAppHandler)
	// v1.DELETE("/clouds/:cloudID/apps/:appID", api.DeleteCloudApphandler)
	// v1.GET("/template/apps", api.GetCloudAvailableAppsHandler)

	// Cloud - Security Verification
	// v1.GET("/clouds/:cloudID/securities", api.GetCloudVCResultsHandler)
	// v1.GET("/clouds/:cloudID/securities/:securityID", api.GetCloudVCResultByIDHandler)
	// v1.POST("/clouds/:cloudID/securities", api.SetCloudVCPeriodHandler)
	// TODO: 추가/삭제/갱신 기능은??

	// Cloud - K8S
	// v1.GET("/clouds/:cloudID/k8s/status", api.GetCloudK8sStatusHandler)
	// v1.POST("/clouds/:cloudID/k8s/:Version", api.UpgradeCloudK8sVersionHandler)
	// v1.GET("/clouds/:cloudID/kore-board", api.GetCloudKoreBoardHandler)

	// Openstack Cluster
	v1.GET("/clouds/:cloudId/clusters", api.GetClusterListHandler)                               // 클러스터 목록 조회
	v1.GET("/clouds/:cloudId/clusters/:clusterId", api.GetClusterHandler)                        // 클러스터 상세 기본 정보 조회
	v1.POST("/clouds/:cloudId/clusters", api.SetClusterHandler)                                  // 클러스터 등록/생성
	v1.PUT("/clouds/:cloudId/clusters/:clusterId", api.UpdateClusterHandler)                     // 클러스터 등록 정보 수정
	v1.DELETE("/clouds/:cloudId/clusters/:clusterId", api.DeleteClusterHandler)                  // 클러스터 삭제
	v1.POST("/clouds/:cloudId/clusters/:clusterId", api.ProvisioningClusterHandler)              // 클러스터 Provisioning
	v1.POST("/clouds/:cloudId/clusters/:clusterId/upgrade", api.UpgradeClusterK8sVersionHandler) // 클러스터 K8s 업그레이드

	// Openstack Cluster - NodeSet
	v1.GET("/clouds/:cloudId/clusters/:clusterId/nodesets", api.GetNodeSetListHandler)        // 클러스터 노드셋 목록 조회
	v1.GET("/clouds/:cloudId/clusters/:clusterId/nodesets/:nodeSetId", api.GetNodeSetHandler) // 클러스터 노드셋 상세 조회
	v1.POST("/clouds/:cloudId/clusters/:clusterId/nodesets", api.SetNodeSetHandler)           // 클러스터 노드셋 추가
	//v1.PUT("/clouds/:cloudId/clusters/:clusterId/nodesets/:nodeSetId", api.UpdateNodeSetHandler)              // 클러스터 노드셋 수정
	v1.DELETE("/clouds/:cloudId/clusters/:clusterId/nodesets/:nodeSetId", api.DeleteNodeSetHandler)           // 클러스터 노드셋 삭제
	v1.GET("/clouds/:cloudId/clusters/:clusterId/nodesets/:nodeSetId/:nodeCount", api.UpdateNodeCountHandler) // 클러스터 노드셋의 노드 카운트 갱신

	// Cloud/Cluster - App
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/apps", api.GetCloudClusterAppsHandler)                   // 클러스터 상세 애플리케이션 목록 조회
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/apps/:appID", api.GetCloudClusterAppHandler)             // 클러스터 상세 애플리케이션 정보 조회
	// v1.POST("/clouds/:cloudID/clusters/:clusterID/apps", api.AddCloudClusterAppsHandler)                  // 클러스터 상세 애플리케이션 추가
	// v1.PUT("/clouds/:cloudID/clusters/:clusterID/apps/:appID", api.UpdateCloudClusterAppsHandler)         // 클러스터 상세 애플리케이션 업데이트
	// v1.DELETE("/clouds/:cloudID/clusters/:clusterID/apps/:appID", api.DeleteCloudClusterAppsHandler)      // 클러스터 상세 애플리케이션 삭제
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/template/apps", api.GetCloudClusterAvailableAppsHandler) // 클러스터에 설치 가능한 애플리케이션 목록

	// Openstack Cluster - CIS Benchmarks
	v1.POST("/clouds/:cloudId/clusters/:clusterId/benchmarks", api.SetBenchmarksHandler)              // 클러스터 벤치마크 실행
	v1.GET("/clouds/:cloudId/clusters/:clusterId/benchmarks", api.GetBenchmarksListHandler)           // 클러스터 벤치마크 결과 목록 조회
	v1.GET("/clouds/:cloudId/clusters/:clusterId/benchmarks/:benchmarksId", api.GetBenchmarksHandler) // 클러스터 벤치마크 결과 상세 조회

	// Cloud/Cluster - Backup
	v1.POST("/clouds/:cloudId/clusters/:clusterId/backup", api.SetBackupHandler)                 // 클러스터 백업
	v1.GET("/clouds/:cloudId/clusters/:clusterId/backup", api.GetBackupListHandler)              // 클러스터 백업 목록 조회
	v1.DELETE("/clouds/:cloudId/clusters/:clusterId/backup/:backresId", api.DeleteBackupHandler) // 클러스터 백업 삭제

	// Cloud/Cluster - Restore
	v1.POST("/clouds/:cloudId/clusters/:clusterId/restore", api.SetRestoreHandler)                 // 클러스터 복원
	v1.GET("/clouds/:cloudId/clusters/:clusterId/restore", api.GetRestoreListHandler)              // 클러스터 복원 목록 조회
	v1.DELETE("/clouds/:cloudId/clusters/:clusterId/restore/:backresId", api.DeleteRestoreHandler) // 클러스터 복원 삭제

	// Cloud/Cluster - Security Verification
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/securities", api.GetCloudClusterVCResultsHandler)            // 클러스터 보안검증 결과 목록 조회
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/securities/:securityID", api.GetCloudClusterVCResultHandler) // 클러스터 보안검증 결과 상세 정보 조회
	// v1.POST("/clouds/:cloudID/clusters/:clusterID/securities", api.SetCloudClusterVCPeriodHandler)            // 클러스터 보안검증 주기 설정

	// Cloud/Cluster - K8S
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/k8s", api.GetCloudClusterK8sStatusHandler)                // 클라우드 상세 k8s 정보 조회
	// v1.POST("/clouds/:cloudID/clusters/:clusterID/k8s/:Version", api.UpgradeCloudClusterK8sVersionHandler) // 클러스터 상세 쿠버네티스 버전 업그레이드
	// v1.GET("/clouds/:cloudID/clusters/:clusterID/kore-board", api.GetCloudClusterKoreBoardHandler)         // 클러스터 상세 kore-board 연계

	// Image
	// v1.GET("/images", api.GetImagesHandler)               // 이미지 목록 조회
	// v1.GET("/images/:imageID", api.GetImageHandler)       // 이미지 상세 기본 정보 조회
	// v1.POST("/images", api.AddImageHandler)               // 이미지 등록(업로드)
	// v1.PUT("/images/:imageID", api.UpdateImageHandler)    // 이미지 등록 정보 수정 및 업데이트
	// v1.DELETE("/images/:imageID", api.DeleteImageHandler) // 이미지 삭제

	// Security Verification
	// v1.GET("/securities/:Version", api.GetVCByVersionHandler)   // 보안검증항목 조회
	// v1.POST("/securities", api.AddVCHandler)                     // 보안검증항목 등록
	// v1.PUT("/securities/:Version", api.UpdateVCByVersionHandler) // 보안검증항목 등록 정보 업데이트

	// Accouts
	v1.GET("/users", api.GetUserListHandler)           // 사용자 목록 조회
	v1.GET("/users/:userId", api.GetUserHandler)       // 사용자 상세 기본 정보 조회
	v1.POST("/users", api.SetUserHandler)              // 사용자 등록
	v1.PUT("/users/:userId", api.UpdateUserHandler)    // 사용자 정보 수정
	v1.DELETE("/users/:userId", api.DeleteUserHandler) // 사용자 삭제

	// Auth
	v1.POST("/auth", api.LoginHandler)         // 사용자 로그인
	v1.POST("/auth/logout", api.LogoutHandler) // 사용자 로그아웃
}
