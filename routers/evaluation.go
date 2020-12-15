package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
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
type getLogReq struct {
	PageNum int `form:"pageNum" json:"pageNum"`
}

type getEvaluationsResp struct {
	Evaluations []*services.Evaluation `json:"evaluations"`
	Total       int                    `json:"total"`
	TotalPage   int                    `json:"totalPage"`
	PageNum     int                    `json:"pageNum"`
	PageSize    int                    `json:"pageSize"`
}

type getEvaluationResp struct {
	Evaluation *services.Evaluation `json:"evaluation"`
	Log        *models.JobLog       `json:"log"`
	Indicator  map[string]string    `json:"indicator"`
	Confusion  map[string]string    `json:"confusion"`
}

// @Summary list evaluations
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

	var req models.GetEvaluationsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	if req.VCName == "" {
		req.VCName = models.DefaultVcName
	}

	evaluations, total, totalPage, err := services.GetEvaluations(username, req)
	if err != nil {
		return AppError(configs.CREATE_EVALUATION_FAILED_CODE, err.Error())
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
// @Param param body services.Evaluation true "name:model name ， datasetName ：dataset name"
// @Success 200 {object} createEvaluationResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations [post]
func createEvaluation(c *gin.Context) error {
	var req services.Evaluation

	err := c.BindJSON(&req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	if req.VCName == "" {
		req.VCName = models.DefaultVcName
	}

	jobId, err := services.CreateEvaluation(c, username, req)
	if err != nil {
		return AppError(configs.CREATE_EVALUATION_FAILED_CODE, err.Error())
	}

	//更新评估参数
	data := createEvaluationResp{
		EvaluationId: jobId,
	}
	return SuccessResp(c, data)
}

// @Summary get evaluation by id
// @Produce  json
// @Param id path int true "evaluation id"
// @Success 200 {object} getEvaluationResp "success indicator:{"accuary":"0.001"},confusion"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations/:id [get]
func getEvaluation(c *gin.Context) error {
	var id getEvaluationReq
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var logReq getLogReq
	err = c.ShouldBindQuery(&logReq)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}
	job, err := services.GetEvaluation(username, id.ID)
	if err != nil {
		return AppError(configs.CREATE_EVALUATION_FAILED_CODE, err.Error())
	}
	log, err := services.GetEvaluationLog(username, id.ID, logReq.PageNum)
	// 请求最后一页日志以获取评估指标
	var maxPageLog *models.JobLog
	var indicator map[string]string
	var confusion map[string]string
	if log != nil {
		maxPageLog, err = services.GetEvaluationLog(username, id.ID, log.MaxPage)
		if maxPageLog != nil {
			indicator, confusion = services.GetRegexpLog(maxPageLog.Log)
		}
	}
	data := getEvaluationResp{
		Evaluation: job,
		Log:        log,
		Indicator:  indicator,
		Confusion:  confusion,
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
	var id getEvaluationReq
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}
	err = services.DeleteEvaluation(username, id.ID)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
