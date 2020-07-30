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
	EngineType   string `json:"engineType"`
	DeviceType   string `json:"deviceType"`
	DeviceNum    int    `json:"deviceNum"`
	StartupFile  string `json:"startupFile"`
	OutputPath   string `json:"outputPath"`
	DatasetPath  string `json:"datasetPath"`
	DatasetName  string `json:"datasetName"`
	ArgumentPath string `json:"argumentPath"`
	CodePath     string `json:"codePath"`

	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

func ListModelSets(page, count int, orderBy, order string, isAdvance bool, name, status, username string) ([]models.Modelset, int, error) {

	offset := count * (page - 1)
	limit := count
	return models.ListModelSets(offset, limit, orderBy, order, isAdvance, name, status, username)
}

func CreateModelset(isAdvance bool, name, description, creator, version, use, jobId,
	dataFormat string, arguments map[string]string, engine, precision, modelPath, argumentPath string) error {
	var size int64
	//获取预制模型的文件size
	if modelPath != "" {
		modelSize, err := GetDirSize(modelPath)
		if err != nil {
			return fmt.Errorf("the model path %s is invaild", modelPath)
		}
		size = modelSize
	} else {
		size = int64(0)
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
		Engine:       engine,
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
