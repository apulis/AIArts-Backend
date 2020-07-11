package models

type Project struct {
	ProjectId string
	Name string
	Info string
	Role string
}

type DataSet struct {
	dataSetId string
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
}

type DatasetsReq struct {
	Successful string
	Msg string
	Datasets []DataSet
}