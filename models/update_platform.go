package models

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type VersionInfoSet struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Description string `gorm:"type:text" json:"description"`
	Version     string `gorm:"not null" json:"version"`
}

type UpgradeYaml struct {
	Version       string
	UpgradeScript string
	Description   string
}

func UploadVersionInfoSet(versionInfo VersionInfoSet) error {
	return db.Create(&versionInfo).Error
}

func GetCurrentVersion() (VersionInfoSet, error) {
	var currentVersion VersionInfoSet
	res := db.Order("created_at desc").First(&currentVersion)
	if res.Error != nil {
		return currentVersion, res.Error
	}
	return currentVersion, nil
}

func GetVersionLogs() ([]VersionInfoSet, error) {
	var versionLog []VersionInfoSet
	res := db.Order("created_at desc").Find(&versionLog)
	if res.Error != nil {
		return versionLog, res.Error
	}
	return versionLog, nil
}

func GetUpgradeConfig() (UpgradeYaml, error) {
	var config UpgradeYaml
	viper.SetConfigName(strings.Replace(UPGRADE_CONFIG_FILE, ".yaml", "", -1))
	viper.AddConfigPath(UPGRADE_FILE_PATH)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error read upgrade config file: %s \n", err))
		return config, err
	}

	viper.Unmarshal(&config)
	return config, err
}
