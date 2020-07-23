package services

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apulis/AIArtsBackend/models"
)

func UploadUpgradeInfo() error {
	versionInfoSet := models.VersionInfoSet{
		Description: "this is a wonderful version",
		Version:     "v0.0.3",
	}
	upgradeFiles, err := ioutil.ReadDir(models.UPGRADE_FILE_PATH)
	if err != nil {
		return err
	}
	for _, file := range upgradeFiles {
		fmt.Println(file.Name())
	}
	return models.UploadVersionInfoSet(versionInfoSet)
}

func GetLocalUpgradeEnv() (bool, bool, error) {
	upgradeFilePath := models.UPGRADE_FILE_PATH
	fmt.Print(upgradeFilePath)
	canUpgrade := true
	if !isFileExists(upgradeFilePath) {
		canUpgrade = false
	}
	packageVersion, err := ioutil.ReadFile(models.UPGRADE_FILE_PATH + "/version")
	if err != nil {
		return false, false, err
	}
	fmt.Printf("local package version :%s", packageVersion)

	return canUpgrade, false, nil
}

func GetCurrentVersion() (models.VersionInfoSet, error) {

	return models.GetCurrentVersion()

}

func GetVersionLogs() ([]models.VersionInfoSet, error) {
	// var versionLogs []models.VersionInfoSet
	// versionLogs, err := models.GetVersionLogs()
	// if err != nil {
	// 	return nil, err
	// }
	// var logs []string = make([]string, 0)
	// for _, versionInfo := range versionLogs {
	// 	logs = append(logs, versionInfo.Version + "update in" + )
	// }
	// return logs, nil
	return models.GetVersionLogs()
}

func isFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
