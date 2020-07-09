package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupCode(r *gin.Engine) {
	group := r.Group("/ai_arts/api/codes")

	group.GET("/", wrapper(getCodeset))
	group.POST("/:id", wrapper(createCodeset))
	group.DELETE("/:id", wrapper(delCodeset))
}


type GetCodesetReq struct {
	PageNum 	int 	`json:"pageNum"`
	PageSize 	int 	`json:"pageSize"`
}

type GetCodesetRsp struct {
	Codesets 	[] *models.CodesetItem `json:"codesets"`
	Total		int   	`json:"total"`
	totalPage	int 	`json:"totalPage"`
}


type CreateCodesetReq struct {
	Name 			string `json:"name"`
	Description 	string `json:"description" binding:"required"`
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
// @Router /ai_arts/api/code [get]
func getCodeset(c *gin.Context) error {

	var req GetCodesetReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	sets, total, totalPage, err := services.GetCodeset(req.PageNum, req.PageSize)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetCodesetRsp {
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary create codeset
// @Produce  json
// @Param name query string true "dataset name"
// @Param description query string true "dataset description"
// @Param creator query string true "dataset creator"
// @Param path query string true "dataset storage path"
// @Success 200 {object} APISuccessRespCreateCodeset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets [post]
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
// @Param description query string true "dataset description"
// @Success 200 {object} APISuccessRespDeleteCodeset "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets/:id [post]
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
