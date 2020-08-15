package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Evaluation struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Engine      string            `json:"engine"`
	DeviceType  string            `json:"deviceType"`
	DeviceNum   int               `json:"deviceNum"`
	CodePath    string            `json:"codePath"`
	StartupFile string            `json:"startupFile"`
	OutputPath  string            `json:"outputPath"`
	DatasetPath string            `json:"datasetPath"`
	Params      map[string]string `json:"params"`
	ParamPath   string            `json:"paramPath"`
	DatasetName string            `json:"datasetName"`
	Status      string            `json:"status"`
	CreateTime  string            `json:"createTime"`
}

func CreateEvaluation(userName string, evaluation Evaluation) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	params := make(map[string]interface{})

	params["userName"] = userName
	params["jobName"] = evaluation.Name
	params["jobType"] = models.JobTypeArtsEvaluation

	params["image"] = ConvertImage(evaluation.Engine)
	params["gpuType"] = evaluation.DeviceType
	params["resourcegpu"] = evaluation.DeviceNum
	params["DeviceNum"] = evaluation.DeviceNum
	params["cmd"] = "" // use StartupFile, params instead
	params["cmd"] = "python " + evaluation.StartupFile
	if len(evaluation.DatasetPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --data_path " + evaluation.DatasetPath
	}
	if len(evaluation.OutputPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --output_path " + evaluation.OutputPath
	}
	if len(evaluation.ParamPath) > 0 {
		params["cmd"] = params["cmd"].(string) + " --checkpoint_path  " + evaluation.ParamPath
	}
	for k, v := range evaluation.Params {
		if len(k) > 0 && len(v) > 0 {
			params["cmd"] = params["cmd"].(string) + " --" + k + " " + v + " "
		}
	}

	logger.Info(fmt.Sprintf("evaluation : %s", params["cmd"]))
	params["startupFile"] = evaluation.StartupFile
	params["datasetPath"] = evaluation.DatasetPath
	params["codePath"] = evaluation.CodePath
	params["outputPath"] = evaluation.OutputPath
	params["scriptParams"] = evaluation.Params
	params["desc"] = fmt.Sprintf("%s^%s", evaluation.DatasetName, evaluation.ParamPath)
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
	params["vcName"] = models.DefaultVcName
	params["team"] = models.DefaultVcName
	id := &models.JobId{}
	err := DoRequest(url, "POST", nil, params, id)
	if err != nil {
		fmt.Printf("create evaluation err[%+v]\n", err)
		return "", err
	}
	return id.Id, nil
}

func GetEvaluations(userName string, page, size int, jobStatus, searchWord, orderBy, order string) ([]*Evaluation, int, int, error) {
	url := fmt.Sprintf(`%s/ListJobsV3?userName=%s&jobOwner=%s&vcName=%s&jobType=%s&pageNum=%d&pageSize=%d&jobStatus=%s&searchWord=%s&orderBy=%s&order=%s`,
		configs.Config.DltsUrl, userName, userName, models.DefaultVcName,
		models.JobTypeArtsEvaluation,
		page, size, jobStatus, url.PathEscape(searchWord),
		orderBy, order)

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all evaluation err[%+v]", err)
		return nil, 0, 0, err
	}

	evaluations := make([]*Evaluation, 0)
	for _, v := range jobList.AllJobs {
		evaluations = append(evaluations, &Evaluation{
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
			DatasetName: v.JobParams.Desc,
		})
	}

	totalJobs := jobList.Meta.TotalJobs
	totalPages := totalJobs / page

	if (totalJobs % page) != 0 {
		totalPages += 1
	}

	return evaluations, totalJobs, totalPages, nil
}

func DeleteEvaluation(userName, id string) error {
	url := fmt.Sprintf("%s/KillJob?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string]interface{})
	job := &models.Job{}
	err := DoRequest(url, "GET", nil, params, job)
	if err != nil {
		fmt.Printf("delete evaluation err[%+v]\n", err)
		return err
	}

	return nil
}

func GetEvaluation(userName, id string) (*Evaluation, error) {
	url := fmt.Sprintf("%s/GetJobDetailV2?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	params := make(map[string]interface{})
	job := &models.Job{}
	evaluation := &Evaluation{}

	err := DoRequest(url, "GET", nil, params, job)
	if err != nil {
		fmt.Printf("create evaluation err[%+v]\n", err)
		return nil, err
	}
	evaluation.Id = job.JobId
	evaluation.Name = job.JobName
	evaluation.Engine = job.JobParams.Image
	evaluation.DeviceNum = job.JobParams.Resourcegpu
	evaluation.DeviceType = job.JobParams.GpuType
	evaluation.Status = job.JobStatus
	evaluation.CreateTime = job.JobTime
	evaluation.CodePath = job.JobParams.CodePath
	evaluation.StartupFile = job.JobParams.StartupFile
	evaluation.OutputPath = job.JobParams.OutputPath
	evaluation.DatasetPath = job.JobParams.DatasetPath
	//解析desc为数据集名称^模型文件名称
	descSplit := strings.Split(job.JobParams.Desc, "^")
	if len(descSplit) > 1 {
		datasetName := descSplit[0]
		evaluation.DatasetName = datasetName
		//workpath为评估参数文件路径
		paramPath := descSplit[1]
		evaluation.DatasetName = datasetName
		evaluation.ParamPath = paramPath
	}
	evaluation.Params = job.JobParams.ScriptParams
	return evaluation, nil
}

func GetEvaluationLog(userName, id string) (*models.JobLog, error) {
	url := fmt.Sprintf("%s/GetJobLog?userName=%s&jobId=%s", configs.Config.DltsUrl, userName, id)
	jobLog := &models.JobLog{}
	err := DoRequest(url, "GET", nil, nil, jobLog)
	if err != nil {
		fmt.Printf("create evaluation err[%+v]\n", err)
		return nil, err
	}
	return jobLog, nil
}

func GetRegexpLog(log string) (map[string]string,map[string]string) {
	acc_reg, _ := regexp.Compile("Accuracy\\[(.*?)\\]")
	recall_5_reg, _ := regexp.Compile("Recall_5\\[(.*?)\\]")
	recall_reg, _ := regexp.Compile("Recall\\[(.*?)\\]")
	precision_reg, _ := regexp.Compile("Precision\\[(.*?)\\]")
	indicator := map[string]string{}
	confusion:= map[string]string{}
	if len(recall_reg.FindStringSubmatch(log)) > 1 {
		indicator["Recall"] = recall_reg.FindStringSubmatch(log)[1]
	}
	if len(recall_5_reg.FindStringSubmatch(log)) > 1 {
		indicator["Recall_5"] = recall_5_reg.FindStringSubmatch(log)[1]
	}
	if len(acc_reg.FindStringSubmatch(log)) > 1 {
		indicator["Accuracy"] = acc_reg.FindStringSubmatch(log)[1]
	}
	if len(precision_reg.FindStringSubmatch(log)) > 1 {
		indicator["Precision"] = precision_reg.FindStringSubmatch(log)[1]
	}
	//二分类混淆矩阵
	VALUE_reg, _ := regexp.Compile("\\[(.+)\\]")
	TP_reg, _ := regexp.Compile("TP: (.+?])")
	if len(TP_reg.FindStringSubmatch(log)) > 1 {
		TP_string := TP_reg.FindStringSubmatch(log)[1]
		label_x1_reg, _ := regexp.Compile("(.*?),")
		label_y1_reg, _ := regexp.Compile(" (.+?)\\[")
		confusion["x1"] = label_x1_reg.FindStringSubmatch(TP_string)[1]
		confusion["y1"] = label_y1_reg.FindStringSubmatch(TP_string)[1]
		confusion["TP"] = VALUE_reg.FindStringSubmatch(TP_string)[1]

		TN_reg, _ := regexp.Compile("TN: (.+?])")
		if len(TN_reg.FindStringSubmatch(log)) > 1 {
			TN_string := TN_reg.FindStringSubmatch(log)[1]
			label_x2_reg, _ := regexp.Compile("(.*?),")
			label_y2_reg, _ := regexp.Compile(" (.+?)\\[")
			confusion["x2"] = label_x2_reg.FindStringSubmatch(TN_string)[1]
			confusion["y2"] = label_y2_reg.FindStringSubmatch(TN_string)[1]
			confusion["TN"] = VALUE_reg.FindStringSubmatch(TN_string)[1]
		}
		FN_reg, _ := regexp.Compile("FN: (.+?])")
		if len(FN_reg.FindStringSubmatch(log)) > 1 {
			FN_string := FN_reg.FindStringSubmatch(log)[1]
			confusion["FN"] = VALUE_reg.FindStringSubmatch(FN_string)[1]
		}
		FP_reg, _ := regexp.Compile("FP: (.+?])")
		if len(FP_reg.FindStringSubmatch(log)) > 1 {
			FP_string := FP_reg.FindStringSubmatch(log)[1]
			confusion["FP"] = VALUE_reg.FindStringSubmatch(FP_string)[1]
		}
		TP, _ := strconv.ParseFloat(confusion["TP"], 32)
		FN, _ := strconv.ParseFloat(confusion["FN"], 32)
		FP, _ := strconv.ParseFloat(confusion["FP"], 32)
		TN, _ := strconv.ParseFloat(confusion["TN"], 32)
		confusion["Recall1"] = strconv.FormatFloat(TP/(TP+FN), 'f', -1, 32)
		confusion["Recall2"] = strconv.FormatFloat(TN/(FP+TN), 'f', -1, 32)
		confusion["Precision1"] = strconv.FormatFloat(TP/(FP+TP), 'f', -1, 32)
		confusion["Precision2"] = strconv.FormatFloat(TN/(FN+TN), 'f', -1, 32)

		Accuracy_reg, _ := regexp.Compile("Accuracy\\[(.+)\\]")
		if len(Accuracy_reg.FindStringSubmatch(log)) > 1 {
			indicator["Accuracy"] = Accuracy_reg.FindStringSubmatch(log)[1]
		}
		Recall_reg, _ := regexp.Compile("Recall\\[(.+)\\]")
		if len(Recall_reg.FindStringSubmatch(log)) > 1 {
			indicator["Recall"] = Recall_reg.FindStringSubmatch(log)[1]
		}
		Precision_reg, _ := regexp.Compile("Precision\\[(.+)\\]")
		if len(Precision_reg.FindStringSubmatch(log)) > 1 {
			indicator["Precision"] = Precision_reg.FindStringSubmatch(log)[1]
		}
		F1_score_reg, _ := regexp.Compile("F1_score\\[(.+)\\]")
		if len(F1_score_reg.FindStringSubmatch(log)) > 1 {
			indicator["F1_score"] = F1_score_reg.FindStringSubmatch(log)[1]
		}
		Auc_ROC_reg, _ := regexp.Compile("Auc_ROC\\[(.+)\\]")
		if len(Auc_ROC_reg.FindStringSubmatch(log)) > 1 {
			indicator["Auc_ROC"] = Auc_ROC_reg.FindStringSubmatch(log)[1]
		}
		Auc_PR_reg, _ := regexp.Compile("Auc_PR\\[(.+)\\]")
		if len(Auc_PR_reg.FindStringSubmatch(log)) > 1 {
			indicator["Auc_PR"] = Auc_PR_reg.FindStringSubmatch(log)[1]
		}
	}

	//目标检测
	mAP_reg, _ := regexp.Compile("mAP@0.5IOU: (.*)")
	localization_loss_reg, _ := regexp.Compile("localization_loss: (.*)")
	classification_loss_reg, _ := regexp.Compile("classification_loss: (.*)")
	regularization_loss_reg, _ := regexp.Compile("regularization_loss: (.*)")
	total_loss_reg, _ := regexp.Compile("total_loss: (.*)")

	if len(mAP_reg.FindStringSubmatch(log)) > 1 {
		indicator["mAP"] = mAP_reg.FindStringSubmatch(log)[1]
		if len(localization_loss_reg.FindStringSubmatch(log)) > 1 {
			indicator["Localization_Loss"] = localization_loss_reg.FindStringSubmatch(log)[1]
		}
		if len(classification_loss_reg.FindStringSubmatch(log)) > 1 {
			indicator["Classification_Loss"] = classification_loss_reg.FindStringSubmatch(log)[1]
		}
		if len(regularization_loss_reg.FindStringSubmatch(log)) > 1 {
			indicator["Regularization_Loss"] = regularization_loss_reg.FindStringSubmatch(log)[1]
		}
		if len(total_loss_reg.FindStringSubmatch(log)) > 1 {
			indicator["Total_Loss"] = total_loss_reg.FindStringSubmatch(log)[1]
		}
	}

	//pytorch

	acc_pytorch_reg, _ := regexp.Compile("Accuracy:(.+)")
	if len(acc_pytorch_reg.FindStringSubmatch(log)) > 1 {
		indicator["Accuracy"] = acc_pytorch_reg.FindStringSubmatch(log)[1]
	}
	avg_loss_pytorch_reg, _ := regexp.Compile("Average loss: (.+?),")
	if len(avg_loss_pytorch_reg.FindStringSubmatch(log)) > 1 {
		indicator["Average_Loss"] = avg_loss_pytorch_reg.FindStringSubmatch(log)[1]
	}
	//mxnet
	acc_mxnet_reg, _ := regexp.Compile("accuracy=(.+)")
	if len(acc_mxnet_reg.FindStringSubmatch(log)) > 1 {
		indicator["Accuracy"] = acc_mxnet_reg.FindStringSubmatch(log)[1]
	}


	return indicator,confusion
}
