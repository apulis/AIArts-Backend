package services

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/apulis/AIArtsBackend/models"
)

func GetUpgradeLog() (string, string, error) {
	var status string
	var Log string
	var err error
	progress := models.Upgrade_Progress
	switch progress {
	case -1:
		status = "not ready"
	case 100:
		status = "success"
		models.Upgrade_Progress = -1
	default:
		status = "upgrading"
		Log, err = acquireLog()
		if err != nil {
			return "error", Log, err
		}
		fmt.Println(string(Log))
	}
	return status, Log, nil
}

func acquireLog() (string, error) {
	if !isFileExists(models.UPGRADE_FILE_PATH + "/" + "/upgrade.log") {
		return "prepare environment", nil
	}
	// cmd := exec.Command("/bin/bash", "-c", "tail -n 2500 "+models.UPGRADE_FILE_PATH+"/upgrade.log")
	cmd := exec.Command("/bin/bash", "-c", "wc -l "+models.UPGRADE_FILE_PATH+"/upgrade.log"+" | awk '{print $1}' | tr -d '\\n' ")
	lineCountOutput, err := cmd.Output()
	if err != nil {
		err = errors.New("count log file lines fail")
		fmt.Println("Execute Command failed:" + err.Error())
		return "", err
	}
	fmt.Println("line count :" + string(lineCountOutput) + ";")
	lineCount, err := strconv.Atoi(string(lineCountOutput))
	if err != nil {
		fmt.Println("converte fail")
		fmt.Println("Execute Command failed:" + err.Error())
		return "", err
	}

	if lineCount > models.Log_Line_Point {
		fmt.Println("latest line: " + strconv.Itoa(lineCount) + "; old line: " + strconv.Itoa(models.Log_Line_Point))
		cmd = exec.Command("/bin/bash", "-c", "sed -n '"+strconv.Itoa(models.Log_Line_Point)+","+strconv.Itoa(lineCount)+"p' "+models.UPGRADE_FILE_PATH+"/upgrade.log")
		models.Log_Line_Point = lineCount
		log, err := cmd.Output()
		if err != nil {
			err = errors.New("get log file fail")
			fmt.Println("Execute Command failed:" + err.Error())
			return "", err
		}
		return string(log), nil
	}

	return "", nil
}

func GetUpgradeProgress() (string, int) {
	var status string
	var progress int
	progress = models.Upgrade_Progress
	switch progress {
	case -1:
		status = "not ready"
		progress = 0
	case 100:
		status = "success"
		models.Upgrade_Progress = -1
	default:
		status = "upgrading"
	}
	return status, progress
}

func UpgradePlatformByLocal() error {
	if models.Upgrade_Progress != -1 {
		fmt.Println("upgrading, please wait until upgrade finish")
		return errors.New("upgrading, please wait until upgrade finish")
	}
	go UpgradePlatformdLocally()
	return nil
}

func UpgradePlatformdLocally() error {
	models.Upgrade_Progress = 0
	// upgradeFiles, err := ioutil.ReadDir(models.UPGRADE_FILE_PATH)
	// if err != nil {
	// 	return err
	// }
	// for _, file := range upgradeFiles {
	// 	fmt.Println(file.Name())
	// }
	upgradeConfig, err := models.GetUpgradeConfig()
	if err != nil {
		return err
	}
	upgradeScript := upgradeConfig.UpgradeScript
	cmd := exec.Command("/bin/bash", "-c", models.UPGRADE_FILE_PATH+"/"+upgradeScript+" > "+models.UPGRADE_FILE_PATH+"/"+"/upgrade.log")
	err = cmd.Run()

	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
		return err
	}
	fmt.Println(upgradeConfig.Version)
	newVersion := upgradeConfig.Version
	description := upgradeConfig.Description
	versionInfoSet := models.VersionInfoSet{
		Description: description,
		Version:     newVersion,
	}
	cmd = exec.Command("/bin/bash", "-c", "mkdir -p /var/log")
	err = cmd.Run()
	if err != nil {
		err = errors.New("mkdir fail")
		fmt.Println("Execute Command failed:" + err.Error())
		return err
	}
	cmd = exec.Command("/bin/bash", "-c", "mv "+models.UPGRADE_FILE_PATH+"/"+"/upgrade.log"+" /var/log/upgrade.log")
	err = cmd.Run()
	if err != nil {
		err = errors.New("move log fail")
		fmt.Println("Execute Command failed:" + err.Error())
		return err
	}
	models.Upgrade_Progress = 100
	return models.UploadVersionInfoSet(versionInfoSet)
}

func GetLocalUpgradeEnv() (bool, bool, error) {
	upgradeFilePath := models.UPGRADE_FILE_PATH
	fmt.Print(upgradeFilePath)
	canUpgrade := true
	if !isFileExists(upgradeFilePath) {
		canUpgrade = false
	}
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
