package routers

import (
	"github.com/gin-gonic/gin"
)

func AddGroupUpdatePlatform(r *gin.Engine) {
	group := r.Group("/ai_arts/api/updatePlatform")

	group.Use(Auth())

	group.GET("/version-info", wrapper(getVersionInfo))
	group.GET("/version-detail/:id", wrapper(getVersionDetailByID))
	group.GET("/online-upgrade-progress", wrapper(getOnlineUpgradeProgress))
	group.GET("/local-upgrade-progress", wrapper(getLocalUpgradeProgress))
	group.GET("/local-upgrade-env-check", wrapper(checkLocalEnv))
	group.POST("/online-upgrade", wrapper(upgradeOnline))
	group.POST("/local-upgrade", wrapper(upgradeLocal))
}

// @Summary acquire version infomation, including the current version and dated ones
// @Produce  json
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /version-info [get]
func getVersionInfo(c *gin.Context) error {
	data := "test"
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
	data := "test"
	return SuccessResp(c, data)
}

func upgradeOnline(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}

func upgradeLocal(c *gin.Context) error {
	data := "test"
	return SuccessResp(c, data)
}
