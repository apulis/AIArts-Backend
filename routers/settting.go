package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGrouSetting(r *gin.Engine) {

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
func upsertPrivilegedSetting(c *gin.Context) error {
	var req models.PrivilegedSetting

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	token := c.GetHeader("Authorization")
	hasPermission, err := HasPermission(token, "MANAGE_PRIVILEGE_JOB")
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	if !hasPermission {
		return AppError(configs.OPERATION_FORBIDDEN, "operation forbidden")
	}

	err = services.UpsertPrivilegedSetting(req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
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
func getPrivilegedSetting(c *gin.Context) error {
	setting, err := services.GetPrivilegedSetting()
	if err != nil {
		return ParameterError(err.Error())
	}

	return SuccessResp(c, setting)
}
