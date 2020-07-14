package services

import (
	"fmt"
	"strings"
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

	url := fmt.Sprintf("%s/ListJobsV3?userName=%s&jobOwner=%s&num=%d&vcName=%s&jobType=%s&jobStatus=all",
						configs.Config.DltsUrl, userName, userName, 1000, "atlas", models.JobTypeCodeEnv)

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
			Desc:        v.JobParams.Desc,
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
			Desc:        v.JobParams.Desc,
		})
	}

	return codes, len(codes), 1, nil
}

func CreateCodeEnv(userName string, codeEnv models.CreateCodeEnv) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string] interface{})

	params["userName"] = userName
	params["jobName"] = codeEnv.Name
	params["jobType"] = models.JobTypeCodeEnv
	params["image"] = codeEnv.Engine

	params["gpuType"] = codeEnv.DeviceType
	params["resourcegpu"] = codeEnv.DeviceNum

	params["codePath"] = codeEnv.CodePath
	params["desc"] = codeEnv.Desc

	params["cmd"] = "sleep infinity"

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

	// create endpoints
	url = fmt.Sprintf("%s/endpoints?userName=%s", configs.Config.DltsUrl, userName)
	endpoints := &models.EndpointsReq{}
	ret := &models.EndpointsRet{}

	endpoints.Endpoints = append(endpoints.Endpoints, "ipython")
	endpoints.JobId = id.Id

	err = DoRequest(url, "POST", nil, params, ret)
	if err != nil {
		fmt.Printf("create endpoints err[%+v]\n", err)
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

func GetJupyterPath(userName, id string) (error, *models.EndpointsDetail) {

	url := fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)

	req := &models.EndpointsReq{}
	endpointDetail := &models.EndpointsDetail{}

	req.Endpoints = append(req.Endpoints, "ipython")
	req.JobId = id

	err := DoRequest(url, "GET", nil, req, endpointDetail)
	if err != nil {
		fmt.Printf("get jupyter path err[%+v]\n", err)
		return err, nil
	}

	for _, v := range endpointDetail.Endpoints {
		if strings.ToLower(v.Name) == "ipython" {
			v.AccessPoint = fmt.Sprintf("http://%s.%s/endpoints/%s/", v.NodeName, v.Domain, v.Port)
			break
		}
	}

	return nil, endpointDetail
}
