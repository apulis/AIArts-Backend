package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupCode(r *gin.Engine) {
	group := r.Group("/ai_arts/api/codes")

	group.GET("/", wrapper(getAllCodeset))
	group.POST("/", wrapper(createCodeset))
	group.DELETE("/:id", wrapper(delCodeset))
}


type GetAllCodesetReq struct {
	PageNum 	int 	`json:"pageNum"`
	PageSize 	int 	`json:"pageSize"`
}

type GetAllCodesetRsp struct {
	Codesets 	[] *models.CodesetItem `json:"codesets"`
	Total		int   	`json:"total"`
	totalPage	int 	`json:"totalPage"`
}

type CreateCodesetReq struct {
	Name 			string `json:"name"`
	Description 	string `json:"description" binding:"required"`
	CodePath		string `json:"codePath"`
	FrameworkInfo 	models.AIFrameworkItem `json:"frameworkInfo"`
}

type CreateCodesetRsp struct {
	Id 				string `json:"id"`
}

type DeleteCodesetReq struct {
	Id 				string `json:"id"`
}

type DeleteCodesetRsp struct {

}

// @Summary list codesets
// @Produce  json
// @Param page query int true "page number"
// @Param pagesize query int true "size per page"
// @Success 200 {object} APISuccessRespGetCodeset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes [get]
func getAllCodeset(c *gin.Context) error {

	var req GetAllCodesetReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	sets, total, totalPage, err := services.GetAllCodeset(req.PageNum, req.PageSize)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetAllCodesetRsp{
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary create codeset
// @Produce  json
// @Param name query string true "codeset name"
// @Param description query string true "codeset description"
// @Param creator query string true "codeset creator"
// @Param path query string true "codeset storage path"
// @Success 200 {object} APISuccessRespCreateCodeset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes [post]
func createCodeset(c *gin.Context) error {

	var req CreateCodesetReq
        var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	id, err = services.CreateCodeset(req.Name, req.Description, req.FrameworkInfo)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary delete codeset
// @Produce  json
// @Param description query string true "codeset description"
// @Success 200 {object} APISuccessRespDeleteCodeset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id [delete]
func delCodeset(c *gin.Context) error {

	var id string
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	var req DeleteCodesetReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	err = services.DeleteCodeset(id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}
