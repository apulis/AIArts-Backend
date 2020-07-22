package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"io/ioutil"
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
		return "",errors.New(resp.String())
	}
	var jobRes models.InferenceJobResp
	json.Unmarshal(resp.Bytes(),&jobRes)
	return jobRes.JobId,err
}

func ListInferenceJob(jobOwner string,vcName string,queryStringParameters models.QueryStringParamInterface) (interface{},error){
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/ListInferenceJobV2?jobOwner="+jobOwner+"&vcName="+vcName+"&page="+
		strconv.Itoa(queryStringParameters.GetPageNum())+"&size="+strconv.Itoa(queryStringParameters.GetPageSize())+"&search="+queryStringParameters.GetName(), nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(resp.String())
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
		return nil,errors.New(resp.String())
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
		return nil,errors.New(resp.String())
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
		return nil,errors.New(resp.String())
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
		return nil,errors.New(resp.String())
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
		return nil,errors.New(resp.String())
	}
	var jobLog interface{}
	json.Unmarshal(resp.Bytes(),&jobLog)
	return jobLog,err
}

func Infer(jobId string,signature_name string,image []byte) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	//fd, err := grequests.FileUploadFromDisk("./"+jobId))
	ro := &grequests.RequestOptions{
		Files: []grequests.FileUpload{{FileName: "image",FileContents: ioutil.NopCloser(bytes.NewReader(image))}},
	}
	if signature_name == "" {
		signature_name = "predict"
	}
	resp, err := grequests.Post(BackendUrl+"/apis/Infer?&jobId="+jobId+"&signature_name="+signature_name,ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(resp.String())
	}
	var respImage interface{}
	json.Unmarshal(resp.Bytes(),&respImage)
	return respImage,err
}


func KillJob(jobId string,userName string) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/KillJob?&jobId="+jobId+"&userName="+userName,nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(resp.String())
	}
	var jobLog interface{}
	json.Unmarshal(resp.Bytes(),&jobLog)
	return jobLog,err
}