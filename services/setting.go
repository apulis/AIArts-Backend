package services

import (
	"errors"

	"github.com/apulis/AIArtsBackend/models"
	"github.com/jinzhu/gorm"
)

var id = "00000000-0000-0000-0000-000000000000"

func UpsertPrivilegedSetting(setting models.PrivilegedSetting) error {
	setting.Id = id

	var settingInDb models.PrivilegedSetting
	result := db.First(&settingInDb, "id = ?", id)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

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
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return settingInDb, result.Error
	} else {
		return settingInDb, nil
	}
}
