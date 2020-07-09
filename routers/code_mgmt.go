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
	Engine          string `json:"engine"`
	DeviceType		string `json:"deviceType"`
	DeviceNum 		int `json:"deviceNum"`
	Desc 			string `json:"desc"`
	CodePath		string `json:"codePath"`
}

type CreateCodesetRsp struct {
	Id 				string `json:"id"`
}

type DeleteCodesetReq struct {
	Id 				string `json:"id"`
}

type DeleteCodesetRsp struct {

}

// @Summary get all codes
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Success 200 {object} APISuccessRespAllGetCodeset "success"
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
// @Param name query string true "code name"
// @Param engine query string true "engine"
// @Param deviceType query string true "device type"
// @Param deviceNum query string true "device number"
// @Param codePath query string true "code path"
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
// @Param id query string true "code set id"
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
