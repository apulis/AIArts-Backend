package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func GetAllTemplate(userName string, page, size int, jobStatus, searchWord string) ([]*models.Template, int, int, error) {

	url := ""

	jobList := &models.JobList{}
	err := DoRequest(url, "GET", nil, nil, jobList)

	if err != nil {
		fmt.Printf("get all Template err[%+v]", err)
		return nil, 0, 0, err
	}

	Templates := make([]*models.Template, 0)
	return Templates, len(Templates), 1, nil
}

func CreateTemplate(userName string, Template models.Template) (string, error) {

	url := fmt.Sprintf("%s/PostJob", configs.Config.DltsUrl)
	return url, nil
}

func DeleteTemplate(userName, id string) error {

	return nil
}

func GetTemplate(userName, id string) (*models.Template, error) {

	template := &models.Template{}
	return template, nil
}
