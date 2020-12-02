package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)



func doRequest(url, method string, headers map[string]string, rawBody interface{}) ([]byte, error) {

	var body io.Reader = nil
	if rawBody != nil {
		switch t := rawBody.(type) {
		case string:
			body = strings.NewReader(t)

		case []byte:
			body = bytes.NewReader(t)

		default:
			data, err := json.Marshal(rawBody)
			if err != nil {
				err = fmt.Errorf("body 序列化失败: %v", err)
				return nil, err
			}

			body = bytes.NewReader(data)
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if len(headers) != 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// read response
	responseData := make([]byte, 0)
	responseData, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 0 {

		Status := resp.Status
		StatusCode := resp.StatusCode

		if StatusCode < 200 || StatusCode >= 400 {
			err = errors.New(Status)
		}
	}

	return responseData, err
}

func DoRequest(url, method string, headers map[string]string, rawBody interface{}, output interface{}) error {

	rspData, err := doRequest(url, method, headers, rawBody)
	if err != nil {
		return err
	}
	if len(rspData) > 0 {
		err = json.Unmarshal(rspData, output)
		if err != nil {
			return err
		}
	}
	logger.Info(url)
	logger.Info(output)
	return nil
}

func DoGetRequest(url string, headers map[string]string, rawBody interface{}) (err error, rawData string) {
	rspData, err := doRequest(url, "GET", headers, rawBody)
	if err != nil {
		return err, ""
	}

	return nil, string(rspData)
}


// 如果配置了私有仓库，则添加私有仓库前缀
func ConvertImage(image string) string {
	imageName := strings.TrimSpace(image)
	if len(configs.Config.PrivateRegistry) > 0 {
		// 不带私有仓库前缀
		if !strings.HasPrefix(imageName, configs.Config.PrivateRegistry) {
			if strings.HasSuffix(configs.Config.PrivateRegistry, "/") {
				imageName = configs.Config.PrivateRegistry + imageName
			} else {
				imageName = configs.Config.PrivateRegistry + "/" + imageName
			}
		}
	}
	return imageName
}

// 如果配置了私有仓库，则删除掉
func UnConvertImage(image string) string {
	imageName := strings.TrimSpace(image)
	if len(configs.Config.PrivateRegistry) > 0 {
		// 如果带私有仓库前缀
		if strings.HasPrefix(imageName, configs.Config.PrivateRegistry) {
			imageName = strings.ReplaceAll(imageName, configs.Config.PrivateRegistry, "")
		}
	}
	return imageName
}

//获取启动文件的类型sh或者python
func CheckStartFileType(filename string) (string, error) {
	for _, filetype := range STARTFILETYPES_SUPPORTED {
		if strings.HasSuffix(filename, filetype) {
			return filetype, nil
		}
	}
	logger.Info("StartFile type not supported: ", filename)
	return "", errors.New("StartFile type not supported")
}
