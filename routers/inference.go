package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGroupInference(r *gin.Engine) {
	group := r.Group("/api/inferences")

	group.GET("/", lsAllInferences)
}

// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/inferences [post]
func lsAllInferences(c *gin.Context) {
	res := Resp{
		Code: 0,
		Msg:  "success",
		Data: gin.H{},
	}
	c.JSON(http.StatusOK, res)
}