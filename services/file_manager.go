package services

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	FILETYPE_TAR_GZ = ".tar.gz"
	FILETYPE_TAR    = ".tar"
	FILETYPE_ZIP    = ".zip"
)

var FILETYPES_SUPPORTED = []string{FILETYPE_TAR_GZ, FILETYPE_TAR, FILETYPE_ZIP}

//解析zip包中的中文名称，utf8编码转为gb解决中文乱码
func transformEncode(fileName string) string {
	tempFile := bytes.NewReader([]byte(fileName))
	decoder := transform.NewReader(tempFile, simplifiedchinese.GB18030.NewDecoder())
	content, _ := ioutil.ReadAll(decoder)
	return string(content)
}

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

func CompressFile(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	dirName := filepath.Dir(path)
	fileName := fileInfo.Name()
	targetPath := dirName + "/" + fileName + ".tar.gz"

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(filepath.ToSlash(file), dirName)
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}
		return err
	})

	if err := tarWriter.Close(); err != nil {
		return "", err
	}
	if err := gzipWriter.Close(); err != nil {
		return "", err
	}

	fileToWrite, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(fileInfo.Mode()))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		return "", err
	}
	return targetPath, nil
}

func ExtractFile(fromPath, filetype, dir, isPrivate, username string) (string, error) {

	var datasetStorageDir string
	if isPrivate == "false" {
		fileConf := configs.Config.File
		datasetStorageDir = fileConf.DatasetDir + "/storage/" + dir
	} else {
		datasetStorageDir = fmt.Sprintf("/home/%s/storage/%s", username,dir)
	}
	//直接使用前端上传的path
	datasetStoragePath:=datasetStorageDir
	//datasetStoragePath := fmt.Sprintf("%s/%d", datasetStorageDir, time.Now().UnixNano())
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
	//关闭reader便于之后删除tmp文件
	defer reader.Close()

	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(toPath, transformEncode(file.Name))
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
func extractTar(fromPath, toPath string) error {
	fileReader, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileReader.Close()
	tarReader := tar.NewReader(fileReader)
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
