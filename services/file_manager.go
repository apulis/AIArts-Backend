package services

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	FILETYPE_TAR_GZ = ".tar.gz"
	FILETYPE_TAR    = ".tar"
	FILETYPE_ZIP    = ".zip"
	FILETYPE_SHELL  = ".sh"
	FILETYPE_PYTHON = ".py"
	FILETYPE_JSON   = ".json"
)

var FILETYPES_SUPPORTED = []string{FILETYPE_TAR_GZ, FILETYPE_TAR, FILETYPE_ZIP}
var STARTFILETYPES_SUPPORTED = []string{FILETYPE_SHELL, FILETYPE_PYTHON}

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

func GetDirSize(path string) (int64, error) {
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
	return size, err
}

func GetDatasetTempPath(dir string) (string, error) {

	fileConf := configs.Config.File
	datasetTempDir := fileConf.DatasetDir + "/tmp"

	_, err := os.Stat(datasetTempDir)
	if err != nil {
		err = os.MkdirAll(datasetTempDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return "",  err
		}
	}

	// 解压目录
	datasetTempPath := fmt.Sprintf("%s/%s", datasetTempDir, dir)
	return datasetTempPath, nil
}

func GetModelTempPath(filetype string) (string, error) {
	fileConf := configs.Config.File
	modelTempDir := fileConf.ModelDir + "/tmp"
	_, err := os.Stat(modelTempDir)
	if err != nil {
		err = os.MkdirAll(modelTempDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	modelTempPath := fmt.Sprintf("%s/%d%s", modelTempDir, time.Now().UnixNano(), filetype)
	return modelTempPath, nil
}

func CompressFile(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	dirName := filepath.Dir(path)
	fileName := fileInfo.Name()
	tmpDir := dirName + "/../tmp/"
	_, err = os.Stat(tmpDir)
	if err != nil {
		err = os.MkdirAll(tmpDir, os.ModeDir|os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	targetPath := tmpDir + fileName + strconv.FormatInt(time.Now().Unix(), 10) + ".tar.gz"

	var cmd *exec.Cmd
	cmd = exec.Command("tar", "-zcf", targetPath, path)
	if _, err := cmd.Output(); err != nil {
		logger.Error(fmt.Sprintf("run cmd %s error: %s", cmd.String(), err.Error()))
		return "", err
	}

	return targetPath, nil
}

func GenerateDatasetStoragePath(username , folderName string, isPrivate string) string {

	var datasetStoragePath string
	fileConf := configs.Config.File

	//直接使用前端上传的path
	if isPrivate == "false" {
		datasetStoragePath = fileConf.DatasetDir + "/storage/" + folderName
	} else {
		datasetStoragePath = fmt.Sprintf("/home/%s/storage/%s", username, folderName)
	}

	return datasetStoragePath
}

func GenerateModelStoragePath(dir, username string) string {
	var datasetStoragePath string
	datasetStoragePath = fmt.Sprintf("/home/%s/storage/%s", username, dir)
	//直接使用前端上传的path
	//debug
	if gin.Mode() == "debug" {
		if username == "kaiyuan.xu" {
			datasetStoragePath = fmt.Sprintf("D:/work/tmp/%s/storage/%s", username, dir)
		}
	}
	return datasetStoragePath
}

func ExtractFile(fromPath, filetype, extractPath string) (string, error) {

	_, err := os.Stat(extractPath)
	if err != nil && os.IsNotExist(err) {

		mask := syscall.Umask(0)  // 改为 0000 八进制
		defer syscall.Umask(mask) // 改为原来的 umask

		err = os.MkdirAll(extractPath, os.ModeDir|os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	logger.Info("Extracting file: ", fromPath, " to ", extractPath)
	switch filetype {
	case FILETYPE_ZIP:
		err = extractZip(fromPath, extractPath)
	case FILETYPE_TAR_GZ:
		err = extractTarGz(fromPath, extractPath)
	case FILETYPE_TAR:
		err = extractTar(fromPath, extractPath)
	default:
		err = errors.New("Unknown file type")
	}

	if err != nil {
		logger.Info("Extracting '", fromPath, "' failed")
		return "", err
	}

	return extractPath, nil
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
		//如果直接递归到底层是文件 比如 /data/pic/train/1.png 那么先要创建pic文件夹,linux与windows的zip压缩包文件夹头结构不一样
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			fileReader.Close()
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {

			fileReader.Close()
			targetFile.Close()

			if os.IsNotExist(err) {
				logger.Error("Ignored not existed file: ", path)
				return nil
			} else {
				return err
			}
		} else if _, err := io.Copy(targetFile, fileReader); err != nil {

			fileReader.Close()
			targetFile.Close()
			return err
		}

		fileReader.Close()
		targetFile.Close()
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
		path := filepath.Join(toPath, transformEncode(head.Name))
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

			if _, err := io.Copy(targetFile, tarReader); err != nil {
				targetFile.Close()
				return err
			}

			targetFile.Close()
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

		path := filepath.Join(toPath, transformEncode(head.Name))
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

			if _, err := io.Copy(targetFile, tarReader); err != nil {
				targetFile.Close()
				return err
			}
			targetFile.Close()

		default:
			logger.Info("Extracting unknown type ", string(head.Typeflag))
		}
	}

	return nil
}