package models


type Training struct {
	Id              string `json:"id"`
	Name 			string `json:"name"`
	Engine          string `json:"engine"`
	DeviceType		string `json:"deviceType"`
	DeviceNum 		int `json:"deviceNum"`
	CodePath		string `json:"codePath"`
	StartupFile		string `json:"startupFile"`
	OutputPath		string `json:"outputPath"`
	DatasetPath		string `json:"datasetPath"`
	Params 			map[string]string `json:"params"`
	Desc 			string `json:"desc"`
	Status 			string `json:"status"`
	CreateTime		string `json:"createTime"`
}

