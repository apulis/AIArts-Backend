package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupUpdatePlatform(r *gin.Engine) {
	group := r.Group("/ai_arts/api/updatePlatform")

	// group.Use(Auth())

	group.GET("/version-info", wrapper(getVersionInfo))
	group.GET("/version-detail/:id", wrapper(getVersionDetailByID))
	group.GET("/online-upgrade-progress", wrapper(getOnlineUpgradeProgress))
	group.GET("/local-upgrade-progress", wrapper(getLocalUpgradeProgress))
	group.GET("/local-upgrade-env-check", wrapper(checkLocalEnv))
	group.POST("/online-upgrade", wrapper(upgradeOnline))
	group.POST("/local-upgrade", wrapper(upgradeLocal))
}

type getVersionInfoResp struct {
	VersionInfo []models.VersionInfoSet `json:"versionInfoLogs"`
}

type getVersionInfoReq struct {
	queryLimit int `form:"limit,default=10"`
}

type getLocalEnvResp struct {
	CanUpgrade bool `json:"canUpgrade"`
	IsLower    bool `json:"isLower"`
}

type getLocalUpgradeProgressResp struct {
	Status  string `json:"status"`
	Percent int    `json:"percent"`
}

// @Summary get version infomation
// @Produce  json
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/updatePlatform/version-info [get]
func getVersionInfo(c *gin.Context) error {

	versionlogs, err := services.GetVersionLogs()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := getVersionInfoResp{
		VersionInfo: versionlogs,
	}
	return SuccessResp(c, data)

}

func getVersionDetailByID(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

func getOnlineUpgradeProgress(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

var progress int
var status string

// @Summary get local upgrade process
// @Produce  json
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/updatePlatform/local-upgrade-progress [get]
func getLocalUpgradeProgress(c *gin.Context) error {
	if status == "" {
		status = "upgrading"
	}
	data := getLocalUpgradeProgressResp{
		Status:  status,
		Percent: progress,
	}
	progress += 10
	if progress > 100 {
		progress = 0
		status = "upgrading"
	}
	if progress == 100 {
		status = "finish"
	}
	return SuccessResp(c, data)
}

// @Summary get local upgrade environment info
// @Produce  json
// @Success 200 {object} APISuccessRespGetDataset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/updatePlatform/local-upgrade-env-check [get]
func checkLocalEnv(c *gin.Context) error {
	canUpgrade, isLower, err := services.GetLocalUpgradeEnv()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := getLocalEnvResp{
		CanUpgrade: canUpgrade,
		IsLower:    isLower,
	}
	return SuccessResp(c, data)
}

func upgradeOnline(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

// @Summary upgrade through local package
// @Produce  json
// @Success 200 {object} APISuccessRespGetDataset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/updatePlatform/local-upgrade [post]
func upgradeLocal(c *gin.Context) error {
	err := services.UploadUpgradeInfo()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
