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
        AIFrameworks            map[string][]string `json:"aiFrameworks"`
        DeviceList              []models.DeviceItem `json:"deviceList"`
}

// @Summary get available resource
// @Produce  json
// @Success 200 {object} APISuccessRespGetResource "success"
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
	return SuccessResp(c, rsp)
}
