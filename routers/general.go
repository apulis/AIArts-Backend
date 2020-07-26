package routers

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupGeneral(r *gin.Engine) {
	group := r.Group("/ai_arts/api/common")

	group.Use(Auth())

	group.GET("/resource", wrapper(getResource))
	group.GET("/job/summary", wrapper(getJobSummary))
}

type GetResourceReq struct {
}

type GetResourceRsp struct {
	AIFrameworks   map[string][]string `json:"aiFrameworks"`
	DeviceList     []models.DeviceItem `json:"deviceList"`
	NodeInfo       models.NodeInfo     `json:"nodeInfo"`
	CodePathPrefix string              `json:"codePathPrefix"`
}

type GetJobSummaryReq struct {
	JobType string `form:"jobType" json:"jobType"`
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

	vcInfo, err := services.GetResource(userName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{}
	rsp.AIFrameworks = make(map[string][]string)

	for k, v := range configs.Config.Image {

		rsp.AIFrameworks[k] = make([]string, 0)
		for _, item := range v {
			rsp.AIFrameworks[k] = append(rsp.AIFrameworks[k], item)
		}
	}

	rsp.DeviceList = make([]models.DeviceItem, 0)
	for k, v := range vcInfo.DeviceAvail {
		rsp.DeviceList = append(rsp.DeviceList, models.DeviceItem{
			DeviceType: k,
			Avail:      v,
		})
	}

	rsp.NodeInfo.TotalNodes = len(vcInfo.Nodes)
	rsp.NodeInfo.CountByDeviceType = make(map[string]int)

	for _, v := range vcInfo.Nodes {
		rsp.NodeInfo.CountByDeviceType[v.GPUType] += 1
	}

	rsp.CodePathPrefix = "/home/" + userName + "/"

	return SuccessResp(c, rsp)
}

// @Summary get job summary
// @Produce  json
// @Param param body UpdateTemplateReq true "params"
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/jobs/summary [get]
func getJobSummary(c *gin.Context) error {

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	var err error
	var req GetJobSummaryReq

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	vcInfo, err := services.GetJobSummary(userName, req.JobType)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{}
	rsp.AIFrameworks = make(map[string][]string)

	rsp.NodeInfo.TotalNodes = len(vcInfo.Nodes)
	rsp.NodeInfo.CountByDeviceType = make(map[string]int)

	for _, v := range vcInfo.Nodes {
		rsp.NodeInfo.CountByDeviceType[v.GPUType] += 1
	}

	rsp.CodePathPrefix = "/home/" + userName + "/"

	return SuccessResp(c, rsp)
}
