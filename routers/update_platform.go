package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupUpdatePlatform(r *gin.Engine) {
	group := r.Group("/ai_arts/api/version")

	// group.Use(Auth())

	group.GET("/info", wrapper(getVersionInfo))
	group.GET("/detail/:id", wrapper(getVersionDetailByID))
	group.GET("/upgradeProgress", wrapper(getLocalUpgradeProgress))
	group.GET("/env/local", wrapper(checkLocalEnv))
	group.POST("/upgrade/online", wrapper(upgradeOnline))
	group.POST("/upgrade/local", wrapper(upgradeLocal))
}

type getVersionInfoResp struct {
	CurrentVersion models.VersionInfoSet   `json:"versionInfo"`
	VersionInfo    []models.VersionInfoSet `json:"versionLogs"`
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
// @Success 200 {object} getVersionInfoResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/version/info [get]
func getVersionInfo(c *gin.Context) error {
	currentversion, err := services.GetCurrentVersion()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	versionlogs, err := services.GetVersionLogs()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := getVersionInfoResp{
		CurrentVersion: currentversion,
		VersionInfo:    versionlogs,
	}
	return SuccessResp(c, data)

}

func getVersionDetailByID(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

var progress int
var status string

// @Summary get local upgrade process
// @Produce  json
// @Success 200 {object} getLocalUpgradeProgressResp
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/version/upgradeProgress [get]
func getLocalUpgradeProgress(c *gin.Context) error {
	status, progress := services.GetUpgradeProgress()
	data := getLocalUpgradeProgressResp{
		Status:  status,
		Percent: progress,
	}
	return SuccessResp(c, data)
}

// @Summary get local upgrade environment info
// @Produce  json
// @Success 200 {object} getLocalEnvResp
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/version/env/local [get]
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
// @Success 200 {object} APISuccessRespGetDataset
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/version/upgrade/local [post]
func upgradeLocal(c *gin.Context) error {
	err := services.UpgradePlatformByLocal()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
