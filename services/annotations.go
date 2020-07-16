package services

import (
	"encoding/json"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"github.com/apulis/AIArtsBackend/configs"
	"errors"
	"strconv"
)

var BackendUrl string

func GetProjects(queryStringParameters models.QueryStringParamInterface) ([]models.Project,int, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects?page="+strconv.Itoa(queryStringParameters.GetPageNum())+"&size="+strconv.Itoa(queryStringParameters.GetPageSize()), ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,0,errors.New(resp.String())
	}
	var projects models.ProjectsReq
	json.Unmarshal(resp.Bytes(),&projects)
	return projects.Projects,projects.TotalCount,err
}

func DeleteProject(projectId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
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
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
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
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
	}
	return err
}

func GetDatasets(projectId string,queryStringParameters models.QueryStringParameters) ([]models.DataSet,int,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets?page="+strconv.Itoa(queryStringParameters.GetPageNum())+"&size="+strconv.Itoa(queryStringParameters.GetPageSize()), ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,0,errors.New(resp.String())
	}
	var datasets models.DatasetsReq
	json.Unmarshal(resp.Bytes(),&datasets)
	return datasets.Datasets,datasets.TotalCount,err
}

func AddDataset(projectId string, dataset models.UpdateDataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: dataset,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
	}
	return err
}

func GetDatasetInfo(projectId string,dataSetId string) (models.DataSet,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	var dataset models.DatasetReq
	json.Unmarshal(resp.Bytes(),&dataset)
	return dataset.Info,err
}

func UpdateDataSet(projectId string,dataSetId string,dataset models.UpdateDataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: dataset,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Patch(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
	}
	return err
}

func RemoveDataSet(projectId string,dataSetId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: "\""+dataSetId+"\"",
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
	}
	return err
}

func GetTasks(projectId string,dataSetId string,queryStringParameters models.QueryStringParameters) ([]interface{},int,error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks?page="+strconv.Itoa(queryStringParameters.GetPageNum())+"&size="+strconv.Itoa(queryStringParameters.GetPageSize()), ro)
	logger.Info(strconv.Itoa(queryStringParameters.Page),strconv.Itoa(queryStringParameters.Size))
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,0,errors.New(resp.String())
	}
	var taskList models.TasksList
	json.Unmarshal(resp.Bytes(),&taskList)
	return taskList.TaskList,taskList.TotalCount,err
}

func GetNextTask(projectId string,dataSetId string,taskId string) (interface{},error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/next/"+taskId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return "",errors.New(resp.String())
	}
	var nextTask models.NextTask
	json.Unmarshal(resp.Bytes(),&nextTask)
	return nextTask.Next,err
}

func GetOneTask(projectId string,dataSetId string,taskId string) (interface{},error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return "",errors.New(resp.String())
	}
	var oneTask models.OneTask
	json.Unmarshal(resp.Bytes(),&oneTask)
	return oneTask.Annotations,err
}

func PostOneTask(projectId string,dataSetId string,taskId string,value string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: value,
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return errors.New(resp.String())
	}
	return err
}

func GetDataSetLabels(projectId string,dataSetId string) (interface{},error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization":"Bearer "+configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/labels", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(resp.String())
	}
	var labels models.LabelReq
	json.Unmarshal(resp.Bytes(),&labels)
	return labels.Annotations,err
}

func ConvertDataFormat(convert models.ConvertDataFormat) (interface{},error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: convert,
	}
	resp, err := grequests.Post(BackendUrl+"/apis/ConvertDataFormat", ro)
	if resp.StatusCode!=200 {
		logger.Error("response code is ",resp.StatusCode,resp.String())
		return nil,errors.New(resp.String())
	}
	var ret interface{}
	json.Unmarshal(resp.Bytes(),&ret)
	return ret,err
}