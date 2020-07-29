package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupModel(r *gin.Engine) {
	group := r.Group("/ai_arts/api")
	group.Use(Auth())

	group.GET("/models/", wrapper(lsModelsets))
	group.GET("/models/:id", wrapper(getModelset))
	group.POST("/models/", wrapper(createModelset))
	group.POST("/models/:id", wrapper(updateModelset))
	group.DELETE("/models/:id", wrapper(deleteModelset))
	group.GET("/models/:id/evaluation", wrapper(getEvaluation))
	group.POST("/models/:id/evaluation", wrapper(createEvaluation))

}

type modelsetId struct {
	ID int `uri:"id" binding:"required"`
}
type createEvaluationResp struct {
	EvaluationId string `json:"jobId"`
}
type getEvaluationResp struct {
	ModelName    string `json:"modelName"`
	EngineType   string `json:"engineType"`
	DeviceType   string `json:"deviceType"`
	DeviceNum    int    `json:"deviceNum"`
	StartupFile  string `json:"startupFile"`
	OutputPath   string `json:"outputPath"`
	CreatedAt    string `json:"createdAt"`
	Status       string `json:"status"`
	DatasetName  string `json:"datasetName"`
	ArgumentPath string `json:"argumentPath"`
	Log          string `json:"log"`
}

type lsModelsetsReq struct {
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize,default=10"`
	Name      string `form:"name"`
	Status    string `form:"status"`
	IsAdvance bool   `form:"isAdvance"`
	OrderBy   string `form:"orderBy,default=created_at"`
	Order     string `form:"order,default=desc"`
}

type createModelsetReq struct {
	Name         string            `json:"name" binding:"required"`
	Description  string            `json:"description" `
	JobId        string            `json:"jobId"`
	Use          string            `json:"use"`
	DataFormat   string            `json:"dataFormat"`
	Arguments    map[string]string `json:"arguments"`
	EngineType   string            `json:"engineType"`
	Precision    string            `json:"precision"`
	IsAdvance    bool              `json:"isAdvance"`
	ModelPath    string            `json:"modelPath"`
	ArgumentPath string            `json:"argumentPath" binding:"required"`
}

type updateModelsetReq struct {
	Description string `json:"description" binding:"required"`
}

type GetModelsetResp struct {
	Model models.Modelset `json:"model"`
}

type GetModelsetsResp struct {
	Models    []models.Modelset `json:"models"`
	Total     int               `json:"total"`
	TotalPage int               `json:"totalPage"`
	PageNum   int               `json:"pageNum"`
	PageSize  int               `json:"pageSize"`
}

// @Summary list models
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Param isAdvance query string true "job status. get all jobs if it is all"
// @Param name query string true "the keyword of search"
// @Param status query string true "the keyword of search"
// @Success 200 {object} APISuccessRespGetModelsets "success"
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
// @Success 200 {object} APISuccessRespGetModelset "success"
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
	data := GetModelsetResp{Model: modelset}
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
	if req.ModelPath != "" {
		err = services.CheckPathExists(req.ModelPath)
		if err != nil {
			return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
		}
	}
	//检查模型参数文件是否存在
	err = services.CheckPathExists(req.ArgumentPath)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.CreateModelset(req.IsAdvance, req.Name, req.Description, username, "0.0.1", req.Use, req.JobId,
		req.DataFormat, req.Arguments, req.EngineType, req.Precision, req.ModelPath, req.ArgumentPath)
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

// @Summary create Training
// @Produce json
// @Param param body services.CreateEvaluationReq true "params"
// @Success 200 {object} APISuccessRespCreateTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id/evaluation [post]
func createEvaluation(c *gin.Context) error {
	var req services.CreateEvaluationReq
	var id int
	err := c.ShouldBindUri(&id)
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelset, err := services.GetModelset(id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	username := getUsername(c)

	//检查模型文件是否存在
	//err = services.CheckPathExists(req.DatasetPath)
	//if err != nil {
	//	return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//}
	////检查模型参数文件是否存在
	//err = services.CheckPathExists(req.ArgumentPath)
	//if err != nil {
	//	return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//}
	////检查输出路径是否存在
	//err = services.CheckPathExists(req.OutputPath)
	//if err != nil {
	//	return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//}

	jobId, err := services.CreateEvaluation(username, req)
	if err != nil {
		return AppError(CREATE_TRAINING_FAILED_CODE, err.Error())
	}
	//更新评估参数
	modelset.DatasetName = req.DatasetName
	modelset.EngineType = req.EngineType
	modelset.DatasetPath = req.DatasetPath
	modelset.OutputPath = req.OutputPath
	modelset.StartupFile = req.StartupFile
	modelset.EvaluationId = jobId

	err = models.UpdateModelset(&modelset)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := createEvaluationResp{
		EvaluationId: jobId,
	}
	return SuccessResp(c, data)
}

// @Summary get evaluation by modelid
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} APISuccessRespGetModelset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id/evaluation [get]
func getEvaluation(c *gin.Context) error {
	var id int
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelset, err := services.GetModelset(id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	username := getUsername(c)
	job, err := services.GetTraining(username, modelset.EvaluationId)
	if err != nil {
		return AppError(CREATE_TRAINING_FAILED_CODE, err.Error())
	}
	log, err := services.GetTrainingLog(username, modelset.EvaluationId)
	logResp := ""
	if log != nil {
		logResp = log.Log
	}

	data := getEvaluationResp{
		ModelName:   modelset.Name,
		EngineType:  modelset.EngineType,
		DeviceType:  job.DeviceType,
		DeviceNum:   job.DeviceNum,
		CreatedAt:   job.CreateTime,
		Status:      job.Status,
		DatasetName: modelset.DatasetName,
		Log:         logResp,
	}

	return SuccessResp(c, data)
}
