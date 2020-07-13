package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/apulis/AIArtsBackend/services"

)

func AddGroupInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/inferences")

	group.POST("/PostInferenceJob", wrapper(PostInferenceJob))
	group.GET("/ListInferenceJob", wrapper(ListInferenceJob))
	group.GET("/GetAllSupportInference", wrapper(GetAllSupportInference))
	group.GET("/GetAllDevice", wrapper(GetAllDevice))
	group.GET("/GetJobDetail", wrapper(GetJobDetail))
	group.GET("/GetJobLog", wrapper(GetJobLog))
	group.GET("/GetJobStatus", wrapper(GetJobStatus))
	group.POST("/Infer", wrapper(Infer))
}

// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/inferences [post]
func PostInferenceJob(c *gin.Context) error {
	var params models.PostInference
	err := c.ShouldBind(&params)
	params.UserName = getUsername(c)
	params.UserId = getUserId(c)
	if params.VcName=="" {
		params.VcName = "platform"
	}
	if err != nil {
		return ParameterError(err.Error())
	}

	jobId,err := services.PostInferenceJob(params)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"jobId":jobId})
}

func ListInferenceJob(c *gin.Context) error {
	vcName := c.Query("vcName")
	if vcName=="" {
		vcName = "platform"
	}
	//jobOwner := c.Query("jobOwner")
	jobOwner := getUsername(c)
	var queryStringParameters models.QueryStringParametersV2
	err := c.ShouldBindQuery(&queryStringParameters)
	jobs,err := services.ListInferenceJob(jobOwner,vcName,queryStringParameters)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetAllSupportInference(c *gin.Context) error {
	inferences,err := services.GetAllSupportInference()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, inferences)
}

func GetAllDevice(c *gin.Context) error {
	userName := getUsername(c)
	jobs,err := services.GetAllDevice(userName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetJobDetail(c *gin.Context) error {
	userName := getUsername(c)
	jobId := c.Query("jobId")
	jobs,err := services.GetJobDetail(userName,jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetJobLog(c *gin.Context) error {
	userName := getUsername(c)
	jobId := c.Query("jobId")
	jobs,err := services.GetJobLog(userName,jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetJobStatus(c *gin.Context) error {
	jobId := c.Query("jobId")
	jobs,err := services.GetJobStatus(jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func Infer(c *gin.Context) error {
	jobId := c.Query("jobId")
	file, err := c.FormFile("image")
	err = c.SaveUploadedFile(file, "./"+jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	resp,err := services.Infer(jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, resp)
}