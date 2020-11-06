package routers

import (
	//"fmt"
	//"github.com/apulis/AIArtsBackend/models"
	//"github.com/apulis/AIArtsBackend/services"
	"strconv"

	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)


//AddGroupExperimentMgr
func AddGroupExperimentMgr(r *gin.Engine) {

	group := r.Group("/ai_arts/api/projects")
	//group.Use(Auth())

	group.GET("/", wrapper(getAllExpProjects))
	group.GET("/:id", wrapper(getExpProject))
	group.PUT("/:id", wrapper(updateExpProject))
	group.POST("/",wrapper(createExpProject))
	group.POST("/:id", wrapper(postExpProject))

	group = r.Group("/ai_arts/api/experiments")
	//group.Use(Auth())

	group.GET("/", wrapper(getAllExperiments))
	group.GET("/:id", wrapper(getExperiment))
	group.PUT("/:id", wrapper(updateExperiment))
	group.POST("/", wrapper(createExperiment))
	group.POST("/:id", wrapper(postExperiment))


}

func doRespWith(c*gin.Context,err error,data interface{}) error{
    if err != nil{
    	return err
	}
	if data == nil {
		return SuccessResp(c,gin.H{})
	}
	return SuccessResp(c,data)
}

// @Summary get all experiments projects
// @Produce  json
// @Param pageNum query int true  "page number"
// @Param pageSize query int true "size per page"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException   "error"
// @Failure 404 {object} APIException   "not found"
// @Router /ai_arts/api/projects [get]
func getAllExpProjects(c *gin.Context) error {

	userName := getUsername(c)

	var queryParams services.CommQueryParams
	queryParams.SetQueryParams(c)

	resultSet, total, err := services.ListExpProjects(queryParams, userName)

	return doRespWith(c,err,gin.H{
		"items": resultSet,
		"total": total,
	})
}

func getExpProject(c *gin.Context) error {
	var project models.ExpProject
	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	err := services.QueryExpProject(id,&project)
	return doRespWith(c,err,&project)
}

func updateExpProject(c *gin.Context) error {

	new_name := c.Query("new_name")
	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	if c.Request.ContentLength == 0 {
		return doRespWith(c,services.UpdateExpProject(id,new_name,nil),nil)
	}else{
		var project models.RequestUpdates
		err := c.ShouldBindJSON(&project)
		if err != nil {
			return ParameterError(err.Error())
		}
		return doRespWith(c, services.UpdateExpProject(id,new_name,&project),nil)
	}

}

func createExpProject(c *gin.Context) error {

	var project models.ExpProject
	err := c.ShouldBindJSON(&project)
	if err != nil {
		return ParameterError(err.Error())
	}
	project.Creator=getUsername(c)
	if len(project.Name) == 0 {
		return AppError(APP_ERROR_CODE,"name cannot be empty")
	}
	err = services.CreateExpProject(&project)
	return doRespWith(c,err,gin.H{ "id":project.ID	})
}

func postExpProject(c *gin.Context) error {

	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	action := c.Query("action")
	switch action {
		case "delete":
			return doRespWith(c,services.MarkExpProject(id, true),nil)
		case "restore":
			return doRespWith(c,services.MarkExpProject(id, false),nil)
		default:
			return AppError(APP_ERROR_CODE, "Unsupport action !!!")
	}

}
// @Summary get all experiments
// @Produce  json
// @Param pageNum query int true  "page number"
// @Param pageSize query int true "size per page"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException   "error"
// @Failure 404 {object} APIException   "not found"
// @Router /ai_arts/api/experiments?project=n [get]
func getAllExperiments(c *gin.Context) error {

	projectID,_ := strconv.ParseUint(c.Query("project"),0,0)
	if projectID == 0 {
		return AppError(PARAMETER_ERROR_CODE,"Invalid project ID !")
	}

	var queryParams services.CommQueryParams
	queryParams.SetQueryParams(c)

	resultSet, total, err := services.ListExperiments(queryParams, projectID)

	return doRespWith(c,err,gin.H{
		"items": resultSet,
		"total": total,
	})
}
func getExperiment(c *gin.Context) error {
	var experiment models.Experiment
	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	err := services.QueryExperiment(id,&experiment)
	return doRespWith(c,err,&experiment)
}
func postExperiment(c *gin.Context) error {
	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	action := c.Query("action")
	switch action {
	case "delete":
		return doRespWith(c,services.MarkExperiment(id, true),nil)
	case "restore":
		return doRespWith(c,services.MarkExperiment(id, false),nil)
	default:
		return AppError(APP_ERROR_CODE, "Unsupport action !!!")
	}
}
func updateExperiment(c *gin.Context) error {

	new_name := c.Query("new_name")
	id , _:=  strconv.ParseUint(c.Param("id"),0,0)
	if c.Request.ContentLength == 0 {
		return doRespWith(c,services.UpdateExperiment(id,new_name,nil),nil)
	}else{
		var experiment models.RequestUpdates
		err := c.ShouldBindJSON(&experiment)
		if err != nil {
			return ParameterError(err.Error())
		}
		return doRespWith(c, services.UpdateExperiment(id,new_name,&experiment),nil)
	}

}
func createExperiment(c*gin.Context)error{
	projectID,_ := strconv.ParseUint(c.Query("project"),0,0)
	if projectID == 0 {
		return AppError(PARAMETER_ERROR_CODE,"Invalid project ID !")
	}

	var experiment models.Experiment
	err := c.ShouldBindJSON(&experiment)
	if err != nil {
		return ParameterError(err.Error())
	}
	experiment.Creator=getUsername(c)
	if len(experiment.Name) == 0 {
		return AppError(APP_ERROR_CODE,"name cannot be empty")
	}
	experiment.ProjectID=uint(projectID)
	err = services.CreateExperiment(&experiment)
	return doRespWith(c,err,gin.H{ "id":experiment.ID	})
}
