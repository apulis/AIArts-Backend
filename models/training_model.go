package models

import (
	"strings"
)

type Training struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Engine          string            `json:"engine"`
	IsPrivateImg    bool              `json:"private"`
	DeviceType      string            `json:"deviceType"`
	DeviceNum       int               `json:"deviceNum"`
	CodePath        string            `json:"codePath"`
	StartupFile     string            `json:"startupFile,omitempty"`
	OutputPath      string            `json:"outputPath"`
	DatasetPath     string            `json:"datasetPath"`
	Params          map[string]string `json:"params"`
	Desc            string            `json:"desc"`
	Status          string            `json:"status"`
	VisualPath      string            `json:"visualPath,omitempty"`
	CreateTime      string            `json:"createTime"`
	JobTrainingType string            `json:"jobTrainingType"`
	NumPs           int               `json:"numPs"`
	NumPsWorker     int               `json:"numPsWorker"`
	VCName          string            `json:"vcName"`
	Command         string            `json:"command"`

	FrameworkType   string            `json:"frameworkType"`
}

func ValidHomePath(userName, path string) bool {
	path = strings.TrimSpace(path)
	usrHome := "/home/" + userName

	if !strings.HasPrefix(path, usrHome) {
		return false
	}

	return true
}

func (t *Training) ValidatePathByUser(userName string) (bool, string) {

	//if !strings.HasSuffix(t.StartupFile, ".py") {
	//	return false, "启动文件非python"
	//}

	//if !ValidHomePath(userName, t.StartupFile) {
	//	return false, "启动文件路径错误"
	//}

	//if !ValidHomePath(userName, t.CodePath) {
	//	return false, "代码路径不在home下"
	//}

	//if !ValidHomePath(userName, t.OutputPath) {
	//	return false, "输出文件路径不在home目录下"
	//}

	return true, ""
}

type GetAllJobsReq struct {
	PageNum    int    `form:"pageNum" json:"pageNum"`
	PageSize   int    `form:"pageSize" json:"pageSize"`
	JobStatus  string `form:"status" json:"status"`
	JobType    string `form:"jobType" json:"jobType"`
	SearchWord string `form:"searchWord" json:"searchWord"`
	OrderBy    string `form:"orderBy" json:"orderBy"`
	Order      string `form:"order" json:"order"`
	VCName     string `form:"vcName" json:"vcName"`
}
