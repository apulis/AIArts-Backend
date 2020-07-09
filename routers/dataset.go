package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupDataset(r *gin.Engine) {
	group := r.Group("/ai_arts/api/datasets")

	group.GET("/", wrapper(lsDatasets))
	group.GET("/:id", wrapper(getDataset))
	group.POST("/", wrapper(createDataset))
	group.POST("/:id", wrapper(updateDataset))
	group.DELETE("/:id", wrapper(DeleteDataset))
}

type datasetId struct {
	ID int `uri:"id" binding:"required"`
}

type lsDatasetsReq struct {
	Page  int `form:"page"`
	Count int `form:"count,default=10"`
}

type createDatasetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Creator     string `json:"creator" binding:"required"`
	Path        string `json:"path" binding:"required"`
}

type updateDatasetReq struct {
	Description string `json:"description" binding:"required"`
}

type GetDatasetResp struct {
	Dataset models.Dataset `json:"dataset"`
}

type GetDatasetsResp struct {
	Datasets []models.Dataset `json:"datasets"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Count    int              `json:"count"`
}

// @Summary list datasets
// @Produce  json
// @Param page query int true "page number"
// @Param count query int true "count per page"
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets [get]
func lsDatasets(c *gin.Context) error {
	var req lsDatasetsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	datasets, total, err := services.ListDatasets(req.Page, req.Count)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{
		"datasets": datasets,
		"total":    total,
		"page":     req.Page,
		"count":    req.Count,
	}
	return SuccessResp(c, data)
}

// @Summary get dataset by id
// @Produce  json
// @Param id query int true "dataset id"
// @Success 200 {object} APISuccessRespGetDataset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id [get]
func getDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	dataset, err := services.GetDataset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := GetDatasetResp{Dataset: dataset}
	return SuccessResp(c, data)
}

// @Summary create dataset
// @Produce  json
// @Param name query string true "dataset name"
// @Param description query string true "dataset description"
// @Param creator query string true "dataset creator"
// @Param path query string true "dataset storage path"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets [post]
func createDataset(c *gin.Context) error {
	var req createDatasetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.CheckDatasetPathValid(req.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_VALID_CODE, err.Error())
	}
	err = services.CheckPathExists(req.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	err = services.CreateDataset(req.Name, req.Description, req.Creator, "0.0.1", req.Path)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary update dataset
// @Produce  json
// @Param description query string true "dataset description"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id [post]
func updateDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req updateDatasetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UpdateDataset(id.ID, req.Description)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary delete dataset by id
// @Produce  json
// @Param id query int true "dataset id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id [delete]
func DeleteDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.DeleteDataset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
