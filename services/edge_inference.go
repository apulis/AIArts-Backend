package services

import (
	"errors"
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func LsEdgeInferences(pageNum, pageSize int, user string) ([]models.ConversionJob, error) {
	fmt.Println(pageNum, pageSize)
	url := fmt.Sprintf("%s/ListModelConversionJob?vcName=%s&jobOwner=%s", configs.Config.DltsUrl, models.DefaultVcName, user)

	var resp models.ConversionList
	var res []models.ConversionJob
	err := DoRequest(url, "GET", nil, nil, &resp)
	if err != nil {
		return res, err
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

	return res, err
}

func CreateEdgeInference(jobName, inputPath, outputPath, convType, userName string) (string, error) {
	url := fmt.Sprintf("%s/PostModelConversionJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = jobName
	params["inputPath"] = inputPath
	params["outputPath"] = outputPath
	params["conversionType"] = convType
	params["vcName"] = models.DefaultVcName

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