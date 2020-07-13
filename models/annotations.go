package models

type Project struct {
	ProjectId string
	Name string
	Info string
	Role string
}

type DataSet struct {
	DataSetId string
	Name string
	Info string
	Type string
	DataSetBindId string
	DataSetPath string
}

type ProjectsReq struct {
	Successful string
	Msg string
	Projects []Project
	TotalCount int `json:"totalCount"`
}

type DatasetsReq struct {
	Successful string
	Msg string
	Datasets []DataSet
	TotalCount int `json:"totalCount"`
}

type DatasetReq struct {
	Successful string
	Msg string
	Info DataSet
}

type UpdateDataSet struct {
	Name string
	Info string
	Type string
	DataSetBindId string
	DataSetPath string
}

type TasksList struct {
	Successful string
	Msg string
	TaskList []interface{}
	TotalCount int `json:"totalCount"`
}

type NextTask struct {
	Successful string
	Msg string
	Next interface{}
}

type OneTask struct {
	Successful string
	Msg string
	Annotations interface{}
}

type QueryStringParameters struct {
	Page int `schema:"page"`
	Size int `schema:"size"`
}