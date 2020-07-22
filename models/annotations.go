package models

import (
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/gin-gonic/gin"
	"strings"
)

type Project struct {
	ProjectId string	`json:"projectId"`
	Name string	`json:"name"`
	Info string	`json:"info"`
	Role string	`json:"role"`
}

type DataSet struct {
	DataSetId string	`json:"dataSetId"`
	Name string	`json:"name"`
	Info string	`json:"info"`
	Type string	`json:"type"`
	Labels []Label `json:"labels"`
	DataSetBindId int `json:"dataSetBindId"`
	DataSetPath string `json:"dataSetPath"`
	ConvertStatus string `json:"convertStatus"`
	ConvertOutPath string `json:"convertOutPath"`
}

type ProjectsReq struct {
	Successful string
	Msg string
	Projects []Project `json:"projects"`
	TotalCount int `json:"totalCount"`
}

type DatasetsReq struct {
	Successful string
	Msg string
	Datasets []DataSet `json:"datasets"`
	TotalCount int `json:"totalCount"`
}

type DatasetReq struct {
	Successful string
	Msg string
	Info DataSet `json:"info"`
}

type UpdateDataSet struct {
	Name string	`json:"name"`
	Info string `json:"info"`
	Type string	`json:"type"`
	DataSetBindId int `json:"dataSetBindId"`
	DataSetPath string `json:"dataSetPath"`
	Labels []Label `json:"labels"`
}

type Label struct {
	Id int				`json:"id"`
	Name string			`json:"name"`
	Type string 		`json:"type"`
	Supercategory string	`json:"supercategory"`
}

type LabelReq struct {
	Successful string
	Msg string
	Annotations interface{}	`json:"annotations"`
}

type TasksList struct {
	Successful string
	Msg string
	TaskList []interface{} `json:"taskList"`
	TotalCount int `json:"totalCount"`
}

type NextTask struct {
	Successful string
	Msg string
	Next interface{} `json:"next"`
}

type PreviousTask struct {
	Successful string
	Msg string
	Previous interface{} `json:"previous"`
}

type OneTask struct {
	Successful string
	Msg string
	Annotations interface{} `json:"annotations"`
}

type QueryStringParameters struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (queryStringParameters QueryStringParameters) GetPageNum() int {
	if queryStringParameters.Page<=0 {
		return 1
	}
	return queryStringParameters.Page
}

func (queryStringParameters QueryStringParameters) GetPageSize() int {
	if queryStringParameters.Size<=0 {
		return 5
	}
	if queryStringParameters.Size>=100 {
		return 100
	}
	return queryStringParameters.Size
}

type QueryStringParamInterface interface {
	GetPageNum () int
	GetPageSize () int
}

type GinContext struct {
	Context *gin.Context
}
func (ct GinContext) SaveToken()  {
	token := ct.Context.GetHeader("Authorization")
	token = strings.Split(token, "Bearer ")[1]
	configs.Config.Token = token
}

type ConvertDataFormat struct {
	ProjectId string `json:"projectId"`
	DatasetId string `json:"datasetId"`
	DatasetType string `json:"type"`
	ConvertTarget string `json:"target"`
}