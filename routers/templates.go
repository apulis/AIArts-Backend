package routers

import (
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
	group.DELETE("/:id", wrapper(delTemplate))
}

type GetAllTemplateReq struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	From     string `json:"from"`
}

type GetAllTemplateRsp struct {
	Templates []*models.Template `json:"Templates"`
	Total     int                `json:"total"`
	totalPage int                `json:"totalPage"`
}

type CreateTemplateReq struct {
	models.Template
}

type CreateTemplateRsp struct {
	Id int `json:"id"`
}

type DeleteTemplateReq struct {
	Id string `json:"id"`
}

type DeleteTemplateRsp struct {
}

type GetTemplateReq struct {
	Id string `json:"id"`
}

type GetTemplateRsp struct {
	models.Template
}

// @Summary get all Templates
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Param status query string true "job status. get all jobs if it is all"
// @Param searchWord query string true "the keyword of search"
// @Success 200 {object} APISuccessRespGetAllTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates [get]
func getAllTemplates(c *gin.Context) error {

	var req GetAllJobsReq
	var err error

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	sets, total, totalPage, err := services.GetAllTemplate(userName, req.PageNum, req.PageSize, req.JobStatus, req.SearchWord)
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

// @Summary create Template
// @Produce json
// @Param param body CreateTemplateReq true "params"
// @Success 200 {object} APISuccessRespCreateTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates [post]
func createTemplate(c *gin.Context) error {

	var req models.Template
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	id, err = services.CreateTemplate(userName, req)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary get specific Template
// @Produce  json
// @Param param body GetTemplateReq true "params"
// @Success 200 {object} APISuccessRespGetTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates/:id [get]
func getTemplate(c *gin.Context) error {

	var id models.UriJobId
	var Template *models.Template

	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	Template, err = services.GetTemplate(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, Template)
}

// @Summary delete one Template
// @Produce  json
// @Param param body DeleteTemplateReq true "params"
// @Success 200 {object} APISuccessRespDeleteTemplate "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/Templates/:id [delete]
func delTemplate(c *gin.Context) error {

	var id models.UriJobId
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
