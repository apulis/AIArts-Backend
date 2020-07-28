package models

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type VersionInfoSet struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Description string `gorm:"type:text" json:"description"`
	Version     string `gorm:"not null" json:"version"`
	Creator     string `json:"creator"`
}

type UpgradeYaml struct {
	Version       string `yaml:"Version"`
	UpgradeScript string `yaml:"UpgradeScript"`
	Description   string `yaml:"Description"`
	Creator       string `yaml:"Creator"`
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
	yamlFile, err := ioutil.ReadFile(UPGRADE_FILE_PATH + "/" + UPGRADE_CONFIG_FILE)
	if err != nil {
		fmt.Errorf("Fatal error read upgrade config file: %s \n", err)
		return config, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Errorf("Fatal error unmarsh upgrade config file: %s \n", err)
		return config, err
	}
	return config, nil
}
