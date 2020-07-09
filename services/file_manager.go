package services

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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

func CheckPathExists(path string) error {
	_, err := os.Stat(path)
	return err
}

func CheckDatasetPathValid(path string) error {
	datasetPathPrefix := configs.Config.File.DatasetDir + "/storage"
	if strings.HasPrefix(path, datasetPathPrefix) {
		return nil
	} else {
		return errors.New(fmt.Sprint("Dataset path should be in: ", datasetPathPrefix))
	}
}

func CheckFileName(filename string) (string, error) {
	for _, filetype := range FILETYPES_SUPPORTED {
		if strings.HasSuffix(filename, filetype) {
			return filetype, nil
		}
	}

	logger.Info("File type not supported: ", filename)
	return "", errors.New("File type not supported")
}

func CheckFileOversize(size int64) bool {
	fileConf := configs.Config.File
	sizeLimit := fileConf.SizeLimit
	logger.Info("Upload file size: ", size, ". Config size limit: ", sizeLimit)
	if int(size) < sizeLimit {
		return false
	}
	return true
}

func GetDirSize(path string) (int, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return int(size), err
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
		err = extractTarGz(fromPath, datasetStoragePath)
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
	fileReader, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	gzipReader, err := gzip.NewReader(fileReader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		head, err := tarReader.Next()
		if head == nil || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		path := filepath.Join(toPath, head.Name)
		fileInfo := head.FileInfo()
		switch head.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(path, fileInfo.Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
			if err != nil {
				return err
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, tarReader); err != nil {
				return err
			}
		default:
			logger.Info("Extracting unknown type ", string(head.Typeflag))
		}
	}

	return nil
}
