package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGrouSettings(r *gin.Engine) {

	group := r.Group("/ai_arts/api/settings")
	group.Use(Auth())

	group.POST("/privileged", wrapper(upsertPrivilegedSetting))
	group.GET("/privileged", wrapper(getPrivilegedSetting))
}

// @Summary update or insert privileged job settings
// @Produce json
// @Param privileged settings
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/settings/privileged [post]
func upsertPrivilegedSetting(c *gin.Context) {
	var req = models.PrivilegedSetting

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	var err = services.UpsertPrivilegedSetting(req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, gin.H{})
}

// @Summary get privileged job settings
// @Produce json
// @Param privileged settings
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/settings/privileged [get]
func getPrivilegedSetting(c *gin.Context) {
	settings, err = services.GetPrivilegedSetting()
	if err != nil {
		return ParameterError(err.Error())
	}

	return SuccessResp(c, settings)
}
