package services

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
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
	NumPs       int               `json:"numPs"`
	NumPsWorker int               `json:"numPsWorker"`
}

func ListModelSets(page, count int, orderBy, order string, isAdvance bool, name, status, use, username string) ([]models.Modelset, int, error) {

	offset := count * (page - 1)
	limit := count
	return models.ListModelSets(offset, limit, orderBy, order, isAdvance, name, status, use, username)
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
func GetPanel(use, username string) (interface{}, error) {
	//获取可用的数据集
	datasets, total, err := ListDatasets(1, 999, "created_at", "desc", "", "all", true, username)
	_, _, datasets = AppendAnnoDataset(datasets, total, 1, 999, "created_at", "desc")
	if err != nil {
		return "", err
	}
	//分类获取panel
	var modelset models.Modelset
	if use == "Avisualis_Classfication" {
		modelset, err = GetModelset(10001)
	} else if use == "Avisualis_ObjectDetection" {
		modelset, err = GetModelset(10001)
	} else if use == "Avisualis_SemanticSegmentation" {
		modelset, err = GetModelset(10002)
	}
	panelJson, err := gabs.ParseJSON([]byte(modelset.Description))
	//生成panel节点
	input := gabs.New()
	children := gabs.New()
	for _, dataset := range datasets {
		config := gabs.New()
		config.Set(dataset.Path, "key")
		config.Set("disabled", "type")
		config.Set("./", "value")
		children.ArrayAppend(config, dataset.Name)
	}
	input.Set("Input", "name")
	input.ArrayAppend(children, "children")
	panelJson.S("panel").SetIndex(input, 0)

	//加入启动训练任务所需要的节点
	panelJson.Set("CodePath", "name")
	panelJson.Set("Engine", "name")
	panelJson.Set("StartupFile", "name")
	if err != nil {
		return "", err
	}
	return panelJson.Data(), nil
}
