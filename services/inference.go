package services

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"errors"
)

func PostInferenceJob(inference models.PostInference) error {
	BackendUrl = configs.Config.Infer.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: inference,
	}
	resp, err := grequests.Post(BackendUrl+"/apis/PostInferenceJob", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func ListInferenceJob(jobOwner string,vcName string,num string) (interface{},error){
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/ListInferenceJob?jobOwner="+jobOwner+"&vcName="+vcName+"&num"+num, nil)
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
	resp, err := grequests.Get(BackendUrl+"/apis/GetJobDetailV2?userName="+userName+"&jobId="+jobId,nil)
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

func Infer(jobId string,image interface{}) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/Infer?&jobId="+jobId,nil)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(string(resp.StatusCode))
	}
	var respImage interface{}
	json.Unmarshal(resp.Bytes(),&respImage)
	return respImage,err
}