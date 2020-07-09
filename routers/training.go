package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupTraining(r *gin.Engine) {
	group := r.Group("/ai_arts/api/trainings")

	group.GET("/", wrapper(getAllTraining))
	group.GET("/:id", wrapper(getTraining))
	group.POST("/", wrapper(createTraining))
	//group.DELETE("/:id", wrapper(stopTraining))
}

type GetAllTrainingReq struct {
	PageNum 	int 	`json:"pageNum"`
	PageSize 	int 	`json:"pageSize"`
}

type GetAllTrainingRsp struct {
	Trainings 	[] *models.Training `json:"Trainings"`
	Total		int   	`json:"total"`
	totalPage	int 	`json:"totalPage"`
}

type CreateTrainingReq struct {
	Name 			string `json:"name"`
	Desc 			string `json:"desc"`
	FrameworkInfo 	models.AIFrameworkItem `json:"frameworkInfo"`
	CodePath		string `json:"codePath"`
	StartupFile		string `json:"startupFile"`
	OutputPath		string `json:"outputPath"`
	DatasetPath		string `json:"datasetPath"`
	Params 			map[string]string `json:"params"`
}

type CreateTrainingRsp struct {
	Id 				string `json:"id"`
}

type DeleteTrainingReq struct {
	Id 				string `json:"id"`
}

type DeleteTrainingRsp struct {

}

type GetTrainingReq struct {
	Id 				string `json:"id"`
}

type GetTrainingRsp struct {

}

// @Summary list Trainings
// @Produce  json
// @Param page query int true "page number"
// @Param pagesize query int true "size per page"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/training [get]
func getAllTraining(c *gin.Context) error {

	var req GetAllTrainingReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	sets, total, totalPage, err := services.GetAllTraining(req.PageNum, req.PageSize)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetAllTrainingRsp {
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary get Training
// @Produce  json
// @Param name query string true "dataset name"
// @Param description query string true "dataset description"
// @Param creator query string true "dataset creator"
// @Param path query string true "dataset storage path"
// @Success 200 {object} APISuccessRespCreateTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/training [post]
func createTraining(c *gin.Context) error {

	var req CreateTrainingReq
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	id, err = services.CreateTraining(req.Name, "", req.FrameworkInfo)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary create Training
// @Produce  json
// @Param name query string true "dataset name"
// @Param description query string true "dataset description"
// @Param creator query string true "dataset creator"
// @Param path query string true "dataset storage path"
// @Success 200 {object} APISuccessRespGetTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/training/:id [get]
func getTraining(c *gin.Context) error {

	var req GetTrainingReq
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	err = services.GetTraining(id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, nil)
}

// @Summary delete Training
// @Produce  json
// @Param description query string true "dataset description"
// @Success 200 {object} APISuccessRespDeleteTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/training/:id [post]
func delTraining(c *gin.Context) error {

	var id string
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	var req DeleteTrainingReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	err = services.DeleteTraining(id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}