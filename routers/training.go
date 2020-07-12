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
	Training models.Training
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

// @Summary get all trainings
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings [get]
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

// @Summary create Training
// @Produce json
// @Param param body CreateTrainingReq true "params"
// @Success 200 {object} APISuccessRespCreateTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings [post]
func createTraining(c *gin.Context) error {

	var req CreateTrainingReq
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	id, err = services.CreateTraining(req.Training)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary get specific training
// @Produce  json
// @Param param body GetTrainingReq true "params"
// @Success 200 {object} APISuccessRespGetTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings/:id [get]
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

// @Summary delete one training
// @Produce  json
// @Param param body DeleteTrainingReq true "params"
// @Success 200 {object} APISuccessRespDeleteTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/trainings/:id [delete]
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