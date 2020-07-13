package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"strings"
)

func ValidHomePath(userName, path string) bool {
	path = strings.TrimSpace(path)
	usrHome := "/home/"+userName

	if !strings.HasPrefix(path, usrHome) {
		return false
	}

	return true
}

func GetAllTraining(userName string, page, size int) ([] *models.Training, int, int, error) {

	url := fmt.Sprintf("%s/ListJobsV2?userName=%s&jobOwner=%s&num=%d&vcName=%s",
						configs.Config.DltsUrl, userName, userName, 1000, "atlas")

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all training err[%+v]", err)
		return nil, 0, 0, err
	}

	trainings := make([] *models.Training, 0)
	for _, v:= range jobList.RunningJobs {
		trainings = append(trainings, &models.Training{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			DeviceType:  v.JobParams.GpuType,
			DeviceNum:   v.JobParams.Resourcegpu,
			CodePath:    v.JobParams.DataPath,
			StartupFile: "",
			OutputPath:  v.JobParams.WorkPath,
			DatasetPath: v.JobParams.DataPath,
			Params:      nil,
			Status:      v.JobStatus,
			CreateTime:  v.JobTime,
			Desc:        "",
		})
	}

	for _, v:= range jobList.FinishedJobs {
		trainings = append(trainings, &models.Training{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			DeviceType:  v.JobParams.GpuType,
			DeviceNum:   v.JobParams.Resourcegpu,
			CodePath:    v.JobParams.DataPath,
			StartupFile: "",
			OutputPath:  v.JobParams.WorkPath,
			DatasetPath: v.JobParams.DataPath,
			Params:      nil,
			Status:      v.JobStatus,
			CreateTime:  v.JobTime,
			Desc:        "",
		})
	}

	return trainings, len(trainings), 1, nil
}

func CreateTraining(userName string, training models.Training) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string] interface{})

	params["userName"] = userName
	params["jobName"] = training.Name
	params["jobType"] = "training"
	params["image"] = training.Engine
	params["gpuType"] = training.DeviceType
	params["resourcegpu"] = training.DeviceNum
	params["DeviceNum"] = training.DeviceNum

	params["CodePath"] = training.CodePath
	//params["cmd"] = "sleep 30m"  // use StartupFile, params instead

	params["cmd"] = training.StartupFile
	for k, v := range training.Params {
		params["cmd"] = " --" + k + " " + v
	}

	params["OutputPath"] = ""  // use OutputPath instead
	params["dataPath"] = training.DatasetPath
	params["Desc"] = training.Desc
	params["containerUserId"] = 0
	params["jobtrainingtype"] = "RegularJob"
	params["preemptionAllowed"] = false
	params["workPath"] = ""

	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "job"

	params["hostNetwork"] = false
	params["isPrivileged"] = false
	params["interactivePorts"] = false

	params["vcName"] = "atlas"
	params["team"] = "atlas"

	id := &models.JobId{}
	err := DoRequest(url, "POST", nil, params, id)

	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return "", err
	}

	return id.Id, nil
}

func DeleteTraining(userName, id string) error {

	url := fmt.Sprintf("%s/KillJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string] interface{})

	job := &models.Job{}
	err := DoRequest(url, "GET", nil, params, job)

	if err != nil {
		fmt.Printf("delete training err[%+v]\n", err)
		return err
	}

	return nil
}

func GetTraining(userName, id string) (*models.Training, error) {

	url := fmt.Sprintf("%s/GetJobDetailV2?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string] interface{})

	job := &models.Job{}
	training := &models.Training{}

	err := DoRequest(url, "GET", nil, params, job)
	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return nil, err
	}

	training.Id = job.JobId
	training.Name = job.JobName
	training.Engine = job.JobParams.Image
	training.DeviceNum = job.JobParams.Resourcegpu
	training.DeviceType = job.JobParams.GpuType
	training.Status = job.JobStatus
	training.CreateTime = job.JobTime

	training.Params = nil
	training.DatasetPath = ""
	training.StartupFile = ""
	training.CodePath = ""
	training.OutputPath = ""
	training.Desc = ""

	return training, nil
}

func GetTrainingLog(userName, id string) (*models.JobLog, error) {

	url := fmt.Sprintf("%s/GetJobLog?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	jobLog := &models.JobLog{}

	err := DoRequest(url, "GET", nil, nil, jobLog)
	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return nil, err
	}

	return jobLog, nil
}
