package services

import (
	"fmt"
	"strings"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func CreateVisualJob(userName string, jobName string, logdir string, description string) error {
	//step1. create a job
	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = jobName
	params["jobType"] = models.JobTypeCodeEnv

	params["image"] = "apulistech/visualjob:1.0"
	params["gpuType"] = "nvidia_gpu_amd64"
	params["resourcegpu"] = 0

	params["codePath"] = logdir
	params["desc"] = description

	params["cmd"] = "sleep infinity"

	params["containerUserId"] = 0
	params["jobtrainingtype"] = models.JobTypeVisualJob
	params["preemptionAllowed"] = false
	params["workPath"] = ""

	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "job"

	params["hostNetwork"] = false

	params["isPrivileged"] = false
	params["interactivePorts"] = false

	params["vcName"] = models.DefaultVcName
	params["team"] = models.DefaultVcName

	id := &models.JobId{}
	err := DoRequest(url, "POST", nil, params, id)
	if err != nil {
		fmt.Printf("create codeEnv err[%+v]\n", err)
		return err
	}
	//step2. create endpoints
	url = fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id.Id)
	req := &models.CreateEndpointsReq{}
	ret := &models.CreateEndpointsRsp{}

	req.Endpoints = append(req.Endpoints, "tensorboard")
	req.JobId = id.Id
	req.Arguments = "{ 'tensorboard_log_dir': '" + logdir + "' }"

	err = DoRequest(url, "POST", nil, req, ret)
	if err != nil {
		fmt.Printf("create endpoints err[%+v]\n", err)
		return err
	}
	//step3. create visual job record
	visualJob := models.VisualJob{
		UserName:    userName,
		Name:        jobName,
		Status:      "pending",
		LogPath:     logdir,
		Description: description,
		RelateJobId: id.Id,
	}

	return models.CreateVisualJob(visualJob)
}

func GetAllVisualJobInfo(userName string, pageNum int, pageSize int, orderBy string, status string, jobName string, order string) ([]models.VisualJob, int, int, error) {
	//step1. renew all visual job status
	err := renewStatusInfo(userName)
	if err != nil {
		fmt.Printf("job status renew fail : err[%+v]\n", err)
		return nil, 0, 0, err
	}
	//step2. get job info and return
	jobList, err := models.GetAllVisualJobByArguments(userName, pageNum, pageSize, status, jobName, orderBy, order)
	if err != nil {
		fmt.Printf("get job list err[%+v]\n", err)
		return nil, 0, 0, err
	}
	totalJobsNum, err := models.GetVisualJobsSumCount()
	if err != nil {
		fmt.Printf("get job list count err[%+v]\n", err)
		return nil, 0, 0, err
	}
	totalPages := totalJobsNum / pageSize

	if (totalJobsNum % pageSize) != 0 {
		totalPages += 1
	}
	return jobList, totalJobsNum, totalPages, nil
}

func GetEndpointsPath(userName string, visualJobId int) (string, error) {
	visualJobDetail, err := models.GetVisualJobById(visualJobId)
	if err != nil {
		fmt.Printf("get visual job detail err[%+v]\n", err)
		return "", err
	}
	err, endpointInfo := GetTensorboardPath(userName, visualJobDetail.RelateJobId)
	if err != nil {
		fmt.Printf("get endpoint path err[%+v]\n", err)
		return "", err
	}
	return endpointInfo.AccessPoint, nil
}

func GetTensorboardPath(userName, jobId string) (error, *models.EndpointWrapper) {

	url := fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, jobId)
	fmt.Println(url)

	rspData := make([]models.Endpoint, 0)
	err := DoRequest(url, "GET", nil, nil, &rspData)

	if err != nil {
		fmt.Printf("get visual job path err[%+v]\n", err)
		return err, nil
	}

	appRspData := &models.EndpointWrapper{}
	for _, v := range rspData {
		if strings.ToLower(v.Name) == "tensorboard" {
			appRspData.Name = v.Name
			appRspData.Status = v.Status

			if v.Status == "running" {
				appRspData.AccessPoint = fmt.Sprintf("http://%s.%s/endpoints/%s/", v.NodeName, v.Domain, v.Port)
			}

			break
		}
	}

	return nil, appRspData
}

func StopVisualJob(userName string, jobId int) {

}

func renewStatusInfo(userName string) error {
	visualJobList, err := models.GetAllVisualJobByArguments(userName, 1, -1, "", "", "", "")
	if err != nil {
		fmt.Printf("get visual job  err[%+v]\n", err)
		return err
	}
	for _, job := range visualJobList {
		backgroundJobId := job.RelateJobId
		url := fmt.Sprintf("%s/GetJobDetailV2?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, backgroundJobId)
		params := make(map[string]interface{})
		backgroundJob := &models.Job{}
		err := DoRequest(url, "GET", nil, params, backgroundJob)
		if err != nil {
			fmt.Printf("get training err[%+v]\n", err)
			return err
		}
		job.Status = backgroundJob.JobStatus
		models.UpdateVisualJob(&job)
		fmt.Printf(backgroundJobId)
	}
	return nil
}
