package services

import (
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/models"
)

var db = database.Db
var logger = loggers.Log


func GetResource() (map[string][]string, []models.DeviceItem, error) {

	fw := make(map[string][]string, 0)
	devices := make([]models.DeviceItem, 0)

	fw["tensorflow"] = make([]string, 0)
	fw["tensorflow"] = append(fw["tensorflow"],"tf_withtools:1.15")
	devices = append(devices, models.DeviceItem{
		DeviceType: "npu",
		Avail:      1,
	})

	return fw, devices, nil
}


