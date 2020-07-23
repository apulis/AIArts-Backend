package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

type UpdateProjectParams struct {
}

func AddGroupAnnotation(r *gin.Engine) {
	group := r.Group("/ai_arts/api/annotations")

	group.Use(Auth())

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
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks/previous/:taskId", wrapper(GetPreviousTask))
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId", wrapper(GetOneTask))
	group.POST("/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId", wrapper(PostOneTask))
	group.GET("/projects/:projectId/datasets/:dataSetId/tasks/labels", wrapper(GetDataSetLabels))
	group.POST("/projects/:projectId/datasets/:dataSetId/ConvertDataFormat", wrapper(ConvertDataFormat))
	group.GET("/projects/:projectId/datasets/:dataSetId/ConvertSupportFormat", wrapper(ConvertSupportFormat))
	group.GET("/datasets", wrapper(ListAllDatasets))
}




// @Summary list projects
// @Description get projects of data-platform
// @Produce  json
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects [get]
func GetProjects(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	logger.Info(queryStringParameters)
	projects, totalCount, err := services.GetProjects(queryStringParameters)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"projects": projects, "totalCount": totalCount})
}

// @Summary delete project
// @Description delete project of data-platform
// @Produce  json
// @Param projectId path string true "project id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId [delete]
func DeleteProject(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	err := services.DeleteProject(projectId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary add project
// @Description add project of data-platform
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects [post]
func AddProject(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var params models.Project
	err := c.ShouldBind(&params)
	if err != nil {
		return ParameterError(err.Error())
	}
	params.Creator = getUsername(c)
	err = services.AddProject(params)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary update projects
// @Description update project of data-platform
// @Produce  json
// @Param projectId path string true "project id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId [patch]
func UpdateProject(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var project models.Project
	projectId := c.Param("projectId")
	err := c.ShouldBind(&project)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateProject(project, projectId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary list datasets
// @Description list datasets of data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets [get]
func GetDatasets(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	datasets, totalCount, err := services.GetDatasets(projectId, queryStringParameters)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"datasets": datasets, "totalCount": totalCount})
}

// @Summary add dataset
// @Description add dataset for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets [post]
func AddDataset(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var dataset models.UpdateDataSet
	projectId := c.Param("projectId")
	err := c.ShouldBind(&dataset)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.AddDataset(projectId, dataset)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary get dataset info
// @Description get dataset info for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId [get]
func GetDatasetInfo(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	dataset, err := services.GetDatasetInfo(projectId, dataSetId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"info": dataset})
}

// @Summary update dataset info
// @Description update dataset info for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId [patch]
func UpdateDataSet(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	var dataset models.UpdateDataSet
	err := c.ShouldBind(&dataset)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateDataSet(projectId, dataSetId, dataset)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary delete dataset
// @Description delete dataset info for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId body string true "dataSet id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets [delete]
func RemoveDataSet(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	var dataSetId string
	err := c.ShouldBind(&dataSetId)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.RemoveDataSet(projectId, dataSetId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary get dataset tasks
// @Description get dataset tasks for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks [get]
func GetTasks(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	var queryStringParameters models.QueryStringParameters
	err := c.ShouldBindQuery(&queryStringParameters)
	tasks, totalCount, err := services.GetTasks(projectId, dataSetId, queryStringParameters)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"taskList": tasks, "totalCount": totalCount})
}

// @Summary get dataset next task id
// @Description get dataset next task id for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param taskId path string true "current task id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks/next/:taskId [get]
func GetNextTask(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	nextTask, err := services.GetNextTask(projectId, dataSetId, taskId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"next": nextTask})
}

// @Summary get dataset previous task id
// @Description get dataset previous task id for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param taskId path string true "current task id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks/previous/:taskId [get]
func GetPreviousTask(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	nextTask, err := services.GetPreviousTask(projectId, dataSetId, taskId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"previous": nextTask})
}

// @Summary get dataset one task detail
// @Description get dataset one task detail for data-platform project
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param taskId path string true "current task id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId [get]
func GetOneTask(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	taskObj, err := services.GetOneTask(projectId, dataSetId, taskId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"annotations": taskObj})
}

// @Summary commit label data to one task
// @Description commit label data to one task
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param taskId path string true "current task id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks/annotations/:taskId [post]
func PostOneTask(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	taskId := c.Param("taskId")
	value, _ := c.GetRawData()
	err := services.PostOneTask(projectId, dataSetId, taskId, string(value))
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{})
}

// @Summary get dataset all labels
// @Description get dataset all labels
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/tasks/labels [get]
func GetDataSetLabels(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	labels, err := services.GetDataSetLabels(projectId, dataSetId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"annotations": labels})
}

// @Summary convert a dataset to specific format
// @Description convert a dataset to specific format
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Param type body string true "dataset type,like image"
// @Param target body string true "convert to specific format"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/ConvertDataFormat [post]
func ConvertDataFormat(c *gin.Context) error {
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	var convert models.ConvertDataFormat
	err := c.ShouldBind(&convert)
	if err != nil {
		return ParameterError(err.Error())
	}
	convert.DatasetId = dataSetId
	convert.ProjectId = projectId
	ret, err := services.ConvertDataFormat(convert)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, ret)
}

// @Summary current support convert's specific format
// @Description current support convert's specific format
// @Produce  json
// @Param projectId path string true "project id"
// @Param dataSetId path string true "dataSet id"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/projects/:projectId/datasets/:dataSetId/ConvertSupportFormat [get]
func ConvertSupportFormat(c *gin.Context) error {
	projectId := c.Param("projectId")
	dataSetId := c.Param("dataSetId")
	ret, err := services.ConvertSupportFormat(projectId,dataSetId)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, ret)
}

// @Summary list all datasets for user
// @Description list all datasets for user
// @Produce  json
// @Param pageNum query int false "page number, from 1"
// @Param pageSize query int false "count per page"
// @Success 200 {object} APISuccessResp "success"
// @Router /ai_arts/api/annotations/datasets [get]
func ListAllDatasets(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var queryStringParameters models.QueryStringParametersV2
	err := c.ShouldBindQuery(&queryStringParameters)
	datasets, totalCount, err := services.ListAllDatasets(queryStringParameters)
	if err != nil {
		return ServeError(REMOTE_SERVE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, gin.H{"datasets": datasets, "totalCount": totalCount})
}