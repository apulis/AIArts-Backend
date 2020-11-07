package models

type FDInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

type ConversionTypes struct {
	ConversionTypes []string `json:"conversionTypes"`
}

type ConversionList struct {
	FinishedJobs []ConversionJob `json:"finishedJobs"`
	Meta         JobMeta         `json:"meta"`
	QueuedJobs   []ConversionJob `json:"queuedJobs"`
	RunningJobs  []ConversionJob `json:"runningJobs"`
	Total        int             `json:"total"`
}

type ConversionJob struct {
	JobId            string    `json:"jobId"`
	JobName          string    `json:"jobName"`
	JobParams        JobParams `json:"jobParams"`
	JobStatus        string    `json:"jobStatus"`
	JobTime          string    `json:"jobTime"`
	JobType          string    `json:"jobType"`
	Priority         int       `json:"priority"`
	UserName         string    `json:"userName"`
	VcName           string    `json:"vcName"`
	InputPath        string    `json:"inputPath"`
	OutputPath       string    `json:"outputPath"`
	ConversionStatus string    `json:"modelconversionStatus"`
	ConversionType   string    `json:"modelconversionType"`
}

type ConversionJobId struct {
	JobId string `json:"jobId"`
}

type PushToFDRes struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

type CreateEdgeInferenceReq struct {
	JobName        string                 `json:"jobName" binding:"required"`
	InputPath      string                 `json:"inputPath" binding:"required"`
	OutputPath     string                 `json:"outputPath" binding:"required"`
	ConversionType string                 `json:"conversionType" binding:"required"`
	ConversionArgs map[string]interface{} `json:"conversionArgs" binding:"required"`
	VCName 		   string  				  `json:"vcName" binding:"required"`
}

type LsEdgeInferencesReq struct {
	PageNum               int    `form:"pageNum,default=1"`
	PageSize              int    `form:"pageSize,default=10"`
	JobName               string `form:"jobName"`
	ModelConversionType   string `form:"modelconversionType"`
	OrderBy               string `form:"orderBy"`
	Order                 string `form:"order,default=desc" binding:"oneof=desc asc"`
	JobStatus             string `form:"jobStatus"`
	ModelConversionStatus string `form:"modelconversionStatus"`
	VCName				  string `form:"vcName"`
}
