package services

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/apulis/AIArtsBackend/models"
	"os"
	"strings"
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
	VisualPath  string `json:"visualPath"`

	//用于可视化建模平台直接启动训练任务
	JobTrainingType string          `json:"jobTrainingType"`
	NumPs           int             `json:"numPs"`
	NumPsWorker     int             `json:"numPsWorker"`
	DeviceType      string          `json:"deviceType"`
	DeviceNum       int             `json:"deviceNum"`
	Nodes           []AvisualisNode `json:"nodes"`
	Edges           []AvisualisEdge `json:"edges"`
	Panel           string `json:"panel"`

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
	use string, size int64, dataFormat, datasetName, datasetPath string, params map[string]string, engine, precision, outputPath, startupFile, deviceType, visualPath string, deviceNum int) error {

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
		VisualPath:  visualPath,

	}
	//只能创建Avisualis模型
	if strings.HasPrefix(use, `Avisualis`) {
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
			DeviceType:  deviceType,
			DeviceNum:   deviceNum,
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
			modelset.VisualPath = job.Params["visualPath"]
		} else {
			return fmt.Errorf("the job id is invaild")
		}
	} else {
		modelset.CodePath = codePath
	}
	return models.CreateModelset(modelset)
}

func UpdateModelset(id int, name, description, version, jobId, codePath, paramPath,
	use string, size int64, dataFormat, datasetName, datasetPath string, params map[string]string, engine, precision, outputPath, startupFile, visualPath string) error {
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
	modelset.Engine = engine
	modelset.JobId = jobId
	modelset.ParamPath = paramPath
	modelset.Use = use
	modelset.Size = size
	modelset.DataFormat = dataFormat
	modelset.DatasetName = datasetName
	modelset.Precision = precision
	modelset.VisualPath = visualPath
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
func GeneratePanel(modelset models.Modelset, username string) (models.Modelset, error) {
	//获取可用的数据集
	datasets, total, err := ListDatasets(1, 999, "created_at", "desc", "", "all", true, username)
	_, _, datasets = AppendAnnoDataset(datasets, total, 1, 999, "created_at", "desc")
	if err != nil {
		return modelset, err
	}
	//分类获取panel
	paramString, _ := json.Marshal(modelset.Params)
	paramJson, err := gabs.ParseJSON(paramString)
	if err != nil {
		return modelset, err
	}
	//去掉\保证可以解析
	panelString := paramJson.S("panel").String()
	panelString = strings.TrimPrefix(panelString, `"`)
	panelString = strings.TrimSuffix(panelString, `"`)
	panelString = strings.ReplaceAll(panelString, `\`, ``)
	panelJson, err := gabs.ParseJSON([]byte(panelString))

	//生成input节点，并插入数据集
	input := gabs.New()
	for _, dataset := range datasets {
		config := gabs.New()
		children := gabs.New()
		_, err = config.Set("data_path", "key")
		_, err = config.Set("disabled", "type")
		_, err = config.Set(dataset.Path, "value")
		_ = children.ArrayAppend(config, dataset.Name)
		_ = input.ArrayAppend(children, "children")
	}
	_, err = input.Set("Input", "name")
	_, err = panelJson.SetIndex(input, 0)

	//加入修改好后的panel
	_, err = paramJson.Set(panelJson.String(), "panel")
	var params models.ParamsItem
	err = json.Unmarshal([]byte(paramJson.String()), &params)
	modelset.Params = &params

	if err != nil {
		return modelset, err
	}
	return modelset, nil
}

func CreateAvisualisTraining(req CreateModelsetReq, username string) (string, error) {
	//存储节点json
	nodesBytes, _ := json.Marshal(req.Nodes)

	//去掉nodes没用的节点并存入json
	pipelineConfigPath, err := GetModelTempPath(FILETYPE_JSON)
	f, err := os.OpenFile(pipelineConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 777)
	if err != nil {
		fmt.Println(pipelineConfigPath + "failed to created")
	}
	_, err = f.Write(nodesBytes)
	defer f.Close()

	//把数据传入params后端算法只需要pipeline_config
	trainParams := make(map[string]string)
	trainParams["pipeline_config"] = pipelineConfigPath
	//trainParams["pipeline_config"] = "/data/premodel/code/ApulisVision/panel.json"

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
		Params:          trainParams,
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
		return jobId, err
	}
	//nodes和edges只用存储然后传给前端
	return jobId, nil
}
