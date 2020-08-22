package routers

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupVisualJob(r *gin.Engine) {
	group := r.Group("/ai_arts/api/visual")

	group.Use(Auth())

	group.POST("/", wrapper(createVisualJob))
	group.GET("/list", wrapper(getVisualJobList))
	group.GET("/enpoints", wrapper(getEndpoints))
	group.DELETE("/", wrapper(deleteVisualJob))
	group.PUT("/:id", wrapper(switchVisualJobStatus)) //
}

type CreateVisualJobReq struct {
	JobName           string `json:jobname`
	TensorboardLogDir string `json:tensorboardLogDir`
	Description       string `json:description`
}

type GetVisualJobListReq struct {
	PageNum  int    `json:pageNum`
	PageSize int    `json:pageSize`
	OrderBy  string `json:orderBy`
	Status   string `json:status`
	JobName  string `json:JobName`
	Order    string `json:order`
}

type GetVisualJobListRsq struct {
	Templates    []VisualJobListRspUnit `json:Templates`
	TotalJobsNum int                    `json:total`
	TotalPages   int                    `json:totalPages`
}

type VisualJobListRspUnit struct {
	Id                int    `json:id`
	JobName           string `json:jobName`
	Status            string `json:status`
	CreateTime        string `json:createTime`
	TensorboardLogDir string `json:TensorboardLogDir`
	Description       string `json:description`
}

type GetRndpointsReq struct {
	JobId int `json:id`
}
type GetRndpointsRsq struct {
	Path string `json:path`
}

// @Summary create visual job
// @Produce json
// @Param param body CreateVisualJobReq true "params"
// @Success 200 {object} APISuccessRespCreateTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/visual [post]
func createVisualJob(c *gin.Context) error {
	var req CreateVisualJobReq
	err := c.ShouldBindJSON(&req)
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

func getVisualJobList(c *gin.Context) error {
	var req GetVisualJobListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)

	visualJobList, totalJobsNum, totalPagesNum, err := services.GetAllVisualJobInfo(userName, req.PageNum, req.PageSize, req.OrderBy, req.Status, req.JobName, req.Order)
	if err != nil {
		return ParameterError(err.Error())
	}
	visualJobListRspUnitArray := make([]VisualJobListRspUnit, len(visualJobList))
	for _, visualJob := range visualJobList {
		newVisualJobListRspUnit := VisualJobListRspUnit{
			Id:                visualJob.Id,
			JobName:           visualJob.Name,
			Status:            visualJob.Status,
			CreateTime:        visualJob.CreateTime,
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

func getEndpoints(c *gin.Context) error {
	var req GetRndpointsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	userName := getUsername(c)
	path, err := services.GetEndpointsPath(userName, req.JobId)
	if err != nil {
		return ParameterError(err.Error())
	}
	rsp := GetRndpointsRsq{
		Path: path,
	}
	return SuccessResp(c, rsp)
}

func deleteVisualJob(c *gin.Context) error {
	return nil
}

func switchVisualJobStatus(c *gin.Context) error {
	return nil
}
