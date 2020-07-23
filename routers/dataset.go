package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupDataset(r *gin.Engine) {
	group := r.Group("/ai_arts/api/datasets")
	group.Use(Auth())
	group.GET("/", wrapper(lsDatasets))
	group.GET("/:id", wrapper(getDataset))
	group.POST("/", wrapper(createDataset))
	group.POST("/:id", wrapper(updateDataset))
	group.DELETE("/:id", wrapper(deleteDataset))
	group.POST("/:id/bind", wrapper(bindDataset))
	group.POST("/:id/unbind", wrapper(unbindDataset))

}

type datasetId struct {
	ID int `uri:"id" binding:"required"`
}

type lsDatasetsReq struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize,default=10"`
	Name     string `form:"name"`
}

type createDatasetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Path        string `json:"path" binding:"required"`
	IsPrivate   bool   `json:"isPrivate"`
}
type bindDatasetReq struct {
	Platform string `json:"platform" binding:"required"`
	Id       string `json:"id" binding:"required"`
}

type updateDatasetReq struct {
	Description string `json:"description" binding:"required"`
}

type GetDatasetResp struct {
	Dataset models.Dataset `json:"dataset"`
}

type GetDatasetsResp struct {
	Datasets  []models.Dataset `json:"datasets"`
	Total     int              `json:"total"`
	TotalPage int              `json:"totalPage"`
	PageNum   int              `json:"pageNum"`
	PageSize  int              `json:"pageSize"`
}

// @Summary list datasets
// @Produce  json
// @Param pageNum query int true "page number, from 1"
// @Param pageSize query int true "count per page"
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
	var datasets []models.Dataset
	var total int
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	if req.Name == "" {
		datasets, total, err = services.ListDatasets(req.PageNum, req.PageSize, username)
	} else {
		datasets, total, err = services.ListDatasetsByName(req.PageNum, req.PageSize, req.Name, username)
	}
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := GetDatasetsResp{
		Datasets:  datasets,
		Total:     total,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		TotalPage: total/req.PageSize + 1,
	}
	return SuccessResp(c, data)
}

// @Summary get dataset by id
// @Produce  json
// @Param id path int true "dataset id"
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
// @Param body body createDatasetReq true "json body"
// @Param description body string true "dataset description"
// @Param path body string true "dataset storage path"
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
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.CreateDataset(req.Name, req.Description, username, "0.0.1", req.Path, req.IsPrivate)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary update dataset
// @Produce  json
// @Param description path string true "dataset description"
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
// @Param id path int true "dataset id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id [delete]
func deleteDataset(c *gin.Context) error {
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

// @Summary bind dataset
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"  "code": 30000, "already bind"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id/bind  [post]
func bindDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req bindDatasetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	err = services.BindDataset(id.ID, req.Platform, req.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary unbind dataset
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"  "code": 30000, "no bind"
// @Router /ai_arts/api/datasets/:id/unbind  [post]
func unbindDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	var req bindDatasetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.UnbindDataset(id.ID, req.Platform, req.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
