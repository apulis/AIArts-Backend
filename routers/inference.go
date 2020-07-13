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
	group.GET("/Infer", wrapper(Infer))
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
	num := c.Query("num")
	vcName := c.Query("vcName")
	jobOwner := c.Query("jobOwner")
	jobs,err := services.ListInferenceJob(jobOwner,vcName,num)
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
	userName := c.Query("userName")
	jobs,err := services.GetAllDevice(userName)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetJobDetail(c *gin.Context) error {
	userName := c.Query("userName")
	jobId := c.Query("jobId")
	jobs,err := services.GetJobDetail(userName,jobId)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, jobs)
}

func GetJobLog(c *gin.Context) error {
	userName := c.Query("userName")
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
	image, err := c.FormFile("image")
	resp,err := services.Infer(jobId,image)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, resp)
}