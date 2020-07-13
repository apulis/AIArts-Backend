package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/models"
)

var db = database.Db
var logger = loggers.Log


func GetResource(userName string) (map[string][]string, []models.DeviceItem, error) {

	url := fmt.Sprintf("http://atlas02.sigsus.cn/apis/GetVC?userName=%s&vcName=%s", userName, models.DefaultVcName)

	vcInfo := &models.VcInfo{}
	err := DoRequest(url, "GET", nil, nil, vcInfo)
	if err != nil {
		fmt.Printf("get resource err[%+v]\n", err)
		return nil, nil, err
	}

	fw := make(map[string][]string)
	fw["tensorflow"], fw["mindspore"] = make([]string, 0), make([]string, 0)

	// todo: must read from config
	fw["tensorflow"] = append(fw["tensorflow"],"apulistech/ubuntu_tf_test:1.22")
	fw["tensorflow"] = append(fw["tensorflow"],"ubuntu:18.04")
	fw["mindspore"] = append(fw["mindspore"],"apulistech/mindspore:0.3.0-withtools")

	devices := make([]models.DeviceItem, 0)
	for k, v := range vcInfo.DeviceAvail {
		if v > 0 {
			devices = append(devices, models.DeviceItem{
				DeviceType: k,
				Avail:      v,
			})
		}
	}

	return fw, devices, nil
}


