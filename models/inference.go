package models

import "strings"

type PostInference struct {
	UserName        string `json:"userName"`
	Image           string `json:"image"`
	UserId          string `json:"userId"`
	JobName         string `json:"jobName"`
	Model_base_path string `json:"model_base_path"`
	VcName          string `json:"vcName"`
	GpuType         string `json:"gpuType"`
	Framework       string `json:"framework"`
	Device          string `json:"device"`
	Resourcegpu     int    `json:"resourcegpu"`
	DESC            string `json:"desc"`
	VERSION         string `json:"version"`
}

type InferenceJobResp struct {
	JobId string `json:"jobId"`
}

type QueryStringParametersV2 struct {
	PageNum  int    `form:"pageNum"`
	PageSize int    `form:"pageSize"`
	Name     string `form:"name"`
	Status   string `form:"status"`
	OrderBy  string `form:"orderBy"`
	Order    string `form:"order"`
}

func (queryStringParameters QueryStringParametersV2) GetPageNum() int {
	if queryStringParameters.PageNum <= 0 {
		return 1
	}
	return queryStringParameters.PageNum
}

func (queryStringParameters QueryStringParametersV2) GetPageSize() int {
	if queryStringParameters.PageSize < 0 {
		return 5
	}
	if queryStringParameters.PageSize >= 100 {
		return 100
	}
	return queryStringParameters.PageSize
}

func (queryStringParameters QueryStringParametersV2) GetName() string {
	return strings.TrimSpace(queryStringParameters.Name)
}
