package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupGeneral(r *gin.Engine) {
	group := r.Group("/ai_arts/api/common")

	group.GET("/resource", wrapper(getResource))
}

type GetResourceReq struct {

}

type GetResourceRsp struct {
        AIFrameworkList         []models.AIFrameworkItem        `json:"ai_framework_list"`
        DeviceList                      []models.DeviceItem             `json:"device_list"`
}

// @Summary get available resource
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/resource [get]
func getResource(c *gin.Context) error {

	framework, devices, err := services.GetResource()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{
		framework,
		devices,
	}
	return SuccessRsp(c, rsp)
}
