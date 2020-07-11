package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
	"github.com/apulis/AIArtsBackend/services"
)

func AddGroupInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/inferences")

	group.GET("/PostInferenceJob", wrapper(PostInferenceJob))
	group.GET("/ListInferenceJob", wrapper(ListInferenceJob))
	group.GET("/GetAllSupportInference", wrapper(GetAllSupportInference))
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
	err = services.PostInferenceJob(params)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
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
