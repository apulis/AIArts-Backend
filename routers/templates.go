package routers

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupTemplate(r *gin.Engine) {
	group := r.Group("/ai_arts/api/templates")

	group.Use(Auth())

	group.GET("/", wrapper(getAllTemplates))
	group.GET("/:id", wrapper(getTemplate))
	group.POST("/", wrapper(createTemplate))
	group.PUT("/:id", wrapper(updateTemplate))
	group.DELETE("/:id", wrapper(delTemplate))
}

type GetAllTemplateReq struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Scope    int    `json:"scope"`
	JobType  string `json:"jobType"`
}

type GetAllTemplateRsp struct {
	Templates []*models.TemplateItem `json:"Templates"`
	Total     int                    `json:"total"`
	totalPage int                    `json:"totalPage"`
}

type CreateTemplateReq struct {
	Scope        int                   `json:"scope"`
	JobType      string                `json:"jobType"`
	TemplateData models.TemplateParams `json:"templateData"`
}

type CreateTemplateRsp struct {
	Id int `json:"id"`
}

type DeleteTemplateReq struct {
	Id int64 `json:"id"`
}

type DeleteTemplateRsp struct {
}

type UpdateTemplateReq struct {
	Id           int64                 `json:"id"`
	Scope        int                   `json:"scope"`
	JobType      string                `json:"jobType"`
	TemplateData models.TemplateParams `json:"templateData"`
}

type GetTemplateReq struct {
	Id int64 `json:"id"`
}

// @Summary get all templates
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Param jobType query string true "training module: artsTraining, code module: codeEnv"
// @Param scope query string true "public: 1, private: 2"
// @Success 200 {object} APISuccessRespGetAllTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates [get]
func getAllTemplates(c *gin.Context) error {

	var req GetAllTemplateReq
	var err error

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	sets, total, totalPage, err := services.GetAllTemplate(userName, req.PageNum, req.PageSize, req.Scope, req.JobType)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetAllTemplateRsp{
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary create template
// @Produce json
// @Param param body CreateTemplateReq true "params"
// @Success 200 {object} APISuccessRespCreateTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates [post]
func createTemplate(c *gin.Context) error {

	var req CreateTemplateReq
	var id int64

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	id, err = services.CreateTemplate(userName, req.Scope, req.JobType, req.TemplateData)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary update template
// @Produce json
// @Param param body CreateTemplateReq true "params"
// @Success 200 {object} APISuccessRespCreateTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates [post]
func updateTemplate(c *gin.Context) error {

	var req UpdateTemplateReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err = services.UpdateTemplate(req.Id, userName, req.Scope, req.JobType, req.TemplateData)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, nil)
}

// @Summary get specific template
// @Produce  json
// @Param param body GetTemplateReq true "params"
// @Success 200 {object} APISuccessRespGetTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates/:id [get]
func getTemplate(c *gin.Context) error {

	var id GetTemplateReq
	var dbRecord *models.Templates

	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	dbRecord, err = services.GetTemplate(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rspData := &models.TemplateItem{
		MetaData: models.TemplateMeta{
			Name:      dbRecord.Name,
			Scope:     dbRecord.Scope,
			JobType:   dbRecord.JobType,
			Creator:   dbRecord.Creator,
			CreatedAt: dbRecord.CreatedAt,
			UpdatedAt: dbRecord.UpdatedAt,
		},
		Params: models.TemplateParams{},
	}

	err = json.Unmarshal([]byte(dbRecord.Data), rspData)
	if err != nil {
		return AppError(TEMPLATE_INVALID_PARAMS, err.Error())
	}

	return SuccessResp(c, rspData)
}

// @Summary delete one template
// @Produce  json
// @Param param body DeleteTemplateReq true "params"
// @Success 200 {object} APISuccessRespDeleteTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates/:id [delete]
func delTemplate(c *gin.Context) error {

	var id DeleteTemplateReq
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err = services.DeleteTemplate(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}
