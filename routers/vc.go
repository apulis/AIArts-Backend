package routers

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupVC(r *gin.Engine) {

	group := r.Group("/ai_arts/api/vc")
	group.Use(Auth())

	group.GET("/", wrapper(getVC))
	group.GET("/list", wrapper(listVC))
	group.POST("/", wrapper(addVC))
	group.PUT("/", wrapper(updateVC))
	group.DELETE("/", wrapper(delVC))
	group.GET("/count", wrapper(getVCStatistic))
}

// @Summary get vc
// @Produce  json
// @Param vcName query string true "name of vc"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc [get]
func getVC(c *gin.Context) error {

	var vcItem models.VCItem
	var err error

	if err = c.ShouldBindQuery(&vcItem); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err = services.OperateVC(userName, models.VC_OPTYPE_GET, &vcItem)
	if err != nil {
		return AppError(configs.VC_ERROR, err.Error())
	}

	return SuccessResp(c, vcItem)
}

// @Summary create vc
// @Produce  json
// @Param vcName query string true "name of vc"
// @Param quota  query string true "vc quota"
// @Param metadata query string true "metadata"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc [post]
func addVC(c *gin.Context) error {

	var vcItem models.VCItem
	var err error

	if err = c.ShouldBindJSON(&vcItem); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err = services.OperateVC(userName, models.VC_OPTYPE_ADD, &vcItem)
	if err != nil {
		return AppError(configs.VC_ERROR, err.Error())
	}

	return SuccessResp(c, nil)
}

// @Summary  delete vc
// @Produce  json
// @Param 	vcName query string true "name of vc"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc [delete]
func delVC(c *gin.Context) error {

	var vcItem models.VCItem
	var err error

	if err = c.ShouldBindJSON(&vcItem); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err = services.OperateVC(userName, models.VC_OPTYPE_DEL, &vcItem)
	if err != nil {
		return AppError(configs.VC_ERROR, err.Error())
	}

	return SuccessResp(c, nil)
}

// @Summary update vc
// @Produce  json
// @Param vcName 	query string true "name of vc"
// @Param quota  	query string true "quota"
// @Param metadata 	query string true "metadata"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc [put]
func updateVC(c *gin.Context) error {

	var vcItem models.VCItem
	var err error

	if err = c.ShouldBindJSON(&	vcItem); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err = services.OperateVC(userName, models.VC_OPTYPE_UPDATE, &vcItem)
	if err != nil {
		return AppError(configs.VC_ERROR, err.Error())
	}

	return SuccessResp(c, nil)
}

// @Summary get vc
// @Produce  json
// @Param userName query string true "username"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc/list [get]
func listVC(c *gin.Context) error {

	var req models.Paging
	var vcRsp *models.VCRsp
	var err error

	if err = c.ShouldBindQuery(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	vcRsp, err = services.ListVC(userName, req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, *vcRsp)
}

// @Summary get jobs of vc
// @Produce  json
// @Param userName query string true "username"
// @Success 200 {object} APISuccessRespGetAllTraining "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/vc/count [get]
func getVCStatistic(c *gin.Context) error {

	var req models.VCStatisticReq
	var vcRsp *models.VCStatisticRsp
	var err error

	if err = c.ShouldBindQuery(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	vcRsp, err = services.GetVCStatistic(userName, req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, vcRsp)
}