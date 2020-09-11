package services

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/apulis/AIArtsBackend/models"
	"os"
)

const (
	MODELSET_STATUS_NORMAL   = "normal"
	MODELSET_STATUS_DELETING = "deleting"
)

type CreateModelsetReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" `
	JobId       string `json:"jobId"`
	CodePath    string `json:"codePath"`
	ParamPath   string `json:"paramPath"`
	IsAdvance   bool   `json:"isAdvance,default=false"`

	Use         string `json:"use"`
	Size        int64  `json:"size"`
	DataFormat  string `json:"dataFormat"`
	DatasetName string `json:"datasetName"`
	DatasetPath string `json:"datasetPath"`
	//omitempty 值为空，不编码
	Params    map[string]string `json:"params"`
	Engine    string            `json:"engine"`
	Precision string            `json:"precision"`
	//指定的模型参数路径
	// 输出文件路径
	OutputPath string `json:"outputPath"`
	//启动文件路径
	StartupFile string `json:"startupFile"`

	//用于可视化建模平台直接启动训练任务
	JobTrainingType string          `json:"jobTrainingType"`
	NumPs           int             `json:"numPs"`
	NumPsWorker     int             `json:"numPsWorker"`
	DeviceType      string          `json:"deviceType"`
	DeviceNum       int             `json:"deviceNum"`
	Nodes           []AvisualisNode `json:"nodes"`
	Arguments       []AvisualisNode `json:"arguments"`
	Edges           []AvisualisEdge `json:"edges"`
}

type AvisualisEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
type AvisualisNode struct {
	ID     string      `json:"id"`
	Name   string      `json:"name"`
	IDX    int         `json:"idx"`
	Config interface{} `json:"config"`
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

func UpdateModelset(id int, name, description, version, jobId, codePath, paramPath,
	use string, size int64, dataFormat, datasetName, datasetPath string, params map[string]string, engine, precision, outputPath, startupFile string) error {
	modelset, err := models.GetModelsetById(id)
	if err != nil {
		return err
	}
	var paramItem models.ParamsItem
	paramItem = params
	modelset.Description = description
	modelset.OutputPath = outputPath
	modelset.CodePath = codePath
	modelset.DatasetPath = datasetPath
	modelset.StartupFile = startupFile
	modelset.Params = &paramItem
	modelset.Name = name
	modelset.Version = version
	modelset.JobId = jobId
	modelset.ParamPath = paramPath
	modelset.Use = use
	modelset.Size = size
	modelset.DataFormat = dataFormat
	modelset.DatasetName = datasetName
	modelset.Precision = precision

	return models.UpdateModelset(&modelset)
}

func GetModelset(id int) (models.Modelset, error) {
	return models.GetModelsetById(id)
}
func GetModelsetByName(name string) (models.Modelset, error) {
	return models.GetModelsetByName(name)
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
	modelset, err := GetModelsetByName(use)
	panelJson, err := gabs.ParseJSON([]byte(modelset.Description))
	if err != nil {
		return "", err
	}
	//生成panel节点
	input := gabs.New()
	for _, dataset := range datasets {
		config := gabs.New()
		children := gabs.New()
		_, _ = config.Set("data_path", "key")
		_, _ = config.Set("disabled", "type")
		_, _ = config.Set(dataset.Path, "value")
		_ = children.ArrayAppend(config, dataset.Name)
		_ = input.ArrayAppend(children, "children")
	}
	_, _ = input.Set("Input", "name")
	_, _ = panelJson.S("panel").SetIndex(input, 0)

	//加入启动训练任务所需要的节点
	_, _ = panelJson.Set(modelset.CodePath, "codePath")
	_, _ = panelJson.Set(modelset.Engine, "engine")
	_, _ = panelJson.Set(modelset.StartupFile, "startupFile")
	if err != nil {
		return "", err
	}
	return panelJson.Data(), nil
}

func CreateAvisualisTraining(req CreateModelsetReq, username string) (CreateModelsetReq, error) {

	//存储节点json
	nodesBytes, _ := json.Marshal(req.Nodes)
	edgesBytes, _ := json.Marshal(req.Edges)

	//去掉nodes没用的节点并存入json
	pipelineConfigPath, err := GetModelTempPath(FILETYPE_JSON)
	f, err := os.OpenFile(pipelineConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)
	if err != nil {
		fmt.Println(pipelineConfigPath + "failed to created")
	}
	_, err = f.Write(nodesBytes)
	defer f.Close()

	//把数据传入params后端算法只需要pipeline_config
	req.Params = make(map[string]string)
	req.Params["pipeline_config"] = pipelineConfigPath

	req.Params["pipeline_config"] = "/data/premodel/code/ApulisVision/panel.json"

	//baseconfig待定，
	//req.Params["config"] = req.CodePath

	training := models.Training{
		Id:              req.JobId,
		Name:            req.Name,
		Engine:          req.Engine,
		CodePath:        req.CodePath,
		StartupFile:     req.StartupFile,
		OutputPath:      req.OutputPath,
		DatasetPath:     req.DatasetPath,
		Params:          req.Params,
		Desc:            req.Description,
		NumPs:           req.NumPs,
		NumPsWorker:     req.NumPsWorker,
		DeviceType:      req.DeviceType,
		DeviceNum:       req.DeviceNum,
		JobTrainingType: req.JobTrainingType,
	}
	//启动训练作业
	jobId, err := CreateTraining(username, training)
	if err != nil {
		return req, err
	}
	req.JobId = jobId
	//nodes和edges只用存储然后传给前端
	req.Params["nodes"] = string(nodesBytes)
	req.Params["edges"] = string(edgesBytes)

	return req, nil
}
