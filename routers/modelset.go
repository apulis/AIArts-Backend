package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
	"strings"
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

type getModelsetResp struct {
	Model    models.Modelset  `json:"model"`
	Training *models.Training `json:"training"`
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
// @Param query query models.LsModelsetsReq true "query"
// @Success 200 {object} getModelsetResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [get]
func lsModelsets(c *gin.Context) error {

	var req models.LsModelsetsReq
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

	modelsets, total, err = services.ListModelSets(username, req)

	if strings.HasPrefix(req.Use, `Avisualis`) {
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
	var training *models.Training

	if err != nil {
		return ParameterError(err.Error())
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	modelset, err := services.GetModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	//插入数据集面板
	if strings.HasPrefix(modelset.Use, `Avisualis`) {
		modelset, err = services.GeneratePanel(modelset, username)
		if err != nil {
			return err
		}
		//如果有Jobid
		if modelset.JobId != "" {
			//获取新的training
			training, _ = services.GetTraining(username, modelset.JobId)
			data := getModelsetResp{Model: modelset, Training: training}
			return SuccessResp(c, data)
		}
	}
	data := getModelsetResp{Model: modelset, Training: training}
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

	var req models.CreateModelsetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	//if len(req.VCName) == 0 {
	//	req.VCName = models.DefaultVcName
	//}

	//如果上传模型文件检查路径是否存在
	if req.CodePath != "" {
		err = services.CheckPathExists(req.CodePath)
		if err != nil {
			return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
		}
	}

	//avisualis不检测参数路径是否存在
	if req.ParamPath != "" && !strings.HasPrefix(req.Use, `Avisualis`) {
		err = services.CheckPathExists(req.ParamPath)
		if err != nil {
			return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
		}
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	if strings.HasPrefix(req.Use, `Avisualis`) && !req.IsAdvance {
		req, err = services.CreateAvisualisTraining(c, req, username)
		if err != nil {
			return err
		}
	}

	err = services.CreateModelset(username, "0.0.1", req)
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
	var req models.CreateModelsetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)

	//更新Avisualis任务
	if strings.HasPrefix(req.Use, `Avisualis`) {
		if req.JobTrainingType != models.TrainingTypeDist && req.JobTrainingType != models.TrainingTypeRegular {
			return AppError(INVALID_TRAINING_TYPE, "任务类型非法")
		}
		if req.JobId != "" {
			//重启新的training
			_ = services.DeleteTraining(username, req.JobId)
		}
		//更新节点的parma配置节点
		req, err = services.CreateAvisualisTraining(c, req, username)
		if err != nil {
			return err
		}
	}

	err = services.UpdateModelset(id.ID, "0.0.1", req)
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
	username := getUsername(c)
	//删除可视化建模同时删除停止任务
	model, _ :=services.GetModelset(id.ID)
	if strings.HasPrefix(model.Use, `Avisualis`) {
		_ = services.DeleteTraining(username, model.JobId)
	}
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
