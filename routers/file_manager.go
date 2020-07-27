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
	//if services.CheckFileOversize(file.Size) {
	//	return AppError(FILE_OVERSIZE_CODE, "File over size limit")
	//}
	filetype, err := services.CheckFileName(file.Filename)
	if err != nil {
		return AppError(FILETYPE_NOT_SUPPORTED_CODE, err.Error())
	}
	filePath, err := services.GetDatasetTempPath(filetype)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}
	logger.Info("starting saving file")
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}
	logger.Info("starting extract file")

	datasetStoragePath := services.GenerateDatasetStoragePath(dir, isPrivate, username)
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
// @Success 200 {object} APISuccessResp "success"
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
	//如果上传模型文件检查模型文件是否存在
	if modelset.ModelPath != "" {
		err = services.CheckPathExists(modelset.ModelPath)
		if err != nil {
			return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
		}
	}
	//检查模型参数文件是否存在
	err = services.CheckPathExists(modelset.ArgumentPath)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	targetPath, err := services.CompressFile(modelset.ModelPath)
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
