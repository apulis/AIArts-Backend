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

func GetResource(userName string) (map[string][]string, []models.DeviceItem, error) {

	url := fmt.Sprintf("%s/GetVC?userName=%s&vcName=%s", configs.Config.DltsUrl, userName, models.DefaultVcName)
	vcInfo := &models.VcInfo{}

	err := DoRequest(url, "GET", nil, nil, vcInfo)
	if err != nil {
		fmt.Printf("get resource err[%+v]\n", err)
		return nil, nil, err
	}

	fw := make(map[string][]string)
	for k, v := range configs.Config.Image {

		fw[k] = make([]string, 0)
		for _, item := range v {
			fw[k] = append(fw[k], item)
		}
	}

	devices := make([]models.DeviceItem, 0)
	for k, v := range vcInfo.DeviceAvail {
		devices = append(devices, models.DeviceItem{
			DeviceType: k,
			Avail:      v,
		})
	}

	return fw, devices, nil
}
