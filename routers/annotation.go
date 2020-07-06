package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddGroupAnnotation(r *gin.Engine) {
	group := r.Group("/api/annotations")

	group.GET("/", lsAllAnnotations)
}

// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/annotations [post]
func lsAllAnnotations(c *gin.Context) {
	res := APISuccessResp{
		Code: 0,
		Msg:  "success",
		Data: gin.H{},
	}
	c.JSON(http.StatusOK, res)
}
