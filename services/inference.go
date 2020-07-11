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