package routers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupFile(r *gin.Engine) {
	group := r.Group("/ai_arts/api/files")
	group.GET("/download/model/:id", wrapper(downloadModelset))
	group.GET("/download/dataset/:id", wrapper(downloadDataset))
	group.Use(Auth())
	group.POST("/upload/dataset", wrapper(uploadDataset))
	group.POST("/upload/model", wrapper(uploadModelset))
}

type UploadFileResp struct {
	Path string `json:"path"`
}

// @Summary upload dataset file
// @Produce  json
// @Param data query string true "upload file key 'data'"
// @Param isPrivate query string true "isPrivate key 'isPrivate'"
// @Param dir query string true "upload file directory 'dir'"
// @Success 200 {object} UploadFileResp "success"
// @Failure 400 {object} APIException "error code:30009,msg:the /tmp direct is full"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/dataset [post]

func uploadDataset(c *gin.Context) error {

	logger.Info("starting upload file")

	file, err := c.FormFile("data")
	isPrivate := c.PostForm("isPrivate")

	//存储文件夹
	dir := c.PostForm("dir")
	if err != nil {
		return AppError(UPLOAD_TEMPDIR_FULL_CODE, err.Error())
	}

	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	} //取消大小限制

	// 获取文件类别
	filetype, err := services.CheckFileName(file.Filename)
		if err != nil {
		return AppError(FILETYPE_NOT_SUPPORTED_CODE, err.Error())
	}

	// 获取数据集真实路径，但文件夹名称依dir结尾目录
	// 待用户创建数据集时，会重命名此目录为真实名称
	datasetStoragePath := services.GenerateDatasetStoragePath(username, dir, isPrivate)
	logger.Info("uploadDataset - datasetStoragePath", datasetStoragePath)

	// 获取文件临时目录
	filePath, err := services.GetDatasetTempPath(dir)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}

	logger.Info("uploadDataset - filePath", filePath)
	logger.Info("starting saving file")

	// 将文件保存为filePath（包括文件名）
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}

	logger.Info("starting extract file")

	// 检查数据集文件是否已存在
	err = services.CheckPathExists(datasetStoragePath)
	if err == nil {
		return AppError(DATASET_IS_EXISTED, "same dataset found! cannot move to target path")
	}

	// 解压并返回解压后的目录
	unzippedPath, err := services.ExtractFile(filePath, filetype, datasetStoragePath)
	if err != nil {
		return AppError(EXTRACT_FILE_ERROR_CODE, err.Error())
	}

	// 删除临时文件
	logger.Info("starting remove file")
	err = os.Remove(filePath)
	if err != nil {
		return AppError(REMOVE_FILE_ERROR_CODE, err.Error())
	}

	// 不移除，等到创建数据集时
	if isPrivate=="false"{
		_ = os.Chmod(unzippedPath, os.ModePerm)
	}

	// 返回解压后的目录，前端调用数据集创建接口会传入此参数
	// 后端根据此参数将文件夹重命名
	logger.Info("unzippedPath - unzippedPath", unzippedPath)
	return SuccessResp(c, UploadFileResp{Path: unzippedPath})
}

// @Summary download dataset by id
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/download/dataset/:id [get]
func downloadDataset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	dataset, err := services.GetDataset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	err = services.CheckPathExists(dataset.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	targetPath, err := services.CompressFile(dataset.Path)
	if err != nil {
		return AppError(COMPRESS_PATH_ERROR_CODE, err.Error())
	}
	fi, _ := os.Stat(targetPath)

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprint("attachment; filename=", fi.Name()))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(targetPath)

	return nil
}

// @Summary upload model file, not implemented yet
// @Produce  json
// @Param data query string true "upload file key 'data'"
// @Param dir query string true "upload file directory 'dir'"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/model [post]
func uploadModelset(c *gin.Context) error {
	logger.Info("starting upload model")
	file, err := c.FormFile("data")
	dir := c.PostForm("dir")
	username := getUsername(c)
	if len(username) == 0 {
		return AppError(NO_USRNAME, "no username")
	}
	//存储文件夹
	if err != nil {
		return AppError(UPLOAD_TEMPDIR_FULL_CODE, err.Error())
	}
	filetype, err := services.CheckFileName(file.Filename)
	if err != nil {
		return AppError(FILETYPE_NOT_SUPPORTED_CODE, err.Error())
	}
	filePath, err := services.GetModelTempPath(filetype)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}
	logger.Info("starting saving file")
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}
	logger.Info("starting extract file")
	datasetStoragePath := services.GenerateModelStoragePath(dir, username)
	unzippedPath, err := services.ExtractFile(filePath, filetype, datasetStoragePath)
	if err != nil {
		return AppError(EXTRACT_FILE_ERROR_CODE, err.Error())
	}
	logger.Info("starting remove file")
	err = os.Remove(filePath)
	if err != nil {
		return AppError(REMOVE_FILE_ERROR_CODE, err.Error())
	}
	return SuccessResp(c, UploadFileResp{Path: unzippedPath})
}

// @Summary download model by id
// @Produce  json
// @Param id path int true "model id"
// @Success 200 {object} APISuccessResp "success  download the code path dir"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/download/model/:id [get]
func downloadModelset(c *gin.Context) error {
	var id modelsetId
	err := c.ShouldBindUri(&id)
	if err != nil {
		return ParameterError(err.Error())
	}
	modelset, err := services.GetModelset(id.ID)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	//just download the code path
	targetPath, err := services.CompressFile(modelset.CodePath)
	if err != nil {
		return AppError(COMPRESS_PATH_ERROR_CODE, err.Error())
	}
	fi, _ := os.Stat(targetPath)

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", fmt.Sprint("attachment; filename=", fi.Name()))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(targetPath)

	return nil
}
