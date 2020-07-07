package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupDataset(r *gin.Engine) {
	group := r.Group("/ai_arts/api/datasets")

	group.GET("/", wrapper(lsAllDatasets))
	group.POST("/", wrapper(createDataset))
}

type createDatasetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Creator     string `json:"creator" binding:"required"`
	Path        string `json:"path" binding:"required"`
}

// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/datasets [post]
func lsAllDatasets(c *gin.Context) error {
	datasets := services.ListDatasets()
	data := gin.H{"datasets": datasets}
	return SuccessResp(c, data)
}

func createDataset(c *gin.Context) error {
	var req createDatasetReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return ParameterError(err.Error())
	}
	err = services.CreateDataset(req.Name, req.Description, req.Creator, "0.0.1", req.Path)
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}
	data := gin.H{}
	return SuccessResp(c, data)
}
