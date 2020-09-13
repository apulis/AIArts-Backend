package routers

import (
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
	Search   string `form:"search" json:"search"`
	OrderBy  string `form:"orderBy" json:"orderBy"`
	Order    string `form:"order" json:"order"`
}

type getEvaluationsResp struct {
	Evaluations []*services.Evaluation `json:"evaluations"`
	Total       int                    `json:"total"`
	TotalPage   int                    `json:"totalPage"`
	PageNum     int                    `json:"pageNum"`
	PageSize    int                    `json:"pageSize"`
}
type getEvaluationResp struct {
	Evaluation services.Evaluation `json:"evaluation"`
	Log        string              `json:"log"`
	Indicator  map[string]string   `json:"indicator"`
	Confusion  map[string]string   `json:"confusion"`
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
	var req getEvaluationsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	evaluations, total, totalPage, err := services.GetEvaluations(username, req.PageNum, req.PageSize,
		req.Status, req.Search, req.OrderBy, req.Order)
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
// @Param param body services.Evaluation true "name:model name ， datasetName ：dataset name"
// @Success 200 {object} createEvaluationResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/evaluations [post]
func createEvaluation(c *gin.Context) error {
	var req services.Evaluation

	err := c.BindJSON(&req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	////检查数据集文件是否存在
	//if req.DatasetPath != "" {
	//	err = services.CheckPathExists(req.DatasetPath)
	//	if err != nil {
	//		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//	}
	//}
	//
	////检查模型参数文件是否存在
	//err = services.CheckPathExists(req.CodePath)
	//if err != nil {
	//	return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	//}

	////检查输出路径是否存在自动去创建
	////err = services.CheckPathExists(req.OutputPath)
	////if err != nil {
	////	return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	////}
	////
	jobId, err := services.CreateEvaluation(username, req)
	if err != nil {
		return AppError(CREATE_EVALUATION_FAILED_CODE, err.Error())
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
	indicator,confusion  := services.GetRegexpLog(logResp)
	data := getEvaluationResp{
		Evaluation: *job,
		Log:        logResp,
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
		return AppError(NO_USRNAME, "no username")
	}
	err = services.DeleteEvaluation(username, id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
