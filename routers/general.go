package routers

import (
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupDataset(r *gin.Engine) {
	group := r.Group("/ai_arts/api/common")

	group.GET("/resource", wrapper(getResource))
}

// @Summary get available resource
// @Produce  json
// @Success 200 {object} APISuccessResp "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/common/resource [get]
func getResourceReq(c *gin.Context) error {

	rsp, err := services.GetResourceReq()
	if err != nil {
		return AppError(APP_ERROR_CODE, err.Error())
	}

	return SuccessRsp(c, rsp)
}