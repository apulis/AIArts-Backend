package services

import (
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
	"time"
)

type CommQueryParams struct {
	Offset  uint
	PageNum uint     `form:"pageNum" `
	Limit   uint     `form:"pageSize"`
	OrderBy string   `form:"orderBy" `
	Order   string   `form:"order" `
	Name    string   `form:"searchWord"`
	All     uint     `form:"all"`
}

var ai_arts_experiment_name_prefix="ai_arts_"

func (s *CommQueryParams) SetQueryParams(c *gin.Context) {
	 c.ShouldBind(s)
	 if s.Limit > 100 {
	 	s.Limit=100
	 } else if s.Limit < 10 {
	 	s.Limit=10
	}
	if s.PageNum >= 1 {
		s.PageNum--;
	}
	 s.Offset= s.PageNum * s.Limit
     if s.Order != "" && s.Order != "asc" && s.Order != "desc" {
     	s.Order=""
	 }
	 //@todo: check for order by ???
	 //s.OrderBy=""

}

func ListExpProjects(queryParam CommQueryParams, username string) ([]models.ExpProject, uint, error) {
	return models.ListAllExpProjects(queryParam.Offset, queryParam.Limit, queryParam.All,
		queryParam.Name,queryParam.OrderBy, queryParam.Order)
}

func CreateExpProject(project* models.ExpProject) error {
	return models.CreateExpProject(project)
}

func  QueryExpProject(id uint64,project*models.ExpProject)error{
	return models.QueryExpProject(uint(id),project)
}

func  UpdateExpProject(id uint64,new_name string,project*models.RequestUpdates)error{
	if project == nil{
		return models.RenameExpProject(uint(id),new_name)
	}else{
		if len(new_name) != 0{
			(*project)["name"]=new_name
		}
		return models.UpdateExpProject(uint(id),project)
	}
}

func MarkExpProject(id uint64, hide bool) error {
	if hide {
		return models.HideExpProject(uint(id))
	}else{
		return models.RestoreExpProject(uint(id))
	}
}

func ListExperiments(queryParam CommQueryParams, projectID uint64) ([]models.Experiment, uint, error) {
	return models.ListAllExperiments(uint(projectID), queryParam.Offset, queryParam.Limit, queryParam.All,
		queryParam.Name, queryParam.OrderBy, queryParam.Order)
}

func CreateExperiment(experiment*models.Experiment)error{
	return models.CreateExperiment(experiment)
}

func  UpdateExperiment(id uint64,new_name string,experiment*models.RequestUpdates)error{
	if experiment == nil{
		return models.RenameExperiment(uint(id),new_name)
	}else{
		if len(new_name) != 0{
			(*experiment)["name"]=new_name
		}
		//@todo: how to store user define data ???
		return models.UpdateExperiment(uint(id),experiment)
	}
}
func MarkExperiment(id uint64, hide bool) error {
	if hide {
		return models.HideExperiment(uint(id))
	}else{
		return models.RestoreExperiment(uint(id))
	}
}

func  QueryExperiment(id uint64,experiment*models.Experiment)error{
	err := models.QueryExperiment(uint(id),experiment)
	if err != nil{
		return err
	}
	/*trackExperiment,err := getMlflowExperiment(id)
	if err != nil{
		return err
	}
	if trackExperiment != nil{
 		experiment.TrackId=trackExperiment.ExperimentID
	}*/
    return nil
}
func JumpExperimentView(id uint64) (map[string]interface{},error){
	 var experiment models.Experiment

	 err := models.QueryExperiment(uint(id),&experiment)
	 if err != nil{
	 	return nil,err
	 }
	experimentPtr,err := getMlflowExperiment(id)
	if err != nil {
		return nil,err
	}
	if experimentPtr == nil {//create if not exists
		id,err := createMlflowExperiment(id)
		if err != nil{
			return nil,err
		}
		experimentPtr=&MlflowExperiment{ExperimentID: id}
	}
	result := make(map[string]interface{})
	result["id"]=id
	result["name"]=experiment.Name
	result["projectId"]=experiment.ProjectID
	result["projectName"]=experiment.ProjectName
	result["trackId"]=experimentPtr.ExperimentID
	return result,nil
}

func getMlflowExperimentName(experiment_id uint64) string{
	 return fmt.Sprintf("%s%d",ai_arts_experiment_name_prefix,experiment_id)
}


type MlflowMetric struct{
     Key       string   `json:"key"`
     Value     float64  `json:"value"`
	 Timestamp int64    `json:"timestamp,string"`
     Step      int64    `json:"step,string"`
}
type RunTag struct{
	 Key   string       `json:"key"`
	 Value string       `json:"value"`
}

type MlflowRunInfo struct{
	RunId        string   `json:"run_id"`
    ExperimentId string   `json:"experiment_id"`
	Status       string    `json:"status"`
	StartTime    int64     `json:"start_time,string"`
	EndTime      int64     `json:"end_time,string"`
	ArtifactUri    string   `json:"artifact_uri"`
	LifecycleStage string   `json:"lifecycle_stage"`
}
type MlflowRunData struct{
	Metrics []MlflowMetric   `json:"metrics"`
	Params  []RunTag         `json:"params"`
	Tags    []RunTag         `json:"tags"`
}

type MlflowRun struct{
	Info  MlflowRunInfo    `json:"info"`
	Data  MlflowRunData    `json:"data"`
}

type MlflowRunResp struct{
	ErrorCode string  `json:"error_code"`
	Message   string  `json:"message"`
	Run MlflowRun     `json:"run"`
}

type MlflowExperiment struct{
	ExperimentID  string       `json:"experiment_id"`
	Name          string       `json:"name"`
	ArtifactLocation string    `json:"artifact_location"`
	LifecycleStage   string    `json:"lifecycle_stage"`
	LastUpdateTime   int64     `json:"last_update_time,string"`
	CreationTime     int64     `json:"creation_time,string"`
	Tags             []RunTag  `json:"tags"`
}

type MlflowInfoExperimentResp struct{
	ErrorCode  string   `json:"error_code"`
	Message    string   `json:"message"`
	Experiment MlflowExperiment `json:"experiment"`
}

func getMlflowExperiment(experimentID uint64) (*MlflowExperiment,error){
	url := fmt.Sprintf("%s/experiments/get-by-name?experiment_name=%s",configs.Config.TrackingUrl,getMlflowExperimentName(experimentID))
	experimentResp := MlflowInfoExperimentResp{}
    err := DoRequest2(url,"GET",nil,nil,&experimentResp)
    if err != nil {
		if experimentResp.ErrorCode == "RESOURCE_DOES_NOT_EXIST"{// supress not found error
			return nil,nil
		}
		return nil,err
	}
    return &experimentResp.Experiment,err
}
func createMlflowExperiment(experimentID uint64) (string,error){
	url := fmt.Sprintf("%s/experiments/create",configs.Config.TrackingUrl)
	Resp := make(map[string]interface{})
	type stCreateRequest struct{
		Name string `json:"name"`
	}
	err := DoRequest2(url,"POST",nil,stCreateRequest{
		Name:getMlflowExperimentName(experimentID),
	},&Resp)
	if err != nil{
		return "",err
	}
	id := Resp["experiment_id"].(string)
	if len(id) == 0{
		return "" , errors.New("Create mlflow experiment ID failed !")
	}
	return id,nil
}
func createMlflowRun(mlflow_exp_id ,user,job string)(*MlflowRun,error){
	url := fmt.Sprintf("%s/runs/create",configs.Config.TrackingUrl)
	Request := map[string]interface{}{
		"experiment_id": mlflow_exp_id,
		"start_time":    time.Now().UnixNano() / 1e6,
	}
	Request["tags"]=[]RunTag{
		         {Key:"mlflow.user",  Value:user},
		         {Key:"mlflow.runName",Value:job},
	    }
	runResp := MlflowRunResp{}
	err := DoRequest(url,"POST",nil,Request,&runResp)
	if err != nil{
		return nil,err
	}
	return &runResp.Run,err
}

//  sync with mlflow to ensure specified experiment exists and allocate a run id
func  StartMlflowRun(experimentID uint64,user,job string) (interface{},error){

	experimentPtr,err := getMlflowExperiment(experimentID)
	if err != nil {
         return nil,err
	}
	if experimentPtr == nil {//create if not exists
		id,err := createMlflowExperiment(experimentID)
		if err != nil{
			return nil,err
		}
		experimentPtr=&MlflowExperiment{ExperimentID: id}
	}
	run,err := createMlflowRun(experimentPtr.ExperimentID,user,job)
	return run,err
}
func QueryMlflowRun(runID string) (interface{},error){
	url := fmt.Sprintf("%s/runs/get?run_id=%s",configs.Config.TrackingUrl,runID)
	runResp := MlflowRunResp{}
	err := DoRequest2(url,"GET",nil,nil,&runResp)
	if err != nil {
		return nil,err
	}
	return &runResp.Run,nil
}
