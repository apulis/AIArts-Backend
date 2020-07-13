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

	group.POST("/upload/dataset", wrapper(uploadDataset))
	group.POST("/upload/model", wrapper(uploadModelset))
	group.GET("/download/model/:id", wrapper(downloadModelset))
}

type UploadFileResp struct {
	Path string `json:"path"`
}

// @Summary upload dataset file
// @Produce  json
// @Param data body string true "upload file key 'data'"
// @Success 200 {object} UploadFileResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/dataset [post]
func uploadDataset(c *gin.Context) error {
	file, err := c.FormFile("data")
	if err != nil {
		return ParameterError(err.Error())
	}

	if services.CheckFileOversize(file.Size) {
		return AppError(FILE_OVERSIZE_CODE, "File over size limit")
	}

	filetype, err := services.CheckFileName(file.Filename)
	if err != nil {
		return AppError(FILETYPE_NOT_SUPPORTED_CODE, err.Error())
	}

	filePath, err := services.GetDatasetTempPath(filetype)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}

	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return AppError(SAVE_FILE_ERROR_CODE, err.Error())
	}

	unzippedPath, err := services.ExtractFile(filePath, filetype)
	if err != nil {
		return AppError(EXTRACT_FILE_ERROR_CODE, err.Error())
	}

	err = os.Remove(filePath)
	if err != nil {
		return AppError(REMOVE_FILE_ERROR_CODE, err.Error())
	}

	return SuccessResp(c, UploadFileResp{Path: unzippedPath})
}

// @Summary upload model file, not implemented yet
// @Produce  json
// @Param data body string true "upload file key 'data'"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/model [post]
func uploadModelset(c *gin.Context) error {
	return nil
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
	err = services.CheckPathExists(modelset.Path)
	if err != nil {
		return AppError(FILEPATH_NOT_EXISTS_CODE, err.Error())
	}
	targetPath, err := services.CompressFile(modelset.Path)
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