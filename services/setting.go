package services

import (
	"github.com/apulis/AIArtsBackend/database"
)

var db = database.Db
var id = "00000000-0000-0000-0000-000000000000"

func UpsertPrivilegedSetting(settings models.PrivilegedSetting) error {
	settings["id"] = id

	var settingsInDb = models.PrivilegedSetting
	result := db.First(&settingsInDb, id)
	if result.RowsAffected > 0 {
		settingsInDb.IsEnable = settings.IsEnable
		settingsInDb.BypassCode = settings.BypassCode
		result = db.Save(&settingsInDb)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result = db.Create(&settings)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func GetPrivilegedSetting() (models.PrivilegedSetting, error) {
	var settingsInDb = models.PrivilegedSetting
	result := db.First(&settingsInDb, id)
	if result.Error != nil {
		return nil, result.Error
	}
	else {
		return settingsInDb, nil
	}
}