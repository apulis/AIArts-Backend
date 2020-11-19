package services

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

const (
	MODELSET_STATUS_NORMAL   = "normal"
	MODELSET_STATUS_DELETING = "deleting"
)

type AvisualisEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}
type AvisualisNode struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	TreeIdx int         `json:"treeIdx"`
	Config  interface{} `json:"config"`
	ComboId string      `json:"comboId"`
}
type AvisualisCombos struct {
	ID           string      `json:"id"`
	Label        string      `json:"label"`
	ParentId     string      `json:"parentId"`
	AnchorPoints interface{} `json:"anchorPoints"`
}

func ListModelSets(username string, req models.LsModelsetsReq) ([]models.Modelset, int, error) {

	offset := req.PageSize * (req.PageNum - 1)
	limit := req.PageSize

	return models.ListModelSets(username, offset, limit, req.OrderBy,
		req.Order, req.IsAdvance, req.Name, req.Status, req.Use)
}

func CreateModelset(username, version string, req models.CreateModelsetReq) error {

	//只能创建非预置模型
	modelset := models.Modelset{
		Name:        req.Name,
		Description: req.Description,
		Creator:     username,
		Version:     version,
		JobId:       req.JobId,
		Status:      MODELSET_STATUS_NORMAL,
		IsAdvance:   req.IsAdvance,
		ParamPath:   req.ParamPath,
		VisualPath:  req.VisualPath,
		//VCName:      req.VCName,
	}

	//只能创建Avisualis模型
	if strings.HasPrefix(req.Use, `Avisualis`) {
		var paramItem models.ParamsItem
		paramItem = req.Params
		modelset = models.Modelset{
			Name:        req.Name,
			Description: req.Description,
			Creator:     username,
			Version:     version,
			JobId:       req.JobId,
			Status:      MODELSET_STATUS_NORMAL,
			IsAdvance:   req.IsAdvance,
			ParamPath:   req.ParamPath,
			Use:         req.Use,
			Size:        req.Size,
			DataFormat:  req.DataFormat,
			DatasetPath: req.DatasetPath,
			DatasetName: req.DatasetName,
			Params:      &paramItem,
			Engine:      req.Engine,
			Precision:   req.Precision,
			OutputPath:  req.OutputPath,
			StartupFile: req.StartupFile,
			DeviceType:  req.DeviceType,
			DeviceNum:   req.DeviceNum,
			//VCName:      req.VCName,
		}
	}
	//获取训练作业输出模型的类型 job
	if req.CodePath == "" {
		job, _ := GetTraining(username, req.JobId)

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
			//modelset.VCName = job.VCName
		} else {
			return fmt.Errorf("the job id is invaild")
		}
	} else {
		modelset.CodePath = req.CodePath
	}

	return models.CreateModelset(modelset)
}

func UpdateModelset(ID int, version string, req models.CreateModelsetReq)  error {

	modelset, err := models.GetModelsetById(ID)
	if err != nil {
		return err
	}

	var paramItem models.ParamsItem

	paramItem = req.Params
	modelset.Description = req.Description
	modelset.OutputPath = req.OutputPath
	modelset.CodePath = req.CodePath
	modelset.DatasetPath = req.DatasetPath
	modelset.StartupFile = req.StartupFile
	modelset.Params = &paramItem
	modelset.Name = req.Name
	modelset.Version = version
	modelset.Engine = req.Engine
	modelset.JobId = req.JobId
	modelset.ParamPath = req.ParamPath
	modelset.Use = req.Use
	modelset.Size = req.Size
	modelset.DataFormat = req.DataFormat
	modelset.DatasetName = req.DatasetName
	modelset.Precision = req.Precision
	modelset.VisualPath = req.VisualPath
	//modelset.VCName = req.VCName

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
	fmt.Println(panelString)
	panelString = strings.TrimPrefix(panelString, `"`)
	panelString = strings.TrimSuffix(panelString, `"`)
	panelString = strings.ReplaceAll(panelString, `\`, ``)
	panelJson, err := gabs.ParseJSON([]byte(panelString))

	//生成input节点，并插入数据集
	input := gabs.New()
	for _, dataset := range datasets {
		config := gabs.New()
		item := gabs.New()
		_, err = config.Set("data_path", "key")
		_, err = config.Set("disabled", "type")
		_, err = config.Set(dataset.Path, "value")

		_ = item.ArrayAppend(config, "config")
		_, err = item.Set(dataset.Name, "name")
		_ = input.ArrayAppend(item, "children")
	}
	_, err = input.Set("Input", "name")
	_, err = panelJson.SetIndex(input, 0)

	//加入修改好后的panel
	_, err = paramJson.Set(panelJson.String(), "panel")
	fmt.Println(panelJson.String())
	var params models.ParamsItem
	err = json.Unmarshal([]byte(paramJson.String()), &params)
	modelset.Params = &params

	if err != nil {
		return modelset, err
	}
	return modelset, nil
}

func CreateAvisualisTraining(c *gin.Context, req models.CreateModelsetReq, username string) (models.CreateModelsetReq, error) {

	//存储节点json
	nodesBytes, _ :=req.Params["nodes"]

	//去掉nodes没用的节点并存入json
	pipelineConfigPath, err := GetModelTempPath(FILETYPE_JSON)
	f, err := os.OpenFile(pipelineConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println(pipelineConfigPath + "failed to created")
	}

	_, err = f.Write([]byte(nodesBytes))
	defer f.Close()

	training := models.Training{
		Id:          req.JobId,
		Name:        req.Name,
		Engine:      req.Engine,
		CodePath:    req.CodePath,
		StartupFile: req.StartupFile,
		OutputPath:  req.OutputPath,
		DatasetPath: req.DatasetPath,
		Params: map[string]string{
			"pipeline_config": pipelineConfigPath,
		},
		Desc:            req.Description,
		NumPs:           req.NumPs,
		NumPsWorker:     req.NumPsWorker,
		DeviceType:      req.DeviceType,
		DeviceNum:       req.DeviceNum,
		JobTrainingType: req.JobTrainingType,
		VCName:          req.VCName,
	}

	//启动训练作业
	jobId, err := CreateTraining(c, username, training)
	//panel不用变
	req.JobId = jobId
	req.Params["pipeline_config"] = pipelineConfigPath
	if err != nil {
		return req, err
	}

	//nodes和edges,combos只用存储然后传给前端
	return req, nil
}
