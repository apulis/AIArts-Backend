package services

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"github.com/apulis/AIArtsBackend/configs"
	"log"
)

var BackendUrl string

func GetProjects() ([]models.Project, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/api/projects", nil)
	log.Fatal(resp)
	var projects []models.Project
	json.Unmarshal(resp.Bytes(),&projects)
	return projects,err
}

func DeleteProject(projectId string) error {
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId, nil)
	log.Println(resp)
	return err
}

func AddProject(project models.Project) error {
	ro := &grequests.RequestOptions{
		JSON: project,
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects", ro)
	log.Println(resp)
	return err
}

func UpdateProject(project models.Project,projectId string) error {
	ro := &grequests.RequestOptions{
		JSON: project,
	}
	resp, err := grequests.Patch(BackendUrl+"/api/projects/"+projectId, ro)
	log.Println(resp)
	return err
}

func GetDatasets(projectId string) ([]models.DataSet,error) {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId, nil)
	var datasets []models.DataSet
	json.Unmarshal(resp.Bytes(),&datasets)
	return datasets,err
}

func AddDataset(projectId string, dataset models.DataSet) error {
	ro := &grequests.RequestOptions{
		JSON: dataset,
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId, ro)
	log.Println(resp)
	return err
}

func GetDatasetInfo(projectId string,dataSetId string) (models.DataSet,error) {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId, nil)
	var dataset models.DataSet
	json.Unmarshal(resp.Bytes(),&dataset)
	return dataset,err
}

func UpdateDataSet(projectId string,dataSetId string,dataset models.DataSet) error {
	ro := &grequests.RequestOptions{
		JSON: dataset,
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId, ro)
	log.Println(resp)
	return err
}

func RemoveDataSet(projectId string,dataSetId string) error {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId, nil)
	log.Println(resp)
	return err
}

func GetTasks(projectId string,dataSetId string) ([]byte,error) {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId+"/tasks", nil)
	return resp.Bytes(),err
}

func GetNextTask(projectId string,dataSetId string,taskId string) (string,error) {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId+"/tasks/next/"+taskId, nil)
	return resp.String(),err
}

func GetOneTask(projectId string,dataSetId string,taskId string) (string,error) {
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId+"/tasks/annotations/"+taskId, nil)
	return resp.String(),err
}

func PostOneTask(projectId string,dataSetId string,taskId string,value string) error {
	ro := &grequests.RequestOptions{
		JSON: value,
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	log.Println(resp)
	return err
}
