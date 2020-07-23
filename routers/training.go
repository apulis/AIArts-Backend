package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupTraining(r *gin.Engine) {
	group := r.Group("/ai_arts/api/trainings")

	group.Use(Auth())

	group.GET("/", wrapper(getAllTraining))
	group.GET("/:id", wrapper(getTraining))
	group.POST("/", wrapper(createTraining))
	group.DELETE("/:id", wrapper(delTraining))
	group.GET("/:id/log", wrapper(getLog))
}

type GetAllTrainingReq struct {
	PageNum    int    `json:"pageNum"`
	PageSize   int    `json:"pageSize"`
	JobStatus  string `json:"status"`
	SearchWord string `json:"searchWord"`
}

type GetAllJobsReq struct {
	PageNum    int    `json:"pageNum"`
	PageSize   int    `json:"pageSize"`
	JobStatus  string `json:"status"`
	SearchWord string `json:"searchWord"`
}

type GetAllTrainingRsp struct {
	Trainings []*models.Training `json:"Trainings"`
	Total     int                `json:"total"`
	totalPage int                `json:"totalPage"`
}

type CreateTrainingReq struct {
	models.Training
}

type CreateTrainingRsp struct {
	Id string `json:"id"`
}

type DeleteTrainingReq struct {
	Id string `json:"id"`
}

type DeleteTrainingRsp struct {
}

type GetTrainingReq struct {
	Id string `json:"id"`
}

type GetTrainingRsp struct {
	models.Training
}

// @Summary get all trainings
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Param status query string true "job status. get all jobs if it is all"
// @Param searchWord query string true "the keyword of search"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings [get]
func getAllTraining(c *gin.Context) error {

	var req GetAllJobsReq
	var err error

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	sets, total, totalPage, err := services.GetAllTraining(userName, req.PageNum, req.PageSize, req.JobStatus, req.SearchWord)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetAllTrainingRsp{
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary create Training
// @Produce json
// @Param param body CreateTrainingReq true "params"
// @Success 200 {object} APISuccessRespCreateTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings [post]
func createTraining(c *gin.Context) error {

	var req models.Training
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	if valid, msg := req.ValidatePathByUser(userName); !valid {
		return AppError(INVALID_CODE_PATH, msg)
	}

	id, err = services.CreateTraining(userName, req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary get specific training
// @Produce  json
// @Param param body GetTrainingReq true "params"
// @Success 200 {object} APISuccessRespGetTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings/:id [get]
func getTraining(c *gin.Context) error {

	var id models.UriJobId
	var training *models.Training

	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	training, err = services.GetTraining(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, training)
}

// @Summary delete one training
// @Produce  json
// @Param param body DeleteTrainingReq true "params"
// @Success 200 {object} APISuccessRespDeleteTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings/:id [delete]
func delTraining(c *gin.Context) error {

	var id models.UriJobId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err = services.DeleteTraining(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary get specific training
// @Produce  json
// @Param param body GetTrainingReq true "params"
// @Success 200 {object} APISuccessRespGetTrainingLog "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings/log/:id [get]
func getLog(c *gin.Context) error {

	var id models.UriJobId
	var jobLog *models.JobLog

	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	jobLog, err = services.GetTrainingLog(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, jobLog)
}
