package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"math/rand"
	"net/url"
	"strings"
	"time"
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

func GetAllCodeEnv(userName string, req models.GetAllJobsReq) ([]*models.CodeEnvItem, int, int, error){

	url := fmt.Sprintf(`%s/ListJobsV3?userName=%s&jobOwner=%s&vcName=%s&jobType=%s&pageNum=%d&pageSize=%d&jobStatus=%s&searchWord=%s&orderBy=%s&order=%s`,
				configs.Config.DltsUrl, userName, userName, req.VCName,
				models.JobTypeCodeEnv,
				req.PageNum, req.PageSize, req.JobStatus, url.QueryEscape(req.SearchWord),
				req.OrderBy, req.Order)

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all code err[%+v]", err)
		return nil, 0, 0, err
	}

	codes := make([]*models.CodeEnvItem, 0)
	for _, v := range jobList.AllJobs {
		codes = append(codes, &models.CodeEnvItem{
			Id:         v.JobId,
			Name:       v.JobName,
			Engine:     v.JobParams.Image,
			CodePath:   v.JobParams.CodePath,
			Status:     v.JobStatus,
			CreateTime: v.JobTime,
			JupyterUrl: "",
			Desc:       v.JobParams.Desc,
		})
	}

	totalJobs := jobList.Meta.TotalJobs
	totalPages := totalJobs / req.PageSize

	if (totalJobs % req.PageSize) != 0 {
		totalPages += 1
	}

	return codes, totalJobs, totalPages, nil
}

func CreateCodeEnv(userName string, codeEnv models.CreateCodeEnv) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = codeEnv.Name
	params["jobType"] = models.JobTypeCodeEnv

	params["image"] = ConvertImage(codeEnv.Engine)
	params["gpuType"] = codeEnv.DeviceType
	params["resourcegpu"] = codeEnv.DeviceNum

	params["codePath"] = codeEnv.CodePath
	params["desc"] = codeEnv.Desc

	params["cmd"] = "sleep infinity"

	params["containerUserId"] = 0
	params["jobtrainingtype"] = codeEnv.JobTrainingType //"RegularJob"
	params["preemptionAllowed"] = false
	params["workPath"] = ""

	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "job"

	if codeEnv.JobTrainingType == "PSDistJob" {
		params["hostNetwork"] = true
	} else {
		params["hostNetwork"] = false
	}

	params["isPrivileged"] = false
	params["interactivePorts"] = false

	params["numpsworker"] = codeEnv.NumPsWorker
	params["numps"] = codeEnv.NumPs

	params["vcName"] = codeEnv.VCName
	params["team"] = codeEnv.VCName

	id := &models.JobId{}
	err := DoRequest(url, "POST", nil, params, id)

	if err != nil {
		fmt.Printf("create codeEnv err[%+v]\n", err)
		return "", err
	}

	// create endpoints
	url = fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id.Id)
	req := &models.CreateEndpointsReq{}
	ret := &models.CreateEndpointsRsp{}

	req.Endpoints = append(req.Endpoints, "ipython")
	req.JobId = id.Id

	err = DoRequest(url, "POST", nil, req, ret)
	if err != nil {
		fmt.Printf("create endpoints err[%+v]\n", err)
		return "", err
	}

	return id.Id, nil
}

func DeleteCodeEnv(userName, id string) error {

	url := fmt.Sprintf("%s/KillJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string]interface{})

	job := &models.Job{}
	err := DoRequest(url, "GET", nil, params, job)

	if err != nil {
		fmt.Printf("delete training err[%+v]\n", err)
		return err
	}

	return nil
}

func GetJupyterPath(userName, id string) (error, *models.EndpointWrapper) {

	url := fmt.Sprintf("%s/endpoints?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	fmt.Println(url)

	rspData := make([]models.Endpoint, 0)
	err := DoRequest(url, "GET", nil, nil, &rspData)

	if err != nil {
		fmt.Printf("get jupyter path err[%+v]\n", err)
		return err, nil
	}

	appRspData := &models.EndpointWrapper{}
	for _, v := range rspData {
		if strings.ToLower(v.Name) == "ipython" {
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
