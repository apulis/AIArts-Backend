package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apulis/AIArtsBackend/configs"
)

const (
	FILETYPE_TAR_GZ = ".tar.gz"
	FILETYPE_TAR    = ".tar"
	FILETYPE_ZIP    = ".zip"
)

var FILETYPES_SUPPORTED = []string{FILETYPE_TAR_GZ, FILETYPE_TAR, FILETYPE_ZIP}

func CheckFileName(filename string) (string, error) {
	for _, filetype := range FILETYPES_SUPPORTED {
		if strings.HasSuffix(filename, filetype) {
			return filetype, nil
		}
	}

	logger.Info("File type not supported: ", filename)
	return "", errors.New("File type not supported")
}

func GetDatasetTempPath(filetype string) (string, error) {
	fileConf := configs.Config.File
	datasetTempDir := fileConf.DatasetDir + "/tmp"
	_, err := os.Stat(datasetTempDir)
	if err != nil {
		err = os.MkdirAll(datasetTempDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	datasetTempPath := fmt.Sprintf("%s/%d%s", datasetTempDir, time.Now().UnixNano(), filetype)
	return datasetTempPath, nil

}

func ExtractFile(fromPath, filetype string) (string, error) {
	fileConf := configs.Config.File
	datasetStorageDir := fileConf.DatasetDir + "/storage"
	datasetStoragePath := fmt.Sprintf("%s/%d", datasetStorageDir, time.Now().UnixNano())
	_, err := os.Stat(datasetStoragePath)
	if err != nil {
		err = os.MkdirAll(datasetStoragePath, os.ModeDir|os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	logger.Info("Extracting file: ", fromPath, " to ", datasetStoragePath)
	switch filetype {
	case FILETYPE_ZIP:
		err = extractZip(fromPath, datasetStoragePath)
	case FILETYPE_TAR_GZ:
		err = extractTarGz(fromPath, datasetStoragePath)
	case FILETYPE_TAR:
		err = extractTar(fromPath, datasetStoragePath)
	default:
		err = errors.New("Unknown file type")
	}

	if err != nil {
		logger.Info("Extracting '", fromPath, "' failed")
		return "", err
	}
	return datasetStoragePath, nil
}

func extractZip(fromPath, toPath string) error {
	reader, err := zip.OpenReader(fromPath)
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		fmt.Println(file)
		path := filepath.Join(toPath, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func extractTarGz(fromPath, toPath string) error {
	return nil
}

func extractTar(fromPath, toPath string) error {
	return nil
}
