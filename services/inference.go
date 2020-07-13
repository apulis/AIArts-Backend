package services

import (
	"encoding/json"
	"errors"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"strconv"
)

func PostInferenceJob(inference models.PostInference) (string,error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: inference,
	}
	resp, err := grequests.Post(BackendUrl+"/apis/PostInferenceJob", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return "",errors.New(string(resp.StatusCode))
	}
	var jobRes models.InferenceJobResp
	json.Unmarshal(resp.Bytes(),&jobRes)
	return jobRes.JobId,err
}

func ListInferenceJob(jobOwner string,vcName string,queryStringParameters models.QueryStringParametersV2) (interface{},error){
	BackendUrl = configs.Config.Infer.BackendUrl
	if queryStringParameters.PageNum==0 {
		queryStringParameters.PageNum = 1
	}
	if queryStringParameters.PageSize==0 {
		queryStringParameters.PageSize = 5
	}
	resp, err := grequests.Get(BackendUrl+"/apis/ListInferenceJobV2?jobOwner="+jobOwner+"&vcName="+vcName+"&page="+
		strconv.Itoa(queryStringParameters.PageNum)+"&size="+strconv.Itoa(queryStringParameters.PageSize), nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var jobs interface{}
	json.Unmarshal(resp.Bytes(),&jobs)
	return jobs,err
}

func GetAllSupportInference() (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetAllSupportInference", nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var inferences interface{}
	json.Unmarshal(resp.Bytes(),&inferences)
	return inferences,err
}

func GetAllDevice(userName string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetAllDevice?userName="+userName, nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var devices interface{}
	json.Unmarshal(resp.Bytes(),&devices)
	return devices,err
}

func GetJobDetail(userName string,jobId string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetInferenceJobDetail?userName="+userName+"&jobId="+jobId,nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var jobDetail interface{}
	json.Unmarshal(resp.Bytes(),&jobDetail)
	return jobDetail,err
}

func GetJobLog(userName string,jobId string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetJobLog?userName="+userName+"&jobId="+jobId,nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var jobLog interface{}
	json.Unmarshal(resp.Bytes(),&jobLog)
	return jobLog,err
}

func GetJobStatus(jobId string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetJobStatus?&jobId="+jobId,nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var jobLog interface{}
	json.Unmarshal(resp.Bytes(),&jobLog)
	return jobLog,err
}

func Infer(jobId string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	fd, err := grequests.FileUploadFromDisk("./jobid")
	logger.Info(fd)
	ro := &grequests.RequestOptions{
		Files: fd,
	}
	resp, err := grequests.Post(BackendUrl+"/apis/Infer?&jobId="+jobId,ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var respImage interface{}
	json.Unmarshal(resp.Bytes(),&respImage)
	return respImage,err
}