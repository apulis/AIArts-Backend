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
	AIFrameworks   map[string][]string `json:"aiFrameworks"`
	DeviceList     []models.DeviceItem `json:"deviceList"`
	CodePathPrefix string              `json:"codePathPrefix"`
}

// @Summary get available resource
// @Produce  json
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/resource [get]
func getResource(c *gin.Context) error {

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	framework, devices, err := services.GetResource(userName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{
		framework,
		devices,
		"/home/" + userName + "/",
	}

	return SuccessResp(c, rsp)
}

func getUsername(c *gin.Context) string {

	data, exists := c.Get("userName")
	if !exists {
		return ""
	}

	userName := fmt.Sprintf("%v", data)
	return userName
}

func getUserId(c *gin.Context) string {

	data, exists := c.Get("userId")
	if !exists {
		return ""
	}

	userId := fmt.Sprintf("%v", data)
	return userId
}