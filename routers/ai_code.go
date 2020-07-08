package routers

import (
	"github.com/gin-gonic/gin"
)


// @Summary sample
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/codes [post]
func lsAllCodes(c *gin.Context) error {
	data := gin.H{}
	return SuccessResp(c, data)
}
