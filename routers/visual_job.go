package routers

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
)

func AddGroupVisualJob(r *gin.Engine) {
	group := r.Group("/ai_arts/api/visual")

	group.Use(Auth())

	group.POST("/", wrapper(createVisualJob))
	group.GET("/list", wrapper(getVisualJobList))
	group.GET("/endpoints", wrapper(getEndpoints))
	group.DELETE("/", wrapper(deleteVisualJob))
	group.PUT("/", wrapper(switchVisualJobStatus)) //
}

type CreateVisualJobReq struct {
	JobName           string `form:"jobName"`
	TensorboardLogDir string `form:"tensorboardLogDir"`
	Description       string `form:"description"`
}

type GetVisualJobListReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	OrderBy  string `form:"orderBy"`
	Status   string `form:"status"`
	JobName  string `form:"jobName"`
	Order    string `form:"order"`
}

type GetVisualJobListRsq struct {
	Templates    []VisualJobListRspUnit `json:"Templates"`
	TotalJobsNum int                    `json:"total"`
	TotalPages   int                    `json:"totalPages"`
}

type VisualJobListRspUnit struct {
	Id                int    `json:"id"`
	JobName           string `json:"jobName"`
	Status            string `json:"status"`
	CreateTime        models.UnixTime `json:"createTime"`
	TensorboardLogDir string `json:"TensorboardLogDir"`
	Description       string `json:"description"`
}

type GetEndpointsReq struct {
	JobId int `form:"id"`
}
type GetEndpointsRsq struct {
	Path string `json:"path"`
}

type SwitchVisualJobStatusReq struct {
	JobId  int    `form:"id"`
	Status string `form:"status"`
}

type DeleteJobReq struct {
	JobId int `form:"id" binding:"required"`
}

// @Summary create visual job
// @Produce json
// @Param param body CreateVisualJobReq true "params"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual [post]
func createVisualJob(c *gin.Context) error {
	var req CreateVisualJobReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	err = services.CreateVisualJob(username, req.JobName, req.TensorboardLogDir, req.Description)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary get visual job List
// @Produce json
// @Param param body GetVisualJobListReq true "params"
// @Success 200 {object} GetVisualJobListRsq "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual/list [get]
func getVisualJobList(c *gin.Context) error {
	var req GetVisualJobListReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)

	visualJobList, totalJobsNum, totalPagesNum, err := services.GetAllVisualJobInfo(userName, req.PageNum, req.PageSize, req.OrderBy, req.Status, req.JobName, req.Order)
	if err != nil {
		return ParameterError(err.Error())
	}
	var visualJobListRspUnitArray []VisualJobListRspUnit
	for _, visualJob := range visualJobList {
		newVisualJobListRspUnit := VisualJobListRspUnit{
			Id:                visualJob.ID,
			JobName:           visualJob.Name,
			Status:            visualJob.Status,
			CreateTime:        visualJob.CreatedAt,
			TensorboardLogDir: visualJob.LogPath,
			Description:       visualJob.Description,
		}
		visualJobListRspUnitArray = append(visualJobListRspUnitArray, newVisualJobListRspUnit)
	}
	fmt.Printf("%d",totalJobsNum)
	rsp := GetVisualJobListRsq{
		Templates:    visualJobListRspUnitArray,
		TotalJobsNum: totalJobsNum,
		TotalPages:   totalPagesNum,
	}
	return SuccessResp(c, rsp)
}

// @Summary get visual job endpoints address
// @Produce json
// @Param param body GetEndpointsReq true "params"
// @Success 200 {object} GetEndpointsRsq "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual/endpoints [get]
func getEndpoints(c *gin.Context) error {
	var req GetEndpointsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)
	path, err := services.GetEndpointsPath(userName, req.JobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	rsp := GetEndpointsRsq{
		Path: path,
	}
	return SuccessResp(c, rsp)
}

// @Summary delete visual job
// @Produce json
// @Param param body DeleteJobReq true "params"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual/ [delete]
func deleteVisualJob(c *gin.Context) error {
	var req DeleteJobReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)
	err = services.DeleteVisualJob(userName, req.JobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary switch job status
// @Produce json
// @Param param body SwitchVisualJobStatusReq true "params"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual/ [put]
func switchVisualJobStatus(c *gin.Context) error {
	var req SwitchVisualJobStatusReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)
	if req.Status == "pause" {
		err = services.StopVisualJob(userName, req.JobId)
		if err != nil {
			return AppError(APP_ERROR_CODE, err.Error())
		}
	}
	if req.Status == "running" {
		err = services.ContinueVisualJob(userName, req.JobId)
		if err != nil {
			return AppError(APP_ERROR_CODE, err.Error())
		}
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
