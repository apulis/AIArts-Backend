package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupFile(r *gin.Engine) {
	group := r.Group("/ai_arts/api/files")

	group.POST("/upload/dataset", wrapper(uploadDataset))
}

func uploadDataset(c *gin.Context) error {
	file, err := c.FormFile("data")
	if err != nil {
		return ParameterError(err.Error())
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
	return SuccessResp(c, gin.H{"path": unzippedPath})
}
