package routers

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
	"os"
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

type createModelsetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" `
	JobId       string `json:"jobId"`
	CodePath    string `json:"codePath"`
	ParamPath   string `json:"paramPath"`
	IsAdvance   bool   `json:"isAdvance,default=false"`

	Use         string `json:"use"`
	Size        int64  `json:"size"`
	DataFormat  string `json:"dataFormat"`
	DatasetName string `json:"datasetName"`
	DatasetPath string `json:"datasetPath"`
	//omitempty 值为空，不编码
	Params    map[string]string `json:"params"`
	Engine    string            `json:"engine"`
	Precision string            `json:"precision"`
	//指定的模型参数路径
	// 输出文件路径
	OutputPath string `json:"outputPath"`
	//启动文件路径
	StartupFile string `json:"startupFile"`

	//用于可视化建模平台直接启动训练任务
	JobTrainingType string          `json:"jobTrainingType"`
	NumPs           int             `json:"numPs"`
	NumPsWorker     int             `json:"numPsWorker"`
	DeviceType      string          `json:"deviceType"`
	DeviceNum       int             `json:"deviceNum"`
	Nodes           []AvisualisNode `json:"nodes"`
	Arguments       []AvisualisNode `json:"arguments"`
	Edges           []AvisualisEdge `json:"edges"`
}
type AvisualisEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
type AvisualisNode struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	IDX    int         `json:"idx"`
	Config interface{} `json:"config"`
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
// @Param body body createModelsetReq true "jsonbody"
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
	//if req.CodePath != "" {
	//	err = services.CheckPathExists(req.CodePath)
	//	if err != nil {
	//		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//	}
	//}
	//if req.ParamPath != "" {
	//	//检查模型参数文件是否存在
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
		if req.JobTrainingType != models.TrainingTypeDist && req.JobTrainingType != models.TrainingTypeRegular {
			return AppError(INVALID_TRAINING_TYPE, "任务类型非法")
		}
		//存储节点json
		nodesBytes, _ := json.Marshal(req.Nodes)
		edgesBytes, _ := json.Marshal(req.Edges)

		//去掉nodes没用的节点并存入json
		pipelineConfigPath, err := services.GetModelTempPath(services.FILETYPE_JSON)
		f, err := os.OpenFile(pipelineConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)
		if err != nil {
			return AppError(FILEPATH_NOT_VALID_CODE, "cannot create new config file")
		}
		_, err = f.Write(nodesBytes)
		defer f.Close()

		//把数据传入params后端算法只需要pipeline_config
		req.Params = make(map[string]string)
		req.Params["pipeline_config"] = pipelineConfigPath

		req.Params["pipeline_config"] = "/data/premodel/code/ApulisVision/panel.json"

		//baseconfig待定，
		//req.Params["config"] = req.CodePath

		training := models.Training{
			Id:              req.JobId,
			Name:            req.Name,
			Engine:          req.Engine,
			CodePath:        req.CodePath,
			StartupFile:     req.StartupFile,
			OutputPath:      req.OutputPath,
			DatasetPath:     req.DatasetPath,
			Params:          req.Params,
			Desc:            req.Description,
			NumPs:           req.NumPs,
			NumPsWorker:     req.NumPsWorker,
			DeviceType:      req.DeviceType,
			DeviceNum:       req.DeviceNum,
			JobTrainingType: req.JobTrainingType,
		}
		//启动训练作业
		jobId, err := services.CreateTraining(username, training)
		if err != nil {
			return AppError(FAILED_START_TRAINING, err.Error())
		}
		req.JobId = jobId
		//nodes和edges只用存储然后传给前端
		req.Params["nodes"] = string(nodesBytes)
		req.Params["edges"] = string(edgesBytes)
	}
	err = services.CreateModelset(req.Name, req.Description, username, "0.0.1", req.JobId, req.CodePath, req.ParamPath, req.IsAdvance,
		req.Use, req.Size, req.DataFormat, req.DatasetName, req.DatasetPath, req.Params, req.Engine, req.Precision, req.OutputPath, req.StartupFile)
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
