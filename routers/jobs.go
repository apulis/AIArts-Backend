package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddJobManager(r *gin.Engine) {
	group := r.Group("/ai_arts/api/jobs/")
	group.Use(Auth())
	group.GET("/", wrapper(GetAllJobs))
	group.GET("/summary", wrapper(GetAllJobSummary))
}

func GetAllJobs(c *gin.Context) error {
	var req models.GetAllJobsReq
	var err error

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	// 兼容老代码
	if req.VCName == "" {
		req.VCName = "platform"
	}

	count, err := services.GetJobsCount(req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	jobs, err := services.GetAllJobs(req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, map[string]interface{}{"jobs": jobs, "total": count})
}

func GetAllJobSummary(c *gin.Context) error {
	var err error
	var req GetJobSummaryReq

	if len(req.VCName) == 0 {
		req.VCName = models.DefaultVcName
	}

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	summary, err := services.GetJobSummary(req.UserName, req.JobType, req.VCName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, summary)
}
