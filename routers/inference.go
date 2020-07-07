package routers

import (
	"github.com/gin-gonic/gin"
)

func AddGroupInference(r *gin.Engine) {
	group := r.Group("/ai_arts/api/inferences")

	group.GET("/", wrapper(lsAllInferences))
}

// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/inferences [post]
func lsAllInferences(c *gin.Context) error {
	data := gin.H{}
	return SuccessResp(c, data)
}
