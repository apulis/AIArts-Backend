package services

import (
	"github.com/apulis/AIArtsBackend/models"
)

func checkUpgradeEnvironment() {

}

func UploadUpgradeInfo() error {
	versionInfoSet := models.VersionInfoSet{
		Description: "this is a wonderful version",
		Version:     "v0.0.2",
	}
	return models.UploadVersionInfoSet(versionInfoSet)
}

func GetCurrentVersion() (models.VersionInfoSet, error) {

	return models.GetCurrentVersion()

}

func GetVersionLogs() ([]models.VersionInfoSet, error) {
	return models.GetVersionLogs()
}
