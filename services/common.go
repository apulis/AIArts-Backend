package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/levigross/grequests"
	"strconv"
	"strings"
)

var db = database.Db
var logger = loggers.Log

func GetResource(userName, vcName string) (*models.VcInfo, error) {

	url := fmt.Sprintf("%s/GetVC?userName=%s&vcName=%s", configs.Config.DltsUrl, userName, vcName)
	vcInfo := &models.VcInfo{}

	err := DoRequest(url, "GET", nil, nil, vcInfo)
	if err != nil {
		fmt.Printf("get resource err[%+v]\n", err)
		return nil, err
	}

	url = fmt.Sprintf("%s/GetAllDevice?userName=%s", configs.Config.DltsUrl, userName)
	devices := make(map[string]models.DeviceItem2)

	err = DoRequest(url, "GET", nil, nil, &devices)
	if err != nil {
		fmt.Printf("get all devices err[%+v]\n", err)
		return nil, err
	}

	for _, v := range vcInfo.Nodes {
		if val, ok := devices[v.GPUType]; ok {
			v.DeviceStr = val.DeviceStr
		}
	}

	return vcInfo, nil
}

func GetJobSummary(userName, jobType, vcName string) (map[string]int, error) {

	url := fmt.Sprintf("%s/GetJobSummary?userName=%s&jobType=%s&vcName=%s", configs.Config.DltsUrl, userName, jobType, vcName)

	summary := make(map[string]int)
	convert_summary := make(map[string]int)

	err := DoRequest(url, "GET", nil, nil, &summary)
	if err != nil {
		fmt.Printf("get job summary err[%+v]\n", err)
		return nil, err
	}

	for k, v := range summary {
		convert_summary[strings.ToLower(k)] = v
	}

	return convert_summary, nil
}

func GetResources(userName string) (interface{}, error) {
	BackendUrl = configs.Config.Infer.BackendUrl
	resp, err := grequests.Get(BackendUrl+"/apis/GetAllDevice?userName="+userName, nil)
	if resp.StatusCode != 200 {
		logger.Error("response code is ", resp.StatusCode, resp.String())
		return "", errors.New("response code: " + (strconv.Itoa(resp.StatusCode)) + ",detail: " + resp.String())
	}
	var resources interface{}
	json.Unmarshal(resp.Bytes(), &resources)
	return resources, err
}

