package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/models"
)

const (
	MODELSET_STATUS_NORMAL   = "normal"
	MODELSET_STATUS_DELETING = "deleting"
)

type CreateEvaluationReq struct {
	EngineType  string            `json:"engineType"`
	DeviceType  string            `json:"deviceType"`
	DeviceNum   int               `json:"deviceNum"`
	StartupFile string            `json:"startupFile"`
	OutputPath  string            `json:"outputPath"`
	DatasetPath string            `json:"datasetPath"`
	DatasetName string            `json:"datasetName"`
	ParamPath   string            `json:"paramPath"`
	CodePath    string            `json:"codePath"`
	Name        string            `json:"name"`
	Params      map[string]string `json:"params"`
}

func ListModelSets(page, count int, orderBy, order string, isAdvance bool, name, status, username string) ([]models.Modelset, int, error) {

	offset := count * (page - 1)
	limit := count
	return models.ListModelSets(offset, limit, orderBy, order, isAdvance, name, status, username)
}

func CreateModelset(name, description, creator, version, jobId, codePath, paramPath string, isAdvance bool,
	use string, size int64, dataFormat, datasetName, datasetPath string, params map[string]string, engine, precision, outputPath, startupFile string) error {

	//只能创建非预置模型
	modelset := models.Modelset{
		Name:        name,
		Description: description,
		Creator:     creator,
		Version:     version,
		JobId:       jobId,
		Status:      MODELSET_STATUS_NORMAL,
		IsAdvance:   isAdvance,
		ParamPath:   paramPath,
	}
	if use != "" {
		var paramItem models.ParamsItem
		paramItem = params
		modelset = models.Modelset{
			Name:        name,
			Description: description,
			Creator:     creator,
			Version:     version,
			JobId:       jobId,
			Status:      MODELSET_STATUS_NORMAL,
			IsAdvance:   isAdvance,
			ParamPath:   paramPath,
			Use:         use,
			Size:        size,
			DataFormat:  dataFormat,
			DatasetPath: datasetPath,
			DatasetName: datasetName,
			Params:      &paramItem,
			Engine:      engine,
			Precision:   precision,
			OutputPath:  outputPath,
			StartupFile: startupFile,
		}
	}
	//获取训练作业输出模型的类型 job
	if codePath == "" {
		job, _ := GetTraining(creator, jobId)
		var paramItem models.ParamsItem
		paramItem = job.Params
		if job != nil {
			modelset.OutputPath = job.OutputPath
			modelset.CodePath = job.CodePath
			modelset.DatasetPath = job.DatasetPath
			modelset.StartupFile = job.StartupFile
			modelset.Params = &paramItem
			modelset.Engine = job.Engine
		} else {
			return fmt.Errorf("the job id is invaild")
		}
	} else {
		modelset.CodePath = codePath
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
