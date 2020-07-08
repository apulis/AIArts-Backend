package services

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
)

var db = database.Db
var logger = loggers.Log


func GetResource() ([]models.AIFrameworkItem, []models.DeviceItem, error) {

	return nil, nil, nil
}
