package models

import (
	"strings"
)

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

func ValidHomePath(userName, path string) bool {
	path = strings.TrimSpace(path)
	usrHome := "/home/"+userName

	if !strings.HasPrefix(path, usrHome) {
		return false
	}

	return true
}

func (t *Training) ValidatePathByUser(userName string) (bool, string) {

	if !strings.HasSuffix(t.StartupFile, ".py") {
		return false, "启动文件非python"
	}

	if !ValidHomePath(userName, t.StartupFile)  {
		return false, "启动文件路径错误"
	}

	if !ValidHomePath(userName, t.CodePath) {
		return false, "代码路径不在home下"
	}

	if !ValidHomePath(userName, t.OutputPath) {
		return false, "输出文件路径不在home目录下"
	}

	return true, ""
} 

