package routers

import (
	"os"

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
	VersionInfo models.VersionInfoSet   `json:"versionInfo"`
	VersionLogs []models.VersionInfoSet `json:"versionLogs"`
}

type getVersionInfoReq struct {
	queryLimit int `form:"limit,default=10"`
}

// @Summary get version infomation
// @Produce  json
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /version-info [get]
func getVersionInfo(c *gin.Context) error {

	curenVersion, err := services.GetCurrentVersion()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	versionlogs, err := services.GetVersionLogs()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := getVersionInfoResp{
		VersionInfo: curenVersion,
		VersionLogs: versionlogs,
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

func getLocalUpgradeProgress(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

func checkLocalEnv(c *gin.Context) error {
	// isUpgradeFileExist := isFileExists("/data/DLTSUpgrade")

	data := "test"
	return SuccessResp(c, data)
}

func upgradeOnline(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

func upgradeLocal(c *gin.Context) error {
	err := services.UploadUpgradeInfo()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

func isFileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
