package routers

import "github.com/gin-gonic/gin"

func AddGroupEdgeInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/edge_inferences")

	group.Use(Auth())

	group.GET("/", wrapper(lsEdgeInferences))
	group.POST("/", wrapper(createEdgeInference))
	group.GET("/conversion_types", wrapper(getConversionTypes))
	group.GET("/fdinfo", wrapper(getFDInfo))
	group.POST("/fdinfo", wrapper(setFDInfo))
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

// @Summary get dataset by id
// @Produce  json
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences [get]
func lsEdgeInferences(c *gin.Context) error {
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
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/conversion_types [get]
func getConversionTypes(c *gin.Context) error {
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary get dataset by id
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/edge_inferences/fdinfo [get]
func getFDInfo(c *gin.Context) error {
	data := gin.H{}
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
	data := gin.H{}
	return SuccessResp(c, data)
}
