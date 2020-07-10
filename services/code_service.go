package services

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/apulis/AIArtsBackend/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetAllCodeset(page, size int) ([] *models.CodesetItem, int, int, error) {

	rand.Seed(time.Now().Unix())
	item := &models.CodesetItem{
		Name: RandStringRunes(16),
		Status: "started",
		Engine: "tf_1.15",
		CodePath: "/home/bifeng.peng/",
		CreateTime: time.Now().Unix(),
		Desc: "test test test",
	}

	url := fmt.Sprintf("http://atlas02.sigsus.cn/apis/ListJobsV2?userName=%s&jobOwner=%s&num=%d&vcName=%s",
							"yunxia.chu", "yunxia.chu", 10, "atlas")
	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Print("request err: %+v", err)
		return nil, 0, 0, err
	}

	codes := make([] *models.CodesetItem, 0)
	for k, v:= range jobList.RunningJobs {
		codes = append(codes, &models.CodesetItem{
			Id:         v.JobId,
			Name:       v.JobName,
			Status:     v.JobStatus,
			Engine:     v.JobParams.Image,
			CodePath:   v.JobParams.JobPath,
			CodeUrl:    "",
			CreateTime: time.Now().Unix() * 1000,
			Desc:       "this is description",
		})
	}

	for k, v:= range jobList.FinishedJobs {
		codes = append(codes, &models.CodesetItem{
			Id:         v.JobId,
			Name:       v.JobName,
			Status:     v.JobStatus,
			Engine:     v.JobParams.Image,
			CodePath:   v.JobParams.JobPath,
			CodeUrl:    "",
			CreateTime: time.Now().Unix() * 1000,
			Desc:       "this is description",
		})
	}

	return codes, 1, 1, nil
}

func CreateCodeset(name, description string, num int) (string, error) {
	return RandStringRunes(16), nil
}

func DeleteCodeset(id string) error {
	return nil
}

