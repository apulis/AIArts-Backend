package services

import (
	"errors"
	"fmt"
	urllib "net/url"
	"strings"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func LsEdgeInferences(user string, req models.LsEdgeInferencesReq) ([]models.ConversionJob, int, error) {

	url := fmt.Sprintf("%s/ListModelConversionJob?vcName=%s&jobOwner=%s&num=%d&size=%d",
		configs.Config.DltsUrl, req.VCName, user, req.PageNum, req.PageSize)

	//urlencode改为%20
	if req.JobName != "" {
		url = url + fmt.Sprintf("&jobName=%s", urllib.PathEscape(req.JobName))
	}

	if req.ModelConversionType != "" {
		url = url + fmt.Sprintf("&convType=%s", req.ModelConversionType)
	}

	if req.OrderBy != "" {
		url = url + fmt.Sprintf("&orderBy=%s", req.OrderBy)
	}

	if req.Order != "" {
		url = url + fmt.Sprintf("&order=%s", req.Order)
	}

	if req.JobStatus != "" {
		url = url + fmt.Sprintf("&jobStatus=%s", urllib.PathEscape(req.JobStatus))
	}

	if req.ModelConversionStatus != "" {
		url = url + fmt.Sprintf("&convStatus=%s", urllib.PathEscape(req.ModelConversionStatus))
	}

	var resp models.ConversionList
	var res []models.ConversionJob

	err := DoRequest(url, "GET", nil, nil, &resp)
	if err != nil {
		return res, 0, err
	}

	for _, v := range resp.QueuedJobs {
		res = append(res, v)
	}
	for _, v := range resp.RunningJobs {
		res = append(res, v)
	}
	for _, v := range resp.FinishedJobs {
		res = append(res, v)
	}

	return res, resp.Total, err
}

func CreateEdgeInference(userName string, req models.CreateEdgeInferenceReq)  (string, error) {

	url := fmt.Sprintf("%s/PostModelConversionJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = req.JobName
	params["inputPath"] = req.InputPath
	params["outputPath"] = req.OutputPath
	params["conversionType"] = req.ConversionType
	params["vcName"] = req.VCName
	params["conversionArgs"] = req.ConversionArgs

	baseImageName := "apulistech/atc:0.0.1"
	if strings.HasPrefix(req.ConversionType, "arm64") {
		params["gpuType"] = "huawei_npu_arm64"
		params["image"] = baseImageName + "-arm64"
	}else{
		params["gpuType"] = "nvidia_gpu_amd64"
		params["image"] = baseImageName + "-amd64"
	}

	params["image"] = ConvertPrivateImage(params["image"].(string)) // give image name a harbor prefix
	var res models.ConversionJobId

	err := DoRequest(url, "POST", nil, params, &res)
	if err != nil {
		return "", err
	}

	return res.JobId, err
}

func GetFDInfo() (models.FDInfo, error) {
	var fd models.FDInfo
	url := fmt.Sprintf("%s/GetFDInfo", configs.Config.DltsUrl)

	err := DoRequest(url, "GET", nil, nil, &fd)
	return fd, err
}

func SetFDInfo(username, password, reqUrl string) (bool, error) {
	url := fmt.Sprintf("%s/SetFDInfo", configs.Config.DltsUrl)
	params := make(map[string]interface{})
	params["username"] = username
	params["password"] = password
	params["url"] = reqUrl

	res := false
	err := DoRequest(url, "POST", nil, params, &res)

	return res, err
}

func GetConversionTypes() (models.ConversionTypes, error) {
	var convTypes models.ConversionTypes
	url := fmt.Sprintf("%s/GetModelConversionTypes", configs.Config.DltsUrl)

	err := DoRequest(url, "GET", nil, nil, &convTypes)
	return convTypes, err
}

func PushToFD(jobId string) error {
	var res models.PushToFDRes
	url := fmt.Sprintf("%s/PushModelToFD", configs.Config.DltsUrl)
	params := make(map[string]interface{})
	params["jobId"] = jobId

	err := DoRequest(url, "POST", nil, params, &res)
	if !res.Success {
		err = errors.New(res.Msg)
	}

	return err
}
