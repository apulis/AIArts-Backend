package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupModel(r *gin.Engine) {
	group := r.Group("/ai_arts/api/models")

	group.GET("/", wrapper(lsModelsets))
	group.GET("/:id", wrapper(getModelset))
	group.POST("/", wrapper(createModelset))
	group.POST("/:id", wrapper(updateModelset))
	group.DELETE("/:id", wrapper(deleteModelset))
}

type modelsetId struct {
	ID int `uri:"id" binding:"required"`
}

type lsModelsetsReq struct {
	Page  int `form:"page"`
	Count int `form:"count,default=10"`
}

type createModelsetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Creator     string `json:"creator" binding:"required"`
	Path        string `json:"path" binding:"required"`
}

type updateModelsetReq struct {
	Description string `json:"description" binding:"required"`
}

type GetModelsetResp struct {
	Model models.Modelset `json:"model"`
}

type GetModelsetsResp struct {
	Models []models.Modelset `json:"models"`
	Total  int               `json:"total"`
	Page   int               `json:"page"`
	Count  int               `json:"count"`
}

// @Summary list models
// @Produce  json
// @Param page query int true "page number, from 1"
// @Param count query int true "count per page"
// @Success 200 {object} APISuccessRespGetModelsets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [get]
func lsModelsets(c *gin.Context) error {
	var req lsModelsetsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelsets, total, err := services.ListModelSets(req.Page, req.Count)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := GetModelsetsResp{
		Models: modelsets,
		Total:  total,
		Page:   req.Page,
		Count:  req.Count,
	}
	return SuccessResp(c, data)
}

// @Summary get model by id
// @Produce  json
// @Param id path int true "dataset id"
// @Success 200 {object} APISuccessRespGetModelset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [get]
func getModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelset, err := services.GetModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := GetModelsetResp{Model: modelset}
	return SuccessResp(c, data)
}

// @Summary create model
// @Produce  json
// @Param name body string true "model name"
// @Param description body string true "model description"
// @Param creator body string true "model creator"
// @Param path body string true "model storage path"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models [post]
func createModelset(c *gin.Context) error {
	var req createModelsetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.CheckPathExists(req.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	err = services.CreateModelset(req.Name, req.Description, req.Creator, "0.0.1", req.Path)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary update model
// @Produce  json
// @Param description path string true "model description"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [post]
func updateModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req updateModelsetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateModelset(id.ID, req.Description)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary delete model by id
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/models/:id [delete]
func deleteModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.DeleteModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
