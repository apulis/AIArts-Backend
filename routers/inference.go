package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/inferences")

	group.Use(Auth())

	group.POST("/PostInferenceJob", wrapper(PostInferenceJob))
	group.GET("/ListInferenceJob", wrapper(ListInferenceJob))
	group.GET("/GetAllSupportInference", wrapper(GetAllSupportInference))
	group.GET("/GetAllDevice", wrapper(GetAllDevice))
	group.GET("/GetJobDetail", wrapper(GetJobDetail))
	group.GET("/GetJobLog", wrapper(GetJobLog))
	group.GET("/GetJobStatus", wrapper(GetJobStatus))
	group.POST("/Infer", wrapper(Infer))
	group.GET("/KillJob", wrapper(KillJob))
}

// @Summary submit a inference job
// @Description submit a inference job
// @Produce  json
// @Param body body models.PostInference true "json body"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/PostInferenceJob [post]
func PostInferenceJob(c *gin.Context) error {
	var params models.PostInference
	err := c.ShouldBind(&params)
	params.UserName = getUsername(c)
	params.UserId = getUserId(c)
	if params.VcName == "" {
		params.VcName = "platform"
	}
	if err != nil {
		return ParameterError(err.Error())
	}

	jobId, err := services.PostInferenceJob(params)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"jobId": jobId})
}

// @Summary list inference jobs
// @Description list inference jobs
// @Produce  json
// @Param vcName query string false "which virtual cluster"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/ListInferenceJob [get]
func ListInferenceJob(c *gin.Context) error {
	vcName := c.Query("vcName")
	if vcName == "" {
		vcName = "platform"
	}
	//jobOwner := c.Query("jobOwner")
	jobOwner := getUsername(c)
	var queryStringParameters models.QueryStringParametersV2
	err := c.ShouldBindQuery(&queryStringParameters)
	jobs, err := services.ListInferenceJob(jobOwner, vcName, queryStringParameters)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

// @Summary get all support inference framework\device
// @Description list inference jobs
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/GetAllSupportInference [get]
func GetAllSupportInference(c *gin.Context) error {
	inferences, err := services.GetAllSupportInference()
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, inferences)
}

// @Summary get all device type detail
// @Description get all device type detail
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/GetAllDevice [get]
func GetAllDevice(c *gin.Context) error {
	userName := getUsername(c)
	jobs, err := services.GetAllDevice(userName)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

// @Summary get inference job detail
// @Description get inference job detail
// @Produce  json
// @Param jobId query string true "inference job Id "
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/GetJobDetail [get]
func GetJobDetail(c *gin.Context) error {
	userName := getUsername(c)
	jobId := c.Query("jobId")
	jobs, err := services.GetJobDetail(userName, jobId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

// @Summary get inference job log
// @Description get inference job log
// @Produce  json
// @Param jobId query string true "inference job Id "
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/GetJobLog [get]
func GetJobLog(c *gin.Context) error {
	userName := getUsername(c)
	jobId := c.Query("jobId")
	jobs, err := services.GetJobLog(userName, jobId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

// @Summary get inference job status
// @Description get inference job status
// @Produce  json
// @Param jobId query string true "inference job Id "
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/GetJobStatus [get]
func GetJobStatus(c *gin.Context) error {
	jobId := c.Query("jobId")
	jobs, err := services.GetJobStatus(jobId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

// @Summary Infer a picture using a running inference job
// @Description Infer a picture using a running inference job
// @Produce  json
// @Param jobId query string true "inference job Id "
// @Param image body byte true "picture upload to infer"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/Infer [post]
func Infer(c *gin.Context) error {
	jobId := c.Query("jobId")
	signature_name := c.Query("signature_name")
	file, err := c.FormFile("image")
	err = c.SaveUploadedFile(file, "./"+jobId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	resp, err := services.Infer(jobId, signature_name)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, resp)
}

// @Summary kill a running inference job
// @Description kill a running inference job
// @Produce  json
// @Param jobId query string true "inference job Id "
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/inferences/KillJob [get]
func KillJob(c *gin.Context) error {
	jobId := c.Query("jobId")
	userName := getUsername(c)
	resp, err := services.KillJob(jobId, userName)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, resp)
}
