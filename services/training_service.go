package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/models"
)

func GetAllTraining(userName string, page, size int) ([] *models.Training, int, int, error) {

	url := fmt.Sprintf("http://atlas02.sigsus.cn/apis/ListJobsV2?userName=%s&jobOwner=%s&num=%d&vcName=%s",
						userName, userName, 1000, "atlas")

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Print("request err: %+v", err)
		return nil, 0, 0, err
	}

	trainings := make([] *models.Training, 0)
	for _, v:= range jobList.RunningJobs {
		trainings = append(trainings, &models.Training{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			DeviceType:  v.JobParams.GpuType,
			DeviceNum:   v.JobParams.Resourcegpu,
			CodePath:    v.JobParams.DataPath,
			StartupFile: "",
			OutputPath:  v.JobParams.WorkPath,
			DatasetPath: v.JobParams.DataPath,
			Params:      nil,
			Desc:        "",
		})
	}

	for _, v:= range jobList.FinishedJobs {
		trainings = append(trainings, &models.Training{
			Id:          v.JobId,
			Name:        v.JobName,
			Engine:      v.JobParams.Image,
			DeviceType:  v.JobParams.GpuType,
			DeviceNum:   v.JobParams.Resourcegpu,
			CodePath:    v.JobParams.DataPath,
			StartupFile: "",
			OutputPath:  v.JobParams.WorkPath,
			DatasetPath: v.JobParams.DataPath,
			Params:      nil,
			Desc:        "",
		})
	}

	return trainings, len(trainings), 1, nil
}

func CreateTraining(training models.Training) (string, error) {
	return RandStringRunes(16), nil
}

func DeleteTraining(id string) error {
	return nil
}

func GetTraining(id string) error {
	return nil
}
