package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupCode(r *gin.Engine) {
	group := r.Group("/ai_arts/api/codes")

	group.GET("/", wrapper(getAllCodeEnv))
	group.POST("/", wrapper(createCodeEnv))
	group.DELETE("/:id", wrapper(delCodeEnv))
}

type GetAllCodeEnvReq struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type GetAllCodeEnvRsp struct {
	CodeEnvs  []*models.CodeEnvItem `json:"CodeEnvs"`
	Total     int                   `json:"total"`
	totalPage int                   `json:"totalPage"`
}

type CreateCodeEnvReq struct {
	models.CreateCodeEnv
}

type CreateCodeEnvRsp struct {
	Id string `json:"id"`
}

type DeleteCodeEnvReq struct {
	Id string `json:"id"`
}

type DeleteCodeEnvRsp struct {
}

// @Summary get all codes
// @Produce  json
// @Param pageNum query int true "page number"
// @Param pageSize query int true "size per page"
// @Success 200 {object} APISuccessRespAllGetCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes [get]
func getAllCodeEnv(c *gin.Context) error {

	var req GetAllCodeEnvReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	sets, total, totalPage, err := services.GetAllCodeEnv(userName, req.PageNum, req.PageSize)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	rsp := &GetAllCodeEnvRsp{
		sets,
		total,
		totalPage,
	}

	return SuccessResp(c, rsp)
}

// @Summary create CodeEnv
// @Produce  json
// @Param param body CreateCodeEnvReq true "params"
// @Success 200 {object} APISuccessRespCreateCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes [post]
func createCodeEnv(c *gin.Context) error {

	var req CreateCodeEnvReq
	var id string

	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	id, err = services.CreateCodeEnv(userName, req.CreateCodeEnv)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, id)
}

// @Summary delete CodeEnv
// @Produce  json
// @Param id query string true "code set id"
// @Success 200 {object} APISuccessRespDeleteCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id [delete]
func delCodeEnv(c *gin.Context) error {

	var id DeleteCodeEnvReq
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err = services.DeleteCodeEnv(userName, id.Id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}
