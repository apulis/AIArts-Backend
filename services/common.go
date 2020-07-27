package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/models"
)

var db = database.Db
var logger = loggers.Log

func GetResource(userName string) (*models.VcInfo, error) {

	url := fmt.Sprintf("%s/GetVC?userName=%s&vcName=%s", configs.Config.DltsUrl, userName, models.DefaultVcName)
	vcInfo := &models.VcInfo{}

	err := DoRequest(url, "GET", nil, nil, vcInfo)
	if err != nil {
		fmt.Printf("get resource err[%+v]\n", err)
		return nil, err
	}

	return vcInfo, nil
}

func GetJobSummary(userName, jobType string) (map[string]int, error) {

	url := fmt.Sprintf("%s/GetJobSummary?userName=%s&jobType=%s", configs.Config.DltsUrl, userName, jobType)
	summary := make(map[string]int)

	err := DoRequest(url, "GET", nil, nil, &summary)
	if err != nil {
		fmt.Printf("get job summary err[%+v]\n", err)
		return nil, err
	}

	return summary, nil
}
