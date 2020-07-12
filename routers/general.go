package routers

import (
	"fmt"
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
        CodePathPrefix 			string `json:"codePathPrefix"`
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
		"/home/username",
	}
	return SuccessResp(c, rsp)
}

func getUsername(c *gin.Context) (string, error) {

	data, err := c.Get("userName")
	if err != nil {
		return nil, err
	}

	userName := fmt.Sprintf("%v", data)
	return userName, nil
}
