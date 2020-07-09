package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APISuccessResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type APISuccessRespGetDataset struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data GetDatasetResp `json:"data"`
}

type APISuccessRespGetDatasets struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data GetDatasetsResp `json:"data"`
}

type APISuccessRespGetModelset struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data GetModelsetResp `json:"data"`
}

type APISuccessRespGetModelsets struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetModelsetsResp `json:"data"`
}

func SuccessResp(c *gin.Context, data interface{}) error {
	res := APISuccessResp{
		Code: SUCCESS_CODE,
		Msg:  "success",
		Data: data,
	}
	c.JSON(http.StatusOK, res)
	return nil
}
