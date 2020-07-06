package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APISuccessResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data gin.H  `json:"data"`
}

func SuccessResp(c *gin.Context, data gin.H) error {
	res := APISuccessResp{
		Code: SUCCESS_CODE,
		Msg:  "success",
		Data: data,
	}
	c.JSON(http.StatusOK, res)
	return nil
}
