package routers

import (
	"fmt"

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
	group.POST("/push/:id", wrapper(pushToFD))
}

type createEdgeInferenceReq struct {
	JobName        string `json:"jobName" binding:"required"`
	InputPath      string `json:"inputPath" binding:"required"`
	OutputPath     string `json:"outputPath" binding:"required"`
	ConvertionType string `json:"convertionType" binding:"required"`
	Device         string `json:"device" binding:"required"`
}

type setFDInfoReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Url      string `json:"url" binding:"url"`
}

type lsEdgeInferencesReq struct {
	PageNum  int `form:"pageNum"`
	PageSize int `form:"pageSize,default=10"`
}

type GetFDInfoResp struct {
	FDInfo models.FDInfo `json:"fdinfo"`
}

// @Summary get dataset by id
// @Produce  json
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences [get]
func lsEdgeInferences(c *gin.Context) error {
	var req lsEdgeInferencesReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.LsEdgeInferences(req.PageNum, req.PageSize, username)
	fmt.Println(err)
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary update dataset
// @Produce  json
// @Param body body createEdgeInferenceReq true "json body"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences [post]
func createEdgeInference(c *gin.Context) error {
	data := gin.H{}
	return SuccessResp(c, data)
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
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
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
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	if fdinfo == (models.FDInfo{}) {
		return ServeError(FDINFO_NOT_SET, "fdinfo not set")
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
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	if res {
		return SuccessResp(c, gin.H{})
	} else {
		return ServeError(FDINFO_SET_ERROR, "fd set failed")
	}
}

// @Summary update dataset
// @Produce  json
// @Param id path string true "job id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/push/:id [post]
func pushToFD(c *gin.Context) error {
	data := gin.H{}
	return SuccessResp(c, data)
}
