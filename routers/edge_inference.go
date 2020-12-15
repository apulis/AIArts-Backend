package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupEdgeInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/edge_inferences")

	group.Use(Auth())

	group.GET("/", wrapper(lsEdgeInferences))
	group.POST("/", wrapper(createEdgeInference))
	group.GET("/conversion_types", wrapper(getConversionTypes))
	group.GET("/fdinfo", wrapper(getFDInfo))
	group.POST("/fdinfo", wrapper(setFDInfo))
	group.POST("/push/:jobId", wrapper(pushToFD))
	group.DELETE("/:jobId", wrapper(deleteEdgeInference))
}

type edgeInferenceId struct {
	ID string `uri:"jobId" binding:"required"`
}

type setFDInfoReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Url      string `json:"url" binding:"url"`
}

type pushToFDReq struct {
	JobId string `uri:"jobId" binding:"required"`
}

type GetFDInfoResp struct {
	FDInfo models.FDInfo `json:"fdinfo"`
}

type LsEdgeInferencesResp struct {
	EdgeInferences []models.ConversionJob `json:"edgeInferences"`
	Total          int                    `json:"total"`
	TotalPage      int                    `json:"totalPage"`
	PageNum        int                    `json:"pageNum"`
	PageSize       int                    `json:"pageSize"`
}

type CreateEdgeInferenceResp struct {
	JobId string `json:"jobId"`
}

// @Summary get dataset by id
// @Produce  json
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Param jobName query string false "job name"
// @Param modelconversionType query string false "model conversion type"
// @Param orderBy query string false "order by item"
// @Param order query string false "desc or asc"
// @Success 200 {object} APISuccessRespLsEdgeInferences "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences [get]
func lsEdgeInferences(c *gin.Context) error {
	var req models.LsEdgeInferencesReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	if req.VCName == "" {
		req.VCName = "platform"
	}

	conversionList, total, err := services.LsEdgeInferences(username, req)
	if err != nil {
		return RemoteServerError(err.Error())
	}

	res := LsEdgeInferencesResp{
		EdgeInferences: conversionList,
		Total:          total,
		PageNum:        req.PageNum,
		PageSize:       req.PageSize,
		TotalPage:      total/req.PageSize + 1,
	}

	return SuccessResp(c, res)
}

// @Summary update dataset
// @Produce  json
// @Param body body string true "json body"
// @Success 200 {object} APISuccessRespCreateEdgeInference "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences [post]
func createEdgeInference(c *gin.Context) error {

	var req models.CreateEdgeInferenceReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	jobId, err := services.CreateEdgeInference(username, req)
	if err != nil {
		return RemoteServerError(err.Error())
	}

	return SuccessResp(c, CreateEdgeInferenceResp{JobId: jobId})
}

// @Summary get dataset by id
// @Produce  json
// @Success 200 {object} APISuccessRespGetConvTypes "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/conversion_types [get]
func getConversionTypes(c *gin.Context) error {
	convTypes, err := services.GetConversionTypes()
	if err != nil {
		return RemoteServerError(err.Error())
	}
	return SuccessResp(c, convTypes)
}

// @Summary get dataset by id
// @Produce  json
// @Success 200 {object} APISuccessRespGetFDInfo "success"
// @Failure 500 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/fdinfo [get]
func getFDInfo(c *gin.Context) error {
	fdinfo, err := services.GetFDInfo()
	if err != nil {
		return RemoteServerError(err.Error())
	}
	data := GetFDInfoResp{FDInfo: fdinfo}
	return SuccessResp(c, data)
}

// @Summary update dataset
// @Produce  json
// @Param body body setFDInfoReq true "json body"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/fdinfo [post]
func setFDInfo(c *gin.Context) error {
	var req setFDInfoReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	res, err := services.SetFDInfo(req.Username, req.Password, req.Url)
	if err != nil {
		return RemoteServerError(err.Error())
	}
	if res {
		return SuccessResp(c, gin.H{})
	} else {
		return ServeError(configs.FDINFO_SET_ERROR, "fd set failed")
	}
}

// @Summary update dataset
// @Produce  json
// @Param id path string true "job id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/push/:jobId [post]
func pushToFD(c *gin.Context) error {
	var req pushToFDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.PushToFD(req.JobId)
	if err != nil {
		return ServeError(configs.FD_PUSH_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary delete edge_inference by jobId
// @Produce  json
// @Param jobId path string true "job id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/:jobId [delete]
func deleteEdgeInference(c *gin.Context) error {
	var jobId edgeInferenceId
	err := c.ShouldBindUri(&jobId)
	if err != nil {
		return ParameterError(err.Error())
	}
	resp, err := services.DeleteJob(jobId.ID)
	if err != nil {
		return RemoteServerError(err.Error())
	}
	return SuccessResp(c, resp)
}
