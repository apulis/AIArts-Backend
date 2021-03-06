package models

type CodeEnvItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"   binding:"required"`
	Status     string `json:"status" binding:"required"`
	Engine     string `json:"engine"  binding:"required"`
	CodePath   string `json:"codePath"`
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

	JobTrainingType string `json:"jobTrainingType"`
	NumPs           int    `json:"numPs"`
	NumPsWorker     int    `json:"numPsWorker"`
}
