package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupModel(r *gin.Engine) {
	group := r.Group("/ai_arts/api/models/")
	group.Use(Auth())
	group.GET("/", wrapper(lsModelsets))
	group.GET("/:id", wrapper(getModelset))
	group.POST("/", wrapper(createModelset))
	group.POST("/:id", wrapper(updateModelset))
	group.DELETE("/:id", wrapper(deleteModelset))

}

type modelsetId struct {
	ID int `uri:"id" binding:"required"`
}
type createEvaluationResp struct {
	EvaluationId string `json:"jobId"`
}

type lsModelsetsReq struct {
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize,default=10"`
	Name      string `form:"name"`
	//all
	Status    string `form:"status"`
	IsAdvance bool   `form:"isAdvance"`
	OrderBy   string `form:"orderBy,default=created_at"`
	Order     string `form:"order,default=desc"`
}

type createModelsetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" `
	JobId       string `json:"jobId"`
	CodePath   string `json:"codePath"`
	ParamPath   string `json:"paramPath" binding:"required"`
	IsAdvance bool   `json:"isAdvance,default=false"`

	Use   string `json:"use"`
	Size   int64 `json:"size"`
	DataFormat  string `json:"dataFormat"`
	DatasetName string `json:"datasetName"`
	DatasetPath string `json:"datasetPath"`
	//omitempty 值为空，不编码
	Params  map[string]string `json:"params"`
	Engine string         `json:"engine"`
	Precision  string         `json:"precision"`
	//指定的模型参数路径
	// 输出文件路径
	OutputPath string `json:"outputPath"`
	//启动文件路径
	StartupFile string `json:"startupFile"`
}

type updateModelsetReq struct {
	Description string `json:"description" binding:"required"`
}

type getModelsetResp struct {
	Model models.Modelset `json:"model"`
}

type GetModelsetsResp struct {
	Models    []models.Modelset `json:"models"`
	Total     int               `json:"total"`
	TotalPage int               `json:"totalPage"`
	PageNum   int               `json:"pageNum"`
	PageSize  int               `json:"pageSize"`
}

// @Summary get model by id
// @Produce  json
// @Param query query lsModelsetsReq true "query"
// @Success 200 {object} getModelsetResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [get]
func lsModelsets(c *gin.Context) error {
	var req lsModelsetsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	var modelsets []models.Modelset
	var total int
	//获取当前用户创建的模型
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	modelsets, total, err = services.ListModelSets(req.PageNum, req.PageSize, req.OrderBy, req.Order, req.IsAdvance, req.Name, req.Status, username)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := GetModelsetsResp{
		Models:    modelsets,
		Total:     total,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		TotalPage: total/req.PageSize + 1,
	}
	return SuccessResp(c, data)
}

// @Summary get model by id
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} getModelsetResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [get]
func getModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelset, err := services.GetModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := getModelsetResp{Model: modelset}
	return SuccessResp(c, data)
}

// @Summary create model
// @Produce  json
// @Param body body createModelsetReq true "json body"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [post]
func createModelset(c *gin.Context) error {
	var req createModelsetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	//如果上传模型文件检查模型文件是否存在
	if req.CodePath != "" {
		err = services.CheckPathExists(req.CodePath)
		if err != nil {
			return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
		}
	}
	//检查模型参数文件是否存在
	err = services.CheckPathExists(req.ParamPath)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err = services.CreateModelset(req.Name, req.Description, username, "0.0.1", req.JobId, req.CodePath, req.ParamPath,req.IsAdvance,
		req.Use,req.Size,req.DataFormat,req.DatasetName,req.DatasetPath,req.Params,req.Engine,req.Precision,req.OutputPath,req.StartupFile)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary update model
// @Produce  json
// @Param description path string true "model description"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [post]
func updateModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req updateModelsetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateModelset(id.ID, req.Description)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary delete model by id
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [delete]
func deleteModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.DeleteModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
