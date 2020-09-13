package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupModel(r *gin.Engine) {
	group := r.Group("/ai_arts/api/models/")
	group.Use(Auth())
	group.GET("/:id/panel", wrapper(getPanel))
	group.GET("/", wrapper(lsModelsets))
	group.GET("/:id", wrapper(getModelset))
	group.POST("/", wrapper(createModelset))
	group.POST("/:id", wrapper(updateModelset))
	group.DELETE("/:id", wrapper(deleteModelset))

}

type modelsetId struct {
	ID int `uri:"id" binding:"required"`
}
type modelsetUse struct {
	USE string `uri:"id" binding:"required"`
}

type createEvaluationResp struct {
	EvaluationId string `json:"jobId"`
}

type lsModelsetsReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize,default=10"`
	Name     string `form:"name"`
	//all
	Use       string `form:"use"`
	Status    string `form:"status"`
	IsAdvance bool   `form:"isAdvance"`
	OrderBy   string `form:"orderBy,default=created_at"`
	Order     string `form:"order,default=desc"`
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

	modelsets, total, err = services.ListModelSets(req.PageNum, req.PageSize, req.OrderBy, req.Order, req.IsAdvance, req.Name, req.Status, req.Use, username)
	if req.Use != "" {
		for i := 0; i < len(modelsets); i++ {
			training, err := services.GetTraining(username, modelsets[i].JobId)
			if err != nil {
				return AppError(FAILED_START_TRAINING, err.Error())
			}
			modelsets[i].Status = training.Status
		}
	}
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

// @Summary get visualis panel
// @Produce  json
// @Param query query modelsetUse true "usetype"
// @Success 200 {object} getModelsetResp "Avisualis_Classfication-Avisualis_ObjectDetection-Avisualis_SemanticSegmentation"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id/panel [get]
func getPanel(c *gin.Context) error {
	var use modelsetUse
	err := c.ShouldBindUri(&use)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	panel, err := services.GetPanel(use.USE, username)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, panel)
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
// @Param body body services.CreateModelsetReq true "jsonbody"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [post]
func createModelset(c *gin.Context) error {
	var req services.CreateModelsetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	////如果上传模型文件检查路径是否存在
	//if req.CodePath != "" {
	//	err = services.CheckPathExists(req.CodePath)
	//	if err != nil {
	//		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//	}
	//}
	////检查模型参数文件是否存在
	//if req.ParamPath != "" {
	//	err = services.CheckPathExists(req.ParamPath)
	//	if err != nil {
	//		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//	}
	//}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	if req.Use != "" {
		req, err = services.CreateAvisualisTraining(req, username)
		if err != nil {
			return err
		}
	}
	err = services.CreateModelset(req.Name, req.Description, username, "0.0.1", req.JobId, req.CodePath, req.ParamPath, req.IsAdvance,
		req.Use, req.Size, req.DataFormat, req.DatasetName, req.DatasetPath, req.Params, req.Engine, req.Precision, req.OutputPath, req.StartupFile,req.DeviceType,req.DeviceNum)
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
	var req services.CreateModelsetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if req.JobId != "" {
		req.JobTrainingType = models.TrainingTypeRegular
		job, _ := services.GetTraining(username, req.JobId)
		req.JobTrainingType = job.JobTrainingType
		_ = services.DeleteTraining(username, req.JobId)
	}
	if req.JobTrainingType != models.TrainingTypeDist && req.JobTrainingType != models.TrainingTypeRegular {
		return AppError(INVALID_TRAINING_TYPE, "任务类型非法")
	}
	if req.Use != "" {
		req, err = services.CreateAvisualisTraining(req, username)
		if err != nil {
			return err
		}
	}
	err = services.UpdateModelset(id.ID, req.Name, req.Description, "0.0.1", req.JobId, req.CodePath, req.ParamPath,
		req.Use, req.Size, req.DataFormat, req.DatasetName, req.DatasetPath, req.Params, req.Engine, req.Precision, req.OutputPath, req.StartupFile)
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
