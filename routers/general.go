package routers

import (
	"encoding/json"
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
	group.GET("/resources", wrapper(getResources))
	group.GET("/job/summary", wrapper(getJobSummary))
	group.DELETE("/DeleteJob", wrapper(DeleteJob))
}

type GetResourceReq struct {
	VCName string `form:"vcName" json:"vcName"`
}

type GetResourceRsp struct {
	AIFrameworks          map[string][]string  `json:"aiFrameworks"`
	DeviceList            []models.DeviceItem  `json:"deviceList"`               //
	NodeInfo              []*models.NodeStatus `json:"nodeInfo"`
	CodePathPrefix        string               `json:"codePathPrefix"`
	NodeCountByDeviceType map[string]int       `json:"nodeCountByDeviceType"`
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

	var req GetResourceReq
	var err error

	if err = c.ShouldBindQuery(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	// 兼容老代码
	if req.VCName == "" {
		req.VCName = models.DefaultVcName
	}

	vcInfo, err := services.GetResource(userName, req.VCName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{}
	rsp.AIFrameworks = make(map[string][]string)

	// 获取平台镜像列表
	for k, v := range configs.Config.Image {

		rsp.AIFrameworks[k] = make([]string, 0)
		for _, item := range v {
			rsp.AIFrameworks[k] = append(rsp.AIFrameworks[k], item)
		}
	}

	// 获取平台配额数据
	quota := make(map[string]int)
	if len(vcInfo.Quota) != 0 {
		err = json.Unmarshal([]byte(vcInfo.Quota), &quota)
		if err != nil {

		}
	}

	// 获取设备列表
	rsp.DeviceList = make([]models.DeviceItem, 0)
	for k, v := range vcInfo.DeviceAvail {

		// 计算设备列表
		deviceInfo := models.DeviceItem{
			DeviceType: k,
			Avail:      v,
		}

		if _, ok := quota[k]; ok {
			deviceInfo.UserQuota = quota[k]
		}

		rsp.DeviceList = append(rsp.DeviceList, deviceInfo)
	}

	// 统计设备类型的节点数
	rsp.NodeCountByDeviceType = make(map[string]int)
	for _, v := range vcInfo.Nodes {
		rsp.NodeCountByDeviceType[v.GPUType] = rsp.NodeCountByDeviceType[v.GPUType] + 1
	}

	// 获取节点信息
	for _, v := range vcInfo.Nodes {
		if len(v.GPUType) != 0 {
			rsp.NodeInfo = append(rsp.NodeInfo, v)
		}
	}

	// 代码环境中的 用户home路径
	rsp.CodePathPrefix = "/home/" + userName + "/"
	return SuccessResp(c, rsp)
}

// @Summary get available resource
// @Produce  json
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/resource [get]
func getResources(c *gin.Context) error {
	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	Resources, err := services.GetResources(userName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"resources": Resources})
}

// @Summary get job summary
// @Produce  json
// @Param param body GetJobSummaryReq true "params"
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/job/summary [get]
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

	summary, err := services.GetJobSummary(userName, req.JobType)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, summary)
}
