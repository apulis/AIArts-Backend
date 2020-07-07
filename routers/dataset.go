package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupDataset(r *gin.Engine) {
	group := r.Group("/ai_arts/api/datasets")

	group.GET("/", wrapper(lsDatasets))
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
// @Router /ai_arts/api/datasets [get]
func lsDatasets(c *gin.Context) error {
	datasets := services.ListDatasets()
	data := gin.H{"datasets": datasets}
	return SuccessResp(c, data)
}

// @Summary create dataset
// @Produce  json
// @Param name query string true "dataset name"
// @Param description query string true "dataset description"
// @Param creator query string true "dataset creator"
// @Param path query string true "dataset storage path"
// @Success 200 {string} json "{"code":0,"data":{},"msg":"success"}"
// @Router /ai_arts/api/datasets [post]
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
