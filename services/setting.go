package services

import (
	"github.com/apulis/AIArtsBackend/models"
)

var id = "00000000-0000-0000-0000-000000000000"

func UpsertPrivilegedSetting(setting models.PrivilegedSetting) error {
	setting.Id = id

	var settingInDb models.PrivilegedSetting
	result := db.First(&settingInDb, "id = ?", id)
	if result.RowsAffected > 0 {
		settingInDb.IsEnable = setting.IsEnable
		settingInDb.BypassCode = setting.BypassCode
		result = db.Save(&settingInDb)
		if result.Error != nil {
			return result.Error
		}
	} else {
		result = db.Create(&setting)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func GetPrivilegedSetting() (models.PrivilegedSetting, error) {
	var settingInDb models.PrivilegedSetting
	result := db.First(&settingInDb, "id = ?", id)
	if result.Error != nil {
		return settingInDb, result.Error
	} else {
		return settingInDb, nil
	}
}
