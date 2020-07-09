package routers

import (
	"os"

	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupFile(r *gin.Engine) {
	group := r.Group("/ai_arts/api/files")

	group.POST("/upload/dataset", wrapper(uploadDataset))
	group.POST("/upload/model", wrapper(uploadModel))
}

// @Summary upload dataset file
// @Produce  json
// @Param data body string true "upload file key 'data'"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/dataset [get]
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

	return SuccessResp(c, gin.H{"path": unzippedPath})
}

// @Summary upload dataset file
// @Produce  json
// @Param data body string true "upload file key 'data'"
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/files/upload/model [get]
func uploadModel(c *gin.Context) error {
	return nil
}
