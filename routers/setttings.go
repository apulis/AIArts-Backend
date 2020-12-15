package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGrouSettings(r *gin.Engine) {

	group := r.Group("/ai_arts/api/settings")
	group.Use(Auth())

	group.POST("/privileged", wrapper(upsertPrivilegedSettings))
	group.GET("/privileged", wrapper(getPrivilegedSettings))
}

// @Summary update or insert privileged job settings
// @Produce json
// @Param privileged settings
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/settings/privileged [post]
func upsertPrivilegedSettings(c *gin.Context) {
	var req = models.PrivilegedSettings

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	var err = services.UpsertPrivilegedSettings(req)
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
func getPrivilegedSettings(c *gin.Context) {
	settings, err = services.GetPrivilegedSettings()
	if err != nil {
		return ParameterError(err.Error())
	}

	return SuccessResp(c, settings)
}
