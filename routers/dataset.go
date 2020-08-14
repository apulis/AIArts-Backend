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
	PageNum      int    `form:"pageNum,default=1"`
	PageSize     int    `form:"pageSize,default=10"`
	Name         string `form:"name"`
	Status       string `form:"status"`
	OrderBy      string `form:"orderBy,default=updated_at"`
	Order        string `form:"order,default=desc"`
	IsTranslated bool   `form:"isTranslated"`
}

type createDatasetReq struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Path         string `json:"path" binding:"required"`
	IsPrivate    bool   `json:"isPrivate" `
	IsTranslated bool   `json:"isTranslated" `
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
	Message   string           `json:"message"`
}

// @Summary list datasets
// @Produce  json
// @Param query query lsDatasetsReq true "isUsable 返回用户可用的数据集"
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets [get]
func lsDatasets(c *gin.Context) error {
	models.GinContext{Context: c}.SaveToken()
	var req lsDatasetsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	var datasets []models.Dataset
	var total = 0
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	//获取该用户能够访问的所有已经标注好的数据库
	var message = "success"
	datasets, total, err = services.ListDatasets(req.PageNum, req.PageSize, req.OrderBy, req.Order, req.Name, req.Status, req.IsTranslated, username)

	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	if req.IsTranslated {
		var annoDatasets []models.DataSet
		queryStringParameters := models.QueryStringParametersV2{
			PageNum:  req.PageNum,
			PageSize: req.PageSize,
			OrderBy:  req.OrderBy,
			Order:    req.Order,
		}
		annoDatasets, _, err := services.ListAllDatasets(queryStringParameters)
		if err != nil {
			message = "label image platform is error"
			//return AppError(FAILED_FETCH_ANNOTATION_CODE, "label plantform is error")
		} else {
			for _, v := range annoDatasets {
				if v.ConvertStatus == "finished" {
					dataset := models.Dataset{
						Name:        v.Name,
						Description: v.Info,
						Path:        v.ConvertOutPath,
						Status:      v.Name,
						//是否是公开数据集
						IsPrivate:    v.IsPrivate,
						IsTranslated: true,
					}
					datasets = append(datasets, dataset)
					total += 1

				}
			}
		}
	}


	data := GetDatasetsResp{
		Datasets:  datasets,
		Total:     total,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		TotalPage: total/req.PageSize + 1,
		Message:   message,
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
// @Param name body string false "dataset name"
// @Param IsPrivate body bool false "dataset auth"
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
	//因为会上传到home目录下
	//err = services.CheckDatasetPathValid(req.Path)
	//if err != nil {
	//	return AppError(FILEPATH_NOT_VALID_CODE, err.Error())
	//}
	err = services.CheckPathExists(req.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.CreateDataset(req.Name, req.Description, username, "0.0.1", req.Path, req.IsPrivate, req.IsTranslated)
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
// @Failure 30010 {object} APIException "dataset is still using"
// @Router /ai_arts/api/datasets/:id [delete]
func deleteDataset(c *gin.Context) error {
	var id datasetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.DeleteDataset(id.ID)
	if err != nil {
		return AppError(DATASET_IS_STILL_USE_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary bind dataset
// @Produce  json
// @Param platform body string true "bind platform's name"
// @Param id body string true "bind platform's id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
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
// @Param platform body string true "bind platform's name"
// @Param id body string true "bind platform's id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
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
