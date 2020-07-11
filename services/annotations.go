package services

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"github.com/apulis/AIArtsBackend/configs"
	"errors"
)

var BackendUrl string

func GetProjects() ([]models.Project, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return nil,errors.New(string(resp.StatusCode))
	}
	var projects models.ProjectsReq
	json.Unmarshal(resp.Bytes(),&projects)
	return projects.Projects,err
}

func DeleteProject(projectId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func AddProject(project models.Project) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: project,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func UpdateProject(project models.Project,projectId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: project,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Patch(BackendUrl+"/api/projects/"+projectId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func GetDatasets(projectId string) ([]models.DataSet,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	logger.Info(resp.String())
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return nil,errors.New(string(resp.StatusCode))
	}
	var datasets models.DatasetsReq
	json.Unmarshal(resp.Bytes(),&datasets)
	return datasets.Datasets,err
}

func AddDataset(projectId string, dataset models.DataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: dataset,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func GetDatasetInfo(projectId string,dataSetId string) (models.DataSet,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	var dataset models.DataSet
	json.Unmarshal(resp.Bytes(),&dataset)
	return dataset,err
}

func UpdateDataSet(projectId string,dataSetId string,dataset models.DataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: dataset,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func RemoveDataSet(projectId string,dataSetId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode)
		return errors.New(string(resp.StatusCode))
	}
	return err
}

func GetTasks(projectId string,dataSetId string) ([]byte,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks", ro)
	return resp.Bytes(),err
}

func GetNextTask(projectId string,dataSetId string,taskId string) (string,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/next/"+taskId, ro)
	return resp.String(),err
}

func GetOneTask(projectId string,dataSetId string,taskId string) (string,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	return resp.String(),err
}

func PostOneTask(projectId string,dataSetId string,taskId string,value string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: value,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	logger.Info(resp)
	return err
}
