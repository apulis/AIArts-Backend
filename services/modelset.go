package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

const (
	MODELSET_STATUS_NORMAL   = "normal"
	MODELSET_STATUS_DELETING = "deleting"
)
type CreateEvaluationReq struct {
	EngineType       string `json:"engineType"`
	DeviceType   string `json:"deviceType"`
	DeviceNum    int    `json:"deviceNum"`
	StartupFile  string `json:"startupFile"`
	OutputPath   string `json:"outputPath"`
	DatasetPath  string `json:"datasetPath"`
	DatasetName  string `json:"datasetName"`
	ArgumentPath string `json:"argumentPath"`
	Name       string `json:"name"`


}
func ListModelSets(page, count int, orderBy, order string, isAdvance bool, name, status, username string) ([]models.Modelset, int, error) {

	offset := count * (page - 1)
	limit := count
	return models.ListModelSets(offset, limit, orderBy, order, isAdvance, name, status, username)
}

func CreateModelset(isAdvance bool, name, description, creator, version, use, jobId,
	dataFormat string, arguments map[string]string, engineType, precision, modelPath, argumentPath string) error {
	size, err := GetDirSize(modelPath)
	if err != nil {
		return err
	}
	//json转换格式
	var argItem models.ArgumentsItem
	argItem = arguments
	modelset := models.Modelset{
		Name:         name,
		Description:  description,
		Creator:      creator,
		Version:      version,
		Size:         size,
		Use:          use,
		JobId:        jobId,
		Status:       MODELSET_STATUS_NORMAL,
		DataFormat:   dataFormat,
		Arguments:    &argItem,
		EngineType:   engineType,
		Precision:    precision,
		IsAdvance:    isAdvance,
		ModelPath:    modelPath,
		ArgumentPath: argumentPath,
	}
	return models.CreateModelset(modelset)
}

func UpdateModelset(id int, description string) error {
	modelset, err := models.GetModelsetById(id)
	if err != nil {
		return err
	}
	modelset.Description = description
	return models.UpdateModelset(&modelset)
}

func GetModelset(id int) (models.Modelset, error) {
	return models.GetModelsetById(id)
}

func DeleteModelset(id int) error {
	modelset, err := models.GetModelsetById(id)
	if err != nil {
		return err
	}

	modelset.Status = MODELSET_STATUS_DELETING
	err = models.UpdateModelset(&modelset)
	if err != nil {
		return err
	}

	//err = os.RemoveAll(modelset.Path)
	//if err != nil {
	//	return err
	//}
	return models.DeleteModelset(&modelset)
}
func CreateEvaluation(userName string, evaluation CreateEvaluationReq) (string, error) {
	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})
	params["userName"] = userName
	params["jobName"] = evaluation.Name
	params["jobType"] = models.JobTypeArtsEvaluation
	params["image"] = evaluation.EngineType
	params["gpuType"] = evaluation.DeviceType
	params["resourcegpu"] = evaluation.DeviceNum
	params["DeviceNum"] = evaluation.DeviceNum
	params["cmd"] = "python " + evaluation.StartupFile
	if len(evaluation.DatasetPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --checkpoint_path=" + evaluation.DatasetPath
	}
	if len(evaluation.OutputPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --eval_dir=" + evaluation.OutputPath
	}
	if len(evaluation.DatasetPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --dataset_dir=" + evaluation.OutputPath
	}
	params["startupFile"] = evaluation.StartupFile
	params["datasetPath"] = evaluation.DatasetPath
	//params["codePath"] = evaluation.CodePath
	params["outputPath"] = evaluation.OutputPath
	//params["scriptParams"] = evaluation.Params
	//params["desc"] = evaluation.Desc
	params["containerUserId"] = 0
	params["jobtrainingtype"] = "RegularJob"
	params["preemptionAllowed"] = false
	params["workPath"] = ""
	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "job"
	params["hostNetwork"] = false
	params["isPrivileged"] = false
	params["interactivePorts"] = false
	//params["numworker"] = training.NumPs
	//params["numps"] = training.NumPsWorker
	params["vcName"] = models.DefaultVcName
	params["team"] = models.DefaultVcName
	logger.Info(params)

	id := &models.JobId{}
	err := DoRequest(url, "POST", nil, params, id)

	if err != nil {
		fmt.Printf("create evaluation err[%+v]\n", err)
		return "", err
	}

	return id.Id, nil
}
