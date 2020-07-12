package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APISuccessResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// dataset
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

// models
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

// codes
type APISuccessRespAllGetCodeset struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetAllCodesetRsp	  `json:"data"`
}

type APISuccessRespCreateCodeset struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data CreateCodesetRsp `json:"data"`
}

type APISuccessRespDeleteCodeset struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data DeleteCodesetRsp `json:"data"`
}

type APISuccessRespGetResource struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetResourceRsp `json:"data"`
}

type APISuccessRespGetAllTraining struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetAllTrainingRsp	  `json:"data"`
}

type APISuccessRespCreateTraining struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data CreateTrainingRsp `json:"data"`
}

type APISuccessRespDeleteTraining struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data DeleteTrainingRsp `json:"data"`
}

type APISuccessRespGetTraining struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetTrainingRsp `json:"data"`
}

type APISuccessRespGetTrainingLog struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data models.JobLog	  `json:"data"`
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
