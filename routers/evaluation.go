package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupEvaluation(r *gin.Engine) {
	group := r.Group("/ai_arts/api/evaluations/")
	group.Use(Auth())
	group.GET("/", wrapper(lsEvaluations))
	group.GET("/:id", wrapper(getEvaluation))
	group.POST("/", wrapper(createEvaluation))
	//group.POST("/:id", wrapper(updateModelset))
	group.DELETE("/:id", wrapper(stopEvaluation))
}

type getEvaluationReq struct {
	ID string `uri:"id" binding:"required"`
}
type getEvaluationsReq struct {
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
	Status   string `form:"status" json:"status"`
	Name     string `form:"name" json:"name"`
	OrderBy  string `form:"orderBy" json:"orderBy"`
	Order    string `form:"order" json:"order"`
}
type getEvaluationsResp struct {
	Evaluations []*models.Training `json:"evaluations"`
	Total       int                `json:"total"`
	TotalPage   int                `json:"totalPage"`
	PageNum     int                `json:"pageNum"`
	PageSize    int                `json:"pageSize"`
}
type getEvaluationResp struct {
	Evaluation models.Training   `json:"evaluation"`
	Log        string            `json:"log"`
	Indicator  map[string]string `json:"indicator"`
}

// @Summary list models
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Param name query string true "the keyword of search"
// @Param status query string true "the keyword of search"
// @Success 200 {object} APISuccessRespGetModelsets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations [get]
func lsEvaluations(c *gin.Context) error {
	var req getEvaluationsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	if req.Status == "" {
		req.Status = "all"
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	evaluations, total, totalPage, err := services.GetEvaluations(username, req.PageNum, req.PageSize,
		req.Status, req.Name, req.OrderBy, req.Order)
	if err != nil {
		return AppError(CREATE_EVALUATION_FAILED_CODE, err.Error())
	}
	data := getEvaluationsResp{
		Evaluations: evaluations,
		Total:       total,
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
		TotalPage:   totalPage,
	}
	return SuccessResp(c, data)
}

// @Summary create Evaluation
// @Produce json
// @Param param body models.Training true "ID:modelID ， NAME : model NAME Desc：dataset Name"
// @Success 200 {object} createEvaluationResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations [post]
func createEvaluation(c *gin.Context) error {
	var req models.Training
	var id int
	err := c.ShouldBindUri(&id)
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
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
		return AppError(CREATE_EVALUATION_FAILED_CODE, err.Error())
	}
	//更新评估参数
	//var argItem models.ArgumentsItem
	//argItem = req.Arguments
	//modelset.DatasetName = req.DatasetName
	//modelset.EngineType = req.EngineType
	//modelset.DatasetPath = req.DatasetPath
	//modelset.OutputPath = req.OutputPath
	//modelset.StartupFile = req.StartupFile
	//modelset.Arguments = &argItem
	//modelset.EvaluationId = jobId
	//err = models.UpdateModelset(&modelset)
	//if err != nil {
	//	return AppError(APP_ERROR_CODE, err.Error())
	//}
	data := createEvaluationResp{
		EvaluationId: jobId,
	}
	return SuccessResp(c, data)
}

// @Summary get evaluation by id
// @Produce  json
// @Param id path int true "evaluation id"
// @Success 200 {object} getEvaluationResp "success {"accuary":"0.001"}"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations/:id [get]
func getEvaluation(c *gin.Context) error {
	var id getEvaluationReq
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	//modelset, err := services.GetModelset(id)
	//if err != nil {
	//	return AppError(APP_ERROR_CODE, err.Error())
	//}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	job, err := services.GetEvaluation(username, id.ID)
	if err != nil {
		return AppError(CREATE_EVALUATION_FAILED_CODE, err.Error())
	}
	log, err := services.GetEvaluationLog(username, id.ID)
	logResp := ""
	if log != nil {
		logResp = log.Log
	}
	indicator := services.GetRegexpLog(logResp)
	data := getEvaluationResp{
		Evaluation: *job,
		Log:        logResp,
		Indicator:  indicator,
	}
	return SuccessResp(c, data)
}

// @Summary delete evaluation by id
// @Produce  json
// @Param id path int true "evaluation id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations/:id [delete]
func stopEvaluation(c *gin.Context) error {
	var id string
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.DeleteEvaluation(username, id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
