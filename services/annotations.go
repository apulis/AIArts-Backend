package services

import (
	"encoding/json"
	"errors"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"net/url"
	"strconv"
)

var BackendUrl string

func GetProjects(queryStringParameters models.QueryStringParameters) ([]models.Project, int, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token, "Host": "127.0.0.1"},
	}

	name, err := url.Parse(queryStringParameters.GetName())
	url := BackendUrl + "/api/projectss?page=" + strconv.Itoa(queryStringParameters.GetPageNum()) + "&size=" + strconv.Itoa(queryStringParameters.GetPageSize()) + "&name=" + name.String()
	if queryStringParameters.OrderBy != "" {
		url += "&orderBy=" + queryStringParameters.OrderBy
	}
	if queryStringParameters.OrderBy != "" {
		url += "&order=" + queryStringParameters.Order
	}
	resp, err := grequests.Get(url, ro)
	logger.Info(string(resp.StatusCode), resp.StatusCode, strconv.Itoa(resp.StatusCode))
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, 0, errors.New("remote server return code is " + (strconv.Itoa(resp.StatusCode)))
	}
	var projects models.ProjectsReq
	json.Unmarshal(resp.Bytes(), &projects)
	return projects.Projects, projects.TotalCount, err
}

func DeleteProject(projectId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	return err
}

func AddProject(project models.Project) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    project,
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects", ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	return err
}

func UpdateProject(project models.Project, projectId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    project,
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Patch(BackendUrl+"/api/projects/"+projectId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	return err
}

func GetDatasets(projectId string, queryStringParameters models.QueryStringParameters) ([]models.DataSet, int, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token, "Host": "127.0.0.1"},
	}
	name, err := url.Parse(queryStringParameters.GetName())
	requrl := BackendUrl + "/api/projects/" + projectId + "/datasets?page=" + strconv.Itoa(queryStringParameters.GetPageNum()) + "&size=" + strconv.Itoa(queryStringParameters.GetPageSize()) + "&name=" + name.String()
	if queryStringParameters.OrderBy != "" {
		requrl += "&orderBy=" + queryStringParameters.OrderBy
	}
	if queryStringParameters.OrderBy != "" {
		requrl += "&order=" + queryStringParameters.Order
	}
	resp, err := grequests.Get(requrl, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, 0, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var datasets models.DatasetsReq
	json.Unmarshal(resp.Bytes(), &datasets)
	return datasets.Datasets, datasets.TotalCount, err
}

func AddDataset(projectId string, dataset models.UpdateDataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    dataset,
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	// add bind dataset record
	var datasetRes models.AddDatasetReq
	json.Unmarshal(resp.Bytes(), &datasetRes)
	ro2 := &grequests.RequestOptions{
		JSON:    map[string]string{"platform": "label", "id": datasetRes.DatasetId},
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp2, err := grequests.Post("http://127.0.0.1:"+strconv.Itoa(configs.Config.Port)+"/ai_arts/api/datasets/"+strconv.Itoa(dataset.DataSetBindId)+"/bind", ro2)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp2.StatusCode, resp2.String())
		return errors.New("remote server return code is " + (string(resp2.StatusCode)))
	}
	return err
}

func GetDatasetInfo(projectId string, dataSetId string) (models.DataSet, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return models.DataSet{}, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var dataset models.DatasetReq
	json.Unmarshal(resp.Bytes(), &dataset)
	return dataset.Info, err
}

func UpdateDataSet(projectId string, dataSetId string, dataset models.UpdateDataSet) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    dataset,
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Patch(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	return err
}

func RemoveDataSet(projectId string, dataSetId string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    "\"" + dataSetId + "\"",
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Delete(BackendUrl+"/api/projects/"+projectId+"/datasets", ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	// delete bind dataset record
	var datasetRes models.DeleteDatasetReq
	json.Unmarshal(resp.Bytes(), &datasetRes)
	ro2 := &grequests.RequestOptions{
		JSON:    map[string]string{"platform": "label", "id": dataSetId},
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp2, err := grequests.Post("http://127.0.0.1:"+strconv.Itoa(configs.Config.Port)+"/ai_arts/api/datasets/"+datasetRes.DataSetBindId+"/unbind", ro2)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp2.StatusCode, resp2.String())
		return errors.New("remote server return code is " + (string(resp2.StatusCode)))
	}
	return err
}

func GetTasks(projectId string, dataSetId string, queryStringParameters models.QueryStringParamInterface) ([]interface{}, int, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks?page="+strconv.Itoa(queryStringParameters.GetPageNum())+"&size="+strconv.Itoa(queryStringParameters.GetPageSize()), ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, 0, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var taskList models.TasksList
	json.Unmarshal(resp.Bytes(), &taskList)
	return taskList.TaskList, taskList.TotalCount, err
}

func GetNextTask(projectId string, dataSetId string, taskId string) (interface{}, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/next/"+taskId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return "", errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var nextTask models.NextTask
	json.Unmarshal(resp.Bytes(), &nextTask)
	return nextTask.Next, err
}

func GetPreviousTask(projectId string, dataSetId string, taskId string) (interface{}, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/previous/"+taskId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return "", errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var nextTask models.PreviousTask
	json.Unmarshal(resp.Bytes(), &nextTask)
	return nextTask.Previous, err
}

func GetOneTask(projectId string, dataSetId string, taskId string) (interface{}, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return "", errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var oneTask models.OneTask
	json.Unmarshal(resp.Bytes(), &oneTask)
	return oneTask.Annotations, err
}

func PostOneTask(projectId string, dataSetId string, taskId string, value string) error {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		JSON:    value,
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Post(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/annotations/"+taskId, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	return err
}

func GetDataSetLabels(projectId string, dataSetId string) (interface{}, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	resp, err := grequests.Get(BackendUrl+"/api/projects/"+projectId+"/datasets/"+dataSetId+"/tasks/labels", ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var labels models.LabelReq
	json.Unmarshal(resp.Bytes(), &labels)
	return labels.Annotations, err
}

func ConvertDataFormat(convert models.ConvertDataFormat) (interface{}, error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	ro := &grequests.RequestOptions{
		JSON: convert,
	}
	resp, err := grequests.Post(BackendUrl+"/apis/ConvertDataFormat", ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var ret interface{}
	json.Unmarshal(resp.Bytes(), &ret)
	return ret, err
}

func ConvertSupportFormat(projectId string, dataSetId string) (interface{}, error) {
	ret := []string{"coco"}
	return ret, nil
}

func ListAllDatasets(queryStringParameters models.QueryStringParamInterface) ([]models.DataSet, int, error) {
	BackendUrl = configs.Config.Anno.BackendUrl
	ro := &grequests.RequestOptions{
		Headers: map[string]string{"Authorization": "Bearer " + configs.Config.Token},
	}
	url := BackendUrl + "/api/listDatasets?page=" + strconv.Itoa(queryStringParameters.GetPageNum()) + "&size=" + strconv.Itoa(queryStringParameters.GetPageSize())
	resp, err := grequests.Get(url, ro)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return nil, 0, errors.New("remote server return code is " + (string(resp.StatusCode)))
	}
	var datasets models.DatasetsReq
	json.Unmarshal(resp.Bytes(), &datasets)
	return datasets.Datasets, datasets.TotalCount, err
}
