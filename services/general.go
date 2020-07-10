package services

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
)

var db = database.Db
var logger = loggers.Log


func GetResource() ([]models.AIFrameworkItem, []models.DeviceItem, error) {

	fw := make([]models.AIFrameworkItem, 0)
	devices := make([]models.DeviceItem, 0)

	engine := models.AIFrameworkItem{
		Name:   "tensorflow",
		Engine: "tf_withtools:1.15",
	}

	device := models.DeviceItem{
		DeviceType: "npu",
		Avail:      1,
	}

	fw = append(fw, engine)
	devices = append(devices, device)

	return fw, devices, nil
}
