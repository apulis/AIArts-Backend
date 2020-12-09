package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	urllib "net/url"
	"strconv"
)

func GetAllTraining(userName string, page, size int, jobStatus,jobGroup, searchWord, orderBy, order string) ([]*models.Training, int, int, error) {
	//把传输过来的searchword空格改为%20urlencode
	url := fmt.Sprintf(`%s/ListJobsV3?userName=%s&jobOwner=%s&vcName=%s&jobType=%s&pageNum=%d&pageSize=%d&jobStatus=%s&searchWord=%s&orderBy=%s&order=%s&jobGroup=%s`,
		configs.Config.DltsUrl, userName, userName, models.DefaultVcName,
		models.JobTypeArtsTraining,
		page, size, jobStatus, urllib.PathEscape(searchWord),
		orderBy, order,jobGroup)

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all training err[%+v]", err)
		return nil, 0, 0, err
	}

	trainings := make([]*models.Training, 0)
	for _, v := range jobList.AllJobs {
		experiment_id,_ :=strconv.Atoi(v.JobParams.JobGroup)
		trainings = append(trainings, &models.Training{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			DeviceType:  v.JobParams.GpuType,
			CodePath:    v.JobParams.CodePath,
			DeviceNum:   v.JobParams.Resourcegpu,
			StartupFile: v.JobParams.StartupFile,
			OutputPath:  v.JobParams.OutputPath,
			DatasetPath: v.JobParams.DatasetPath,
			Params:      nil,
			Status:      v.JobStatus,
			CreateTime:  v.JobTime,
			Desc:        v.JobParams.Desc,
			ExperimentID: uint64(experiment_id),
			Track:       v.JobParams.Track,
		})
	}

	totalJobs := jobList.Meta.TotalJobs
	totalPages := totalJobs / size

	if (totalJobs % size) != 0 {
		totalPages += 1
	}

	return trainings, totalJobs, totalPages, nil
}

func CreateTraining(userName string, training models.Training) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})
	params["userName"] = userName
	params["jobName"] = training.Name
	params["jobType"] = models.JobTypeArtsTraining

	params["image"] = ConvertImage(training.Engine)
	params["gpuType"] = training.DeviceType
	params["resourcegpu"] = training.DeviceNum
	params["DeviceNum"] = training.DeviceNum
	params["cmd"] = "" // use StartupFile, params instead

	if configs.Config.InteractiveModeJob {
		params["cmd"] = "sleep infinity" // use StartupFile, params instead
	} else {
		fileType, err := CheckStartFileType(training.StartupFile)
		if fileType == FILETYPE_PYTHON {
			params["cmd"] = "python " + training.StartupFile
		} else if fileType == FILETYPE_SHELL {
			params["cmd"] = "bash " + training.StartupFile
		}
		if err != nil {
			fmt.Printf("startupfile is invalid[%+v]\n", err)
			return "", err
		}
		for k, v := range training.Params {
			if k == "sudo" {
				//添加sudo权限
				params["cmd"] = "sudo " + v + " " + params["cmd"].(string)
			} else if len(k) > 0 && len(v) > 0 {
				params["cmd"] = params["cmd"].(string) + " --" + k + " " + v + " "
			}

		}
		if len(training.DatasetPath) > 0 {
			params["cmd"] = params["cmd"].(string) + " --data_path " + training.DatasetPath
		}
		//params中加入visualpath
		if len(training.VisualPath) > 0 {
			training.Params["visualPath"] = training.VisualPath
		}
		if len(training.OutputPath) > 0 {
			params["cmd"] = params["cmd"].(string) + " --output_path " + training.OutputPath
		}
	}

	params["startupFile"] = training.StartupFile
	params["datasetPath"] = training.DatasetPath
	params["codePath"] = training.CodePath
	params["outputPath"] = training.OutputPath
	params["scriptParams"] = training.Params
	params["desc"] = training.Desc

	params["containerUserId"] = 0
	params["jobtrainingtype"] = training.JobTrainingType // "RegularJob"
	params["preemptionAllowed"] = false
	params["workPath"] = "./"

	params["enableworkpath"] = true
	params["enabledatapath"] = true
	params["enablejobpath"] = true
	params["jobPath"] = "./"

	params["hostNetwork"] = false
	params["isPrivileged"] = false
	params["interactivePorts"] = false

	params["numpsworker"] = training.NumPs
	params["numps"] = training.NumPsWorker

	params["vcName"] = models.DefaultVcName
	params["team"] = models.DefaultVcName

	err := checkCreateTrainingWithMlflow(training.ExperimentID,&params)
	if err != nil{
		return "", err
	}

	id := &models.JobId{}
	err = DoRequest(url, "POST", nil, params, id)

	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return "", err
	}

	return id.Id, nil
}

func DeleteTraining(userName, id string) error {

	url := fmt.Sprintf("%s/KillJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string]interface{})

	job := &models.Job{}
	err := DoRequest(url, "GET", nil, params, job)

	if err != nil {
		fmt.Printf("delete training err[%+v]\n", err)
		return err
	}

	return nil
}

func GetTraining(userName, id string) (*models.Training, error) {

	url := fmt.Sprintf("%s/GetJobDetailV2?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string]interface{})

	job := &models.Job{}
	training := &models.Training{}

	err := DoRequest(url, "GET", nil, params, job)
	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return nil, err
	}

	training.Id = job.JobId
	training.Name = job.JobName
	training.Engine = UnConvertImage(job.JobParams.Image)
	training.DeviceNum = job.JobParams.Resourcegpu
	training.DeviceType = job.JobParams.GpuType
	training.Status = job.JobStatus
	training.CreateTime = job.JobTime
	training.JobTrainingType = job.JobParams.Jobtrainingtype

	training.Params = nil
	training.CodePath = job.JobParams.CodePath
	training.StartupFile = job.JobParams.StartupFile
	training.OutputPath = job.JobParams.OutputPath
	training.DatasetPath = job.JobParams.DatasetPath
	training.Status = job.JobStatus
	training.Desc = job.JobParams.Desc
	training.Params = job.JobParams.ScriptParams

	training.ExperimentID,_ = strconv.ParseUint(job.JobParams.JobGroup,0,0)
	training.Track        = job.JobParams.Track

	return training, nil
}

func GetTrainingLog(userName, id string, pageNum int) (*models.JobLog, error) {

	url := fmt.Sprintf("%s/GetJobLog?userName=%s&jobId=%s&page=%d", configs.Config.DltsUrl, userName, id, pageNum)
	jobLog := &models.JobLog{}

	jobLogFromDlts := &struct {
		Cursor  string `json:"cursor,omitempty"`
		Log     string `json:"log,omitempty"`
		MaxPage int    `json:"max_page"`
	}{}

	err := DoRequest(url, "GET", nil, nil, jobLogFromDlts)
	if err != nil {
		fmt.Printf("create training err[%+v]\n", err)
		return nil, err
	}

	jobLog.Cursor = jobLogFromDlts.Cursor
	jobLog.Log = jobLogFromDlts.Log
	jobLog.MaxPage = jobLogFromDlts.MaxPage

	return jobLog, nil
}

//@add:  support for intergration with mlflow
func checkCreateTrainingWithMlflow( experimentId uint64, params*map[string]interface{})error{
   if experimentId == 0 {
	   (*params)["jobGroup"] = ""
   	   return nil
   }
	(*params)["jobGroup"] = strconv.Itoa(int(experimentId))
	if !configs.Config.EnableTrack {
		return nil
	}
	resp , err := StartMlflowRun(experimentId,(*params)["userName"].(string),(*params)["jobName"].(string))
	if err != nil{
		return err
	}
	run := resp.(*MlflowRun)
	(*params)["jobId"]=run.Info.RunId
	(*params)["track"]=1
	return nil
}
