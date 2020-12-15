package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupJobManager(r *gin.Engine) {
	group := r.Group("/ai_arts/api/jobs/")
	group.Use(Auth())
	group.GET("/", wrapper(GetAllJobs))
	group.GET("/summary", wrapper(GetAllJobSummary))
	group.GET("/resume/:jobId", wrapper(ResumeJob))
	group.GET("/pause/:jobId", wrapper(PauseJob))
}

func GetAllJobs(c *gin.Context) error {
	var req models.GetAllJobsReq
	var err error

	//todo: user authorization
	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	count, err := services.GetJobsCount(req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	jobs, err := services.GetAllJobs(req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, map[string]interface{}{"jobs": jobs, "total": count})
}

func GetAllJobSummary(c *gin.Context) error {
	var err error
	var req GetJobSummaryReq

	//todo: user authorization
	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}


	summary, err := services.GetJobSummary(userName, req.JobType, req.VCName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, summary)
}

func ResumeJob(c *gin.Context) error {
	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	jobId := c.Param("jobId")

	ret, err := services.ResumeJob(jobId, userName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, ret)
}

func PauseJob(c *gin.Context) error {
	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	jobId := c.Param("jobId")

	ret, err := services.PauseJob(jobId, userName)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, ret)
}
