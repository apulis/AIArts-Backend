package models

type PostInference struct {
	UserName string	`json:"userName"`
	Image string	`json:"image"`
	UserId string	`json:"userId"`
	JobName string	`json:"jobName"`
	Model_base_path string	`json:"model_base_path"`
	VcName string	`json:"vcName"`
	GpuType string	`json:"gpuType"`
	Framework string	`json:"framework"`
	Device string	`json:"device"`
	Resourcegpu int	`json:"resourcegpu"`
}

type InferenceJobResp struct {
	JobId string `json:"jobId"`
}

type QueryStringParametersV2 struct {
	PageNum int `form:"pageNum"`
	PageSize int `form:"pageSize"`
}