package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
	"time"
)

func AddGroupModel(r *gin.Engine) {
	group := r.Group("/ai_arts/api/models")

	group.Use(Auth())

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
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize,default=10"`
	Name     string `form:"name"`
	Status   string `form:"status"`
	IsAdvance bool `form:"isAdvance"`
}

type createModelsetReq struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description" `
	Path        string            `json:"path" binding:"required"`
	JobId       string            `json:"jobId" binding:"required"`
	Use         string            `json:"use"`
	DataFormat  string            `json:"dataFormat"`
	Arguments   map[string]string `json:"arguments"`
	EngineType  string            `json:"engineType"`
	Precision   string            `json:"precision"`
	IsAdvance bool `json:"isAdvance"`

}

type updateModelsetReq struct {
	Description string `json:"description" binding:"required"`
}

type GetModelsetResp struct {
	Model models.Modelset `json:"model"`
}
type UnixTime struct {
	time.Time
}
type JsonModel struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	IsAdvance  bool              `json:"isAdvance"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Version     string `json:"version"`
	Path        string `json:"path"`
	Status      string `json:"status"`
	Size        int    `json:"size"`
	//模型类型 计算机视觉
	Use        string            `json:"use"`
	JobId      string            `json:"jobId"`
	DataFormat string            `json:"dataFormat"`
	Arguments  map[string]string `json:"arguments,omitempty"`
	EngineType string            `json:"engineType"`
	Precision  string            `json:"precision"`
}
type GetModelsetsResp struct {
	JsonModels []models.Modelset `json:"models"`
	Total      int               `json:"total"`
	TotalPage  int               `json:"totalPage"`
	PageNum    int               `json:"pageNum"`
	PageSize   int               `json:"pageSize"`
}

// @Summary list models
// @Produce  json
// @Param body url lsModelsetsReq true "url form"
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
	var modelsets []models.Modelset
	var total int
	//获取当前用户创建的模型
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	modelsets, total, err = services.ListModelSets(req.PageNum, req.PageSize,req.IsAdvance, req.Name, req.Status, username)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := GetModelsetsResp{
		JsonModels: modelsets,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPage:  total/req.PageSize + 1,
	}
	return SuccessResp(c, data)
}

// @Summary get model by id
// @Produce  json
// @Param id path int true "model id"
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
// @Param body body createModelsetReq true "json body"
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
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	err = services.CreateModelset(req.IsAdvance,req.Name, req.Description, username, "0.0.1", req.Path, req.Use, req.JobId, req.DataFormat, req.Arguments, req.EngineType, req.Precision)
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
