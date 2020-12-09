package models

type CodeEnvItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"   binding:"required"`
	Status     string `json:"status" binding:"required"`
	Engine     string `json:"engine"  binding:"required"`
	CodePath   string `json:"codePath"`
	Cmd        string `json:"cmd"`
	JupyterUrl string `json:"JupyterUrl"`
	CreateTime string `json:"createTime"`
	DeviceType string `json:"deviceType"`
	DeviceNum  int    `json:"deviceNum"`
	Desc       string `json:"desc"`
}

type CreateCodeEnv struct {
	Name       string `json:"name"   binding:"required"`
	Engine     string `json:"engine"  binding:"required"`
	CodePath   string `json:"codePath"`
	DeviceType string `json:"deviceType"`
	DeviceNum  int    `json:"deviceNum"`
	Desc       string `json:"desc"`
	VCName     string `json:"vcName"`
	Cmd        string `json:"cmd"`

	IsPrivateImage  bool   `json:"private"`
	JobTrainingType string `json:"jobTrainingType"`
	NumPs           int    `json:"numPs"`
	NumPsWorker     int    `json:"numPsWorker"`

	FrameworkType string `json:"frameworkType"`
}

type EndpointsRsp struct {
	IdentityFile  string      `json:"identityFile"`
	EndpointsInfo interface{} `json:"endpointsInfo"`
}

type AddEndportReq struct {
	JobId     string      `json:"jobId"`
	Endpoints interface{} `json:"endpoints"`
}
