package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/apulis/AIArtsBackend/models"
	"strings"
	"github.com/apulis/AIArtsBackend/configs"
)

type UpdateProjectParams struct {

}

func AddGroupAnnotation(r *gin.Engine) {
	group := r.Group("/ai_arts/api/annotations")

	group.GET("/projects", wrapper(GetProjects))
	group.DELETE("/projects/:projectId", wrapper(DeleteProject))
	group.POST("/projects", wrapper(AddProject))
	group.PATCH("/projects/:projectId", wrapper(UpdateProject))
	group.GET("/projects/:projectId/datasets", wrapper(GetDatasets))
	group.POST("/projects/:projectId/datasets", wrapper(AddDataset))
	group.GET("/projects/:projectId/datasets/:dataSetId", wrapper(GetDatasetInfo))
	group.PATCH("/projects/:projectId/datasets/:dataSetId", wrapper(UpdateDataSet))
	group.DELETE("/projects/:projectId/datasets", wrapper(RemoveDataSet))
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks", wrapper(GetTasks))
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks/next/:taskId", wrapper(GetNextTask))
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId", wrapper(GetOneTask))
	group.POST("/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId", wrapper(PostOneTask))
	group.POST("/projects/:projectId/datasets/:dataSetId/tasks/labels", wrapper(GetDataSetLabels))
}



// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/annotations [post]
func GetProjects(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	logger.Info("token is ",token)
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	projects,totalCount,err := services.GetProjects(queryStringParameters)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	return SuccessResp(c,gin.H{"projects":projects,"totalCount":totalCount})
}

func DeleteProject(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	err := services.DeleteProject(projectId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func AddProject(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	var params models.Project
	err := c.ShouldBind(&params)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.AddProject(params)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func UpdateProject(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	var project models.Project
	projectId := c.Param("projectId")
	err := c.ShouldBind(&project)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateProject(project,projectId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func GetDatasets(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	datasets,totalCount,err := services.GetDatasets(projectId,queryStringParameters)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"datasets":datasets,"totalCount":totalCount})
}

func AddDataset(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	var dataset models.UpdateDataSet
	projectId := c.Param("projectId")
	err := c.ShouldBind(&dataset)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.AddDataset(projectId,dataset)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func GetDatasetInfo(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	dataset,err := services.GetDatasetInfo(projectId,dataSetId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"info":dataset})
}

func UpdateDataSet(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	var dataset models.UpdateDataSet
	err := c.ShouldBind(&dataset)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateDataSet(projectId,dataSetId,dataset)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func RemoveDataSet(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	var dataSetId string
	err := c.ShouldBind(&dataSetId)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.RemoveDataSet(projectId,dataSetId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func GetTasks(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	tasks,totalCount,err := services.GetTasks(projectId,dataSetId,queryStringParameters)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"taskList":tasks,"totalCount":totalCount})
}

func GetNextTask(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	nextTask,err := services.GetNextTask(projectId,dataSetId,taskId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"next":nextTask})
}

func GetOneTask(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	taskObj,err := services.GetOneTask(projectId,dataSetId,taskId)
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"annotations":taskObj})
}

func PostOneTask(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	value,_ := c.GetRawData()
	err := services.PostOneTask(projectId,dataSetId,taskId,string(value))
	if err != nil {
		return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{})
}

func GetDataSetLabels(c *gin.Context) error {
	token := c.GetHeader("Authorization")
	token = strings.Split(token,"Bearer ")[1]
	configs.Config.Token = token
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	labels,err := services.GetDataSetLabels(projectId,dataSetId)
	if err != nil {
	return AppError(APP_ERROR_CODE,err.Error())
	}
	return SuccessResp(c,gin.H{"annotations":labels})
}