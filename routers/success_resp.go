package routers

import (
	"net/http"

	"github.com/apulis/AIArtsBackend/models"

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

// edge inferences
type APISuccessRespGetFDInfo struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data GetFDInfoResp `json:"data"`
}

type APISuccessRespGetConvTypes struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data models.ConversionTypes `json:"data"`
}

type APISuccessRespLsEdgeInferences struct {
	Code int                  `json:"code"`
	Msg  string               `json:"msg"`
	Data LsEdgeInferencesResp `json:"data"`
}

type APISuccessRespCreateEdgeInference struct {
	Code int                     `json:"code"`
	Msg  string                  `json:"msg"`
	Data CreateEdgeInferenceResp `json:"data"`
}

// models
type APISuccessRespGetModelset struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data getModelsetResp `json:"data"`
}

type APISuccessRespGetModelsets struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetModelsetsResp `json:"data"`
}

// codes
type APISuccessRespAllGetCodeEnv struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data GetAllCodeEnvRsp `json:"data"`
}

type APISuccessRespCreateCodeEnv struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data CreateCodeEnvRsp `json:"data"`
}

type APISuccessRespDeleteCodeEnv struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data DeleteCodeEnvRsp `json:"data"`
}

type APISuccessRespGetCodeEnvJupyter struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data models.EndpointWrapper `json:"data"`
}

// training
type APISuccessRespGetAllTraining struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data GetAllTrainingRsp `json:"data"`
}

type APISuccessRespCreateTraining struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data CreateTrainingRsp `json:"data"`
}

type APISuccessRespDeleteTraining struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data DeleteTrainingRsp `json:"data"`
}

type APISuccessRespGetTraining struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data GetTrainingRsp `json:"data"`
}

type APISuccessRespGetTrainingLog struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data models.JobLog `json:"data"`
}

// template
type APISuccessRespGetAllTemplate struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data GetAllTemplateRsp `json:"data"`
}

type APISuccessRespCreateTemplate struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data CreateTemplateRsp `json:"data"`
}

type APISuccessRespDeleteTemplate struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data DeleteTemplateRsp `json:"data"`
}

type APISuccessRespGetTemplate struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data models.TemplateItem `json:"data"`
}

// common
type APISuccessRespGetResource struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data GetResourceRsp `json:"data"`
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
