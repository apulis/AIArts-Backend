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
	Role string
	dataSetBindId string
	dataSetPath string
}
