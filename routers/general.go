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
	group.GET("/images", wrapper(getImageList))
	group.GET("/job/summary", wrapper(getJobSummary))
	group.DELETE("/DeleteJob", wrapper(DeleteJob))
	group.GET("/jobs/:id/raw_log", wrapper(getJobRawLog))
}

type GetResourceReq struct {
	VCName string `form:"vcName" json:"vcName"`
}

type GetResourceRsp struct {
	AIFrameworks          map[string][]string  `json:"aiFrameworks"`
	DeviceList            []models.DeviceItem  `json:"deviceList"` //
	NodeInfo              []*models.NodeStatus `json:"nodeInfo"`
	CodePathPrefix        string               `json:"codePathPrefix"`
	NodeCountByDeviceType map[string]int       `json:"nodeCountByDeviceType"`
}

type GetJobSummaryReq struct {
	UserName string `form:"userName" json:"userName"`
	JobType  string `form:"jobType" json:"jobType"`
	VCName   string `form:"vcName" json:"vcName"`
}

type ImageItem struct {
	Image       			string `json:"image"`
	ImageType   			string `json:"imageType"`
	*models.ImageParams
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
		return AppError(configs.NO_USRNAME, "no username")
	}

	// 兼容老代码
	if req.VCName == "" {
		req.VCName = models.DefaultVcName
	}

	vcInfo, err := services.GetResource(userName, req.VCName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	rsp := &GetResourceRsp{}
	rsp.AIFrameworks = make(map[string][]string)

	// 获取平台镜像列表
	db_images, _, _ := models.ListImages(0, 1000)
	for _, v := range db_images {

		if _, ok := rsp.AIFrameworks[v.ImageType]; !ok {
			rsp.AIFrameworks[v.ImageType] = make([]string, 0)
		}

		rsp.AIFrameworks[v.ImageType] = append(rsp.AIFrameworks[v.ImageType], v.ImageFullName)
	}

	// 获取平台配额数据
	quota := make(map[string]int)                   // vc配额数据
	user_quota := make(map[string]models.UserQuota) // vc下用户配额数据
	user_used := make(map[string]int)

	if len(vcInfo.Quota) != 0 {
		err = json.Unmarshal([]byte(vcInfo.Quota), &quota)
		if err != nil {

		}
	}

	// 获取用户配额数据
	if len(vcInfo.Metadata) != 0 {
		err = json.Unmarshal([]byte(vcInfo.Metadata), &user_quota)
		if err != nil {

		}
	}

	// 整理用户已使用数据
	for _, v := range vcInfo.UserStatus {
		if v.UserName == userName {
			for dev, used := range v.UserGPU {
				user_used[dev] = used
			}
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

		// 1. 用户配置了user_quota，直接返回配置数据
		// 2. 用户未配置user_quota，返回VC配置数据
		if _, ok := user_quota[k]; ok {
			deviceInfo.UserQuota = user_quota[k].Quota
		} else if _, ok := quota[k]; ok {
			deviceInfo.UserQuota = quota[k]
		} else {
			deviceInfo.UserQuota = 0
		}

		// 用户配额 - 用户已使用设备数量
		if _, ok := user_used[k]; ok {
			deviceInfo.Avail = deviceInfo.UserQuota - user_used[k]
		} else {
			deviceInfo.Avail = deviceInfo.UserQuota
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
		return AppError(configs.NO_USRNAME, "no username")
	}
	Resources, err := services.GetResources(userName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
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
		return AppError(configs.NO_USRNAME, "no username")
	}

	var err error
	var req GetJobSummaryReq

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	summary, err := services.GetJobSummary(userName, req.JobType, req.VCName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, summary)
}

// @Summary get image list
// @Produce  json
// @Param param body GetJobSummaryReq true "params"
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/images/ [get]
func getImageList(c *gin.Context) error {

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	images := make([]ImageItem, 0)
	db_images, _, _ := models.ListImages(0, 1000)

	for _, v := range db_images {

		item := ImageItem{
			Image:    v.ImageFullName,
			ImageParams: new(models.ImageParams),
		}

		item.Desc = v.Details.Desc
		item.Category = v.Details.Category

		images = append(images, item)
	}

	return SuccessResp(c, images)
}

// @Summary get job raw log
// @Produce  json
// @Param jobId uri string "job id"
// @Success 200 {object} APISuccessRespGetResource "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/jobs/:id/raw_log [get]
func getJobRawLog(c *gin.Context) error {

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	id := c.Param("id")
	if len(id) == 0 {
		return AppError(configs.PARAMETER_ERROR_CODE, "job id invalid")
	}

	jobRawLog, err := services.GetJobRawLog(userName, id)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, jobRawLog)
}
