package routers

import (
	_ "encoding/json"
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupCode(r *gin.Engine) {

	group := r.Group("/ai_arts/api/codes")
	group.Use(Auth())

	group.GET("/", wrapper(getAllCodeEnv))
	group.POST("/", wrapper(createCodeEnv))
	group.DELETE("/:id", wrapper(delCodeEnv))
	group.GET("/:id/jupyter", wrapper(getJupyterPath))
	group.GET("/:id/endpoints", wrapper(getCodeEnvEndpoints))
	group.PATCH("/:id/endpoints", wrapper(addCodeEnvEndpoints))
	group.POST("/upload", wrapper(uploadCode))
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
// @Param status query string true "job status. get all jobs if it is all"
// @Param searchWord query string true "the keyword of search"
// @Success 200 {object} APISuccessRespAllGetCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes [get]
func getAllCodeEnv(c *gin.Context) error {

	var req models.GetAllJobsReq
	var err error

	if err = c.Bind(&req); err != nil {
		return ParameterError(err.Error())
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	// 兼容老代码
	if req.VCName == "" {
		req.VCName = "platform"
	}

	sets, total, totalPage, err := services.GetAllCodeEnv(userName, req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
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
		return AppError(configs.NO_USRNAME, "no username")
	}

	// 兼容老代码
	if req.VCName == "" {
		req.VCName = models.DefaultVcName
	}

	if req.JobTrainingType != models.TrainingTypeDist && req.JobTrainingType != models.TrainingTypeRegular {
		return AppError(configs.INVALID_TRAINING_TYPE, "任务类型非法")
	}

	imageName, err := services.ConvertImage(req.CreateCodeEnv.Engine, req.CreateCodeEnv.IsPrivateImage)
	if err != nil {
		return AppError(configs.DOCKER_IMAGE_NOT_FOUNT, "docker image not exist")
	}

	req.CreateCodeEnv.Engine = imageName
	id, err = services.CreateCodeEnv(c, userName, req.CreateCodeEnv)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
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
		return AppError(configs.NO_USRNAME, "no username")
	}

	err := services.DeleteCodeEnv(userName, id)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
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
		return AppError(configs.NO_USRNAME, "no username")
	}

	err, rspData = services.GetJupyterPath(userName, id)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, rspData)
}

// @Summary get CodeEnv endpoints
// @Produce  json
// @Param id query string true "code id"
// @Success 200 {object} APISuccessRespGetCodeEnvJupyter "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id/endpoints [post]
func getCodeEnvEndpoints(c *gin.Context) error {

	var err error
	var id string
	var rspData interface{}

	id = c.Param("id")
	fmt.Println(id)

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err, rspData = services.GetEndpoints(userName, id)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, rspData)
}

// @Summary add CodeEnv endpoints
// @Produce  json
// @Param id query string true "code id"
// @Success 200 {object} APISuccessRespGetCodeEnvJupyter "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/:id/endpoints [patch]
func addCodeEnvEndpoints(c *gin.Context) error {
	var err error

	var req models.AddEndportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return ParameterError(err.Error())
	}

	var id string
	id = c.Param("id")
	req.JobId = id

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	err, result := services.AddEndpoints(userName, id, req)
	if err != nil {
		return AppError(configs.APP_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, result)
}

// @Summary upload code
// @Produce  json
// @Param data body string true "upload file key 'data'"
// @Success 200 {object} APISuccessRespCreateCodeEnv "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/codes/upload [post]
func uploadCode(c *gin.Context) error {

	//多文件list
	logger.Info("starting upload file")

	file, err := c.FormFile("file")
	if err != nil {
		return ParameterError(err.Error())
	}

	codePath := c.PostForm("codePath")
	if len(codePath) == 0 {
		return ParameterError("缺少代码路径")
	}

	userName := getUsername(c)
	if len(userName) == 0 {
		return AppError(configs.NO_USRNAME, "no username")
	}

	//outputDir, err := services.ConvertPath(userName, codePath)
	//if err != nil {
	//	return AppError(INVALID_CODE_PATH, err.Error())
	//}

	logger.Info("starting saving file")
	err = c.SaveUploadedFile(file, codePath+"/"+file.Filename)
	if err != nil {
		return AppError(configs.SAVE_FILE_ERROR_CODE, err.Error())
	}

	logger.Info("starting change file mode")

	filePath := fmt.Sprintf("%s/%s", codePath, file.Filename)
	if err = services.UploadDone(userName, filePath); err != nil {
		return AppError(configs.COMPLETE_UPLOAD_ERR, err.Error())
	}

	return SuccessResp(c, nil)
}
