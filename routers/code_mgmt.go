package routers

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupCode(r *gin.Engine) {
	group := r.Group("/ai_arts/api/codes")

	group.GET("/", wrapper(getAllCodeEnv))
	group.POST("/", wrapper(createCodeEnv))
	group.DELETE("/:id", wrapper(delCodeEnv))
	group.GET("/:id/jupyter", wrapper(getJupyterPath))
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

type DeleteCodeEnvRsp struct {
}

type CodeEnvId struct {
	Id string `json:"id"`
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
// @Param id query string true "codeEnv id"
// @Success 200 {object} APISuccessRespDeleteCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id [delete]
func delCodeEnv(c *gin.Context) error {

	id := c.Param("id")
	fmt.Println(id)

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err := services.DeleteCodeEnv(userName, id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	data := gin.H{}
	return SuccessResp(c, data)
}

// @Summary get CodeEnv jupyter path
// @Produce  json
// @Param id query string true "code id"
// @Success 200 {object} APISuccessRespGetCodeEnvJupyter "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id/jupyter [get]
func getJupyterPath(c *gin.Context) error {

	var err error
	var id string
	var rspData *models.EndpointWrapper

	id = c.Param("id")
	fmt.Println(id)

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(NO_USRNAME, "no username")
	}

	err, rspData = services.GetJupyterPath(userName, id)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, rspData)
}
