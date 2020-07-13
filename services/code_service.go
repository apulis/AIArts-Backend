package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"time"
	"math/rand"
	"github.com/apulis/AIArtsBackend/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetAllCodeEnv(userName string, page, size int) ([] *models.CodeEnvItem, int, int, error) {

	url := fmt.Sprintf("%s/ListJobsV2?userName=%s&jobOwner=%s&num=%d&vcName=%s",
		configs.Config.DltsUrl, userName, userName, 1000, "atlas")

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all code err[%+v]", err)
		return nil, 0, 0, err
	}

	codes := make([] *models.CodeEnvItem, 0)
	for _, v:= range jobList.RunningJobs {
		codes = append(codes, &models.CodeEnvItem{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			CodePath:    v.JobParams.DataPath,
			Status:      v.JobStatus,
			CreateTime:  v.JobTime,
			JupyterUrl:  "",
			Desc:        "",
		})
	}

	for _, v:= range jobList.FinishedJobs {
		codes = append(codes, &models.CodeEnvItem{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			CodePath:    v.JobParams.DataPath,
			Status:      v.JobStatus,
			CreateTime:  v.JobTime,
			JupyterUrl:  "",
			Desc:        "",
		})
	}

	return codes, len(codes), 1, nil
}

func CreateCodeEnv(userName string, codeEnv models.CreateCodeEnv) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string] interface{})

	params["userName"] = userName
	params["jobName"] = codeEnv.Name
	params["jobType"] = "codeEnv"
	params["image"] = codeEnv.Engine
	params["gpuType"] = codeEnv.DeviceType
	params["resourcegpu"] = codeEnv.DeviceNum
	params["DeviceNum"] = codeEnv.DeviceNum

	params["CodePath"] = codeEnv.CodePath
	params["cmd"] = "sleep infinity"

	params["OutputPath"] = ""  // use OutputPath instead
	params["dataPath"] = ""
	params["Desc"] = codeEnv.Desc
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
		fmt.Printf("create codeEnv err[%+v]\n", err)
		return "", err
	}

	return id.Id, nil
}

func DeleteCodeEnv(userName, id string) error {

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

