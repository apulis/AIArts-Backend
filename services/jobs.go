package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"net/url"
)

func GetJobsCount(req models.GetAllJobsReq) (int, error) {
	url := fmt.Sprintf(`%s/GetJobCount?vcName=%s&jobType=%s&jobStatus=%s&searchWord=%s`,
		configs.Config.DltsUrl, req.VCName,
		req.JobType, req.JobStatus, url.QueryEscape(req.SearchWord))

	ret := struct {
		Count int `json:"count"`
	}{}
	err := DoRequest(url, "GET", nil, nil, &ret)
	if err != nil {
		fmt.Printf("get job count err[%+v]", err)
		return 0, err
	}

	return ret.Count, nil
}

func GetAllJobs(req models.GetAllJobsReq) ([]models.Job, error) {
	url := fmt.Sprintf(`%s/ListAllJobs?vcName=%s&jobType=%s&pageNum=%d&pageSize=%d&jobStatus=%s&searchWord=%s&orderBy=%s&order=%s`,
		configs.Config.DltsUrl, req.VCName,
		req.JobType, req.PageNum, req.PageSize, req.JobStatus, url.QueryEscape(req.SearchWord),
		req.OrderBy, req.Order)

	ret := struct {
		Jobs []models.Job `json:"jobs"`
	}{}
	err := DoRequest(url, "GET", nil, nil, &ret)

	if err != nil {
		fmt.Printf("get all code err[%+v]", err)
		return nil, err
	}
	return ret.Jobs, nil
}