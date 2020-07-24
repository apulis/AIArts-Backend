package services

import (
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
)

func LsEdgeInferences(pageNum, pageSize int, user string) error {
	fmt.Println(pageNum, pageSize)
	test := make(map[string]interface{})
	url := fmt.Sprintf("%s/ListModelConversionJob?vcName=%s&jobOwner=%s", configs.Config.DltsUrl, models.DefaultVcName, user)

	err := DoRequest(url, "GET", nil, nil, &test)
	return err
}

func GetFDInfo() (models.FDInfo, error) {
	var fd models.FDInfo
	url := fmt.Sprintf("%s/GetFDInfo", configs.Config.DltsUrl)

	err := DoRequest(url, "GET", nil, nil, &fd)
	return fd, err
}

func SetFDInfo(username, password, reqUrl string) (bool, error) {
	url := fmt.Sprintf("%s/SetFDInfo", configs.Config.DltsUrl)
	params := make(map[string]interface{})
	params["username"] = username
	params["password"] = password
	params["url"] = reqUrl

	res := false
	err := DoRequest(url, "POST", nil, params, &res)

	return res, err
}

func GetConversionTypes() (models.ConvertionTypes, error) {
	var convTypes models.ConvertionTypes
	url := fmt.Sprintf("%s/GetModelConversionTypes", configs.Config.DltsUrl)

	err := DoRequest(url, "GET", nil, nil, &convTypes)
	return convTypes, err
}
