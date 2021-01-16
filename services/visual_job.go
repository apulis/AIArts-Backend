package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func CreateVisualJob(userName string, vcName string, jobName string, logdir string, description string) error {
	//step1. create a background job
	relateJobId, err := createBackgroundJob(userName, vcName, jobName, logdir, description)
	if err != nil {
		fmt.Printf("create background job failed : [%+v]\n", err)
		return err
	}
	//step2. create visual job record
	visualJob := models.VisualJob{
		UserName:    userName,
		Name:        jobName,
		Status:      "scheduling",
		VCName:      vcName,
		LogPath:     logdir,
		Description: description,
		RelateJobId: relateJobId,
	}
	err = models.CreateVisualJob(visualJob)
	if err != nil {
		fmt.Printf("create visual job record failed : [%+v]\n", err)
		return err
	}
	return nil
}

func GetAllVisualJobInfo(userName string, req models.GetVisualJobListReq) ([]models.VisualJob, int, int, error) {

	//step1. renew all visual job status
	err := renewStatusInfo(userName)
	if err != nil {
		fmt.Printf("job status renew fail : err[%+v]\n", err)
		return nil, 0, 0, err
	}

	//step2. get job info and return
	jobList, err := models.GetAllVisualJobByArguments(userName, req.VCName, req.PageNum, req.PageSize, req.Status, req.JobName, req.Order, req.OrderBy)
	if err != nil {
		fmt.Printf("get job list err[%+v]\n", err)
		return nil, 0, 0, err
	}

	totalJobsNum, err := models.GetVisualJobCountByArguments(userName, req.VCName, req.Status, req.JobName)
	if err != nil {
		fmt.Printf("get job list count err[%+v]\n", err)
		return nil, 0, 0, err
	}

	totalPages := totalJobsNum / req.PageSize
	if (totalJobsNum % req.PageSize) != 0 {
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

	protocol := "http"
	if len(configs.Config.ExtranetProtocol) > 0 {
		protocol = configs.Config.ExtranetProtocol
	}

	appRspData := &models.EndpointWrapper{}
	for _, v := range rspData {
		if strings.ToLower(v.Name) == "tensorboard" {
			appRspData.Name = v.Name
			appRspData.Status = v.Status

			if v.Status == "running" {
				param := models.EndpointURLCode{Port: v.Port, UserName: userName}
				val, _ := json.Marshal(param)
				appRspData.AccessPoint = fmt.Sprintf("%s://%s.%s/endpoints/%s/",
					protocol,
					v.NodeName, v.Domain,
					base64.StdEncoding.EncodeToString(val),
				)
			}

			break
		}
	}

	return nil, appRspData
}

func StopVisualJob(userName string, jobId int) error {
	targetJob, err := models.GetVisualJobById(jobId)
	if err != nil {
		logger.Error("get job detail err ", err)
		return err
	}
	backgroundJobId := targetJob.RelateJobId
	url := fmt.Sprintf("%s/KillJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, backgroundJobId)
	params := make(map[string]interface{})

	job := &models.Job{}
	err = DoRequest(url, "GET", nil, params, job)
	if err != nil {
		logger.Error("delete background job err: ", err)
		return err
	}
	targetJob.Status = "paused"
	err = models.UpdateVisualJob(&targetJob)
	if err != nil {
		logger.Error("kill background job err %s", err)
		return err
	}

	return nil
}

func PauseVisualJob(userName string, jobId int) error {
	targetJob, err := models.GetVisualJobById(jobId)
	if err != nil {
		logger.Error("get job detail err ", err)
		return err
	}
	url := fmt.Sprintf("%s/PauseJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, targetJob.RelateJobId)
	params := make(map[string]interface{})

	job := &models.Job{}
	err = DoRequest(url, "GET", nil, params, job)
	if err != nil {
		logger.Error("pause background job err: ", err)
		return err
	}
	targetJob.Status = "paused"
	err = models.UpdateVisualJob(&targetJob)
	if err != nil {
		logger.Error("update background job err %s", err)
		return err
	}

	return nil
}

func ResumeVisualJob(userName string, jobId int) error {
	targetJob, err := models.GetVisualJobById(jobId)
	if err != nil {
		logger.Error("get job detail err ", err)
		return err
	}
	url := fmt.Sprintf("%s/ResumeJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, targetJob.RelateJobId)
	params := make(map[string]interface{})

	job := &models.Job{}
	err = DoRequest(url, "GET", nil, params, job)
	if err != nil {
		logger.Error("resume background job err: ", err)
		return err
	}
	targetJob.Status = "scheduling"
	err = models.UpdateVisualJob(&targetJob)
	if err != nil {
		logger.Error("update background job err %s", err)
		return err
	}

	return nil
}

func ContinueVisualJob(userName string, vcName string, jobId int) error {
	targetJob, err := models.GetVisualJobById(jobId)
	if err != nil {
		fmt.Printf("get job detail err[%+v]\n", err)
		return err
	}
	relateJobId, err := createBackgroundJob(userName, vcName, targetJob.Name, targetJob.LogPath, targetJob.Description)
	if err != nil {
		fmt.Printf("create background job failed : [%+v]\n", err)
		return err
	}
	targetJob.RelateJobId = relateJobId
	targetJob.Status = "scheduling"
	err = models.UpdateVisualJob(&targetJob)
	if err != nil {
		fmt.Printf("update visual job info failed: [%+v]\n", err)
		return err
	}
	return nil
}

func DeleteVisualJob(userName string, jobId int) error {
	err := renewStatusInfo(userName)
	if err != nil {
		logger.Error("job status renew fail : ", err)
		return err
	}
	job, err := models.GetVisualJobById(jobId)
	if err != nil {
		logger.Error("get job detail err: ", err)
		return err
	}

	if job.Status == "running" {
		err := StopVisualJob(userName, jobId)
		if err != nil {
			logger.Error("stop job error: ", err)
			return err
		}
	}

	_, err = DeleteJob(job.RelateJobId)
	if err != nil {
		logger.Error("delete visual job info fail: ", err)
		return err
	}

	err = models.DeleteVisualJob(&job)
	if err != nil {
		logger.Error("delete visual job record error : %s", err)
		return err
	}

	return nil
}

func createBackgroundJob(userName string, vcName string, jobName string, logdir string, description string) (string, error) {
	//step1. create a job
	// * get cluster availuable gpu type, and randomly select one to run visual job
	requestClusterStatusURL := fmt.Sprintf("%s/GetVC?userName=%s&vcName=%s", configs.Config.DltsUrl, userName, vcName)
	vcInfo := &models.VcInfo{}

	err := DoRequest(requestClusterStatusURL, "GET", nil, nil, vcInfo)
	if err != nil {
		fmt.Printf("visual job process get cluster status err[%+v]\n", err)
		return "", err
	}
	var selectNodeDevice string
	for key := range vcInfo.DeviceCapacity {
		selectNodeDevice = key
		break
	}
	// create one job
	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = jobName
	params["jobType"] = models.JobTypeVisualJob

	var visualJob_image_name = "apulistech/visualjob:1.0"
	// TODO: use docker manifest to avoid suffix name
	//if find := strings.Contains(selectNodeDevice, "arm"); find {
	//	visualJob_image_name = visualJob_image_name + "-arm64"
	//}
	params["image"] = ConvertPrivateImage(visualJob_image_name)
	params["gpuType"] = selectNodeDevice
	params["resourcegpu"] = 0

	params["codePath"] = logdir
	params["desc"] = description

	params["cmd"] = "sleep infinity"

	params["containerUserId"] = 0
	params["jobtrainingtype"] = "RegularJob"
	params["preemptionAllowed"] = false
	params["workPath"] = ""

	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "./"
	params["hostNetwork"] = false

	params["isPrivileged"] = false
	params["interactivePorts"] = false

	params["vcName"] = vcName
	params["team"] = models.DefaultVcName

	id := &models.CreateJobReq{}
	err = DoRequest(url, "POST", nil, params, id)
	if err != nil {
		fmt.Printf("post dlts err[%+v]\n", err)
		return "", err
	}

	//step2. create endpoints
	url = fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id.Id)
	req := &models.CreateEndpointsReq{}
	ret := &models.CreateEndpointsRsp{}

	req.Endpoints = append(req.Endpoints, "tensorboard")
	req.JobId = id.Id
	req.Arguments = `{ "tensorboard_log_dir" : "` + logdir + `" }`

	err = DoRequest(url, "POST", nil, req, ret)
	if err != nil {
		fmt.Printf("create endpoints err[%+v]\n", err)
		return "", err
	}
	return id.Id, nil
}

func renewStatusInfo(userName string) error {

	visualJobList, err := models.GetAllVisualJobByArguments(userName, "", 1, -1, "", "", "", "")
	if err != nil {
		fmt.Printf("get visual job  err[%+v]\n", err)
		return err
	}

	for _, job := range visualJobList {
		backgroundJobId := job.RelateJobId
		if job.Status == "paused" {
			continue
		}

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
