package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"net/http"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/models"
)

var db = database.Db
var logger = loggers.Log


func GetResource() (map[string][]string, []models.DeviceItem, error) {

	fw := make(map[string][]string, 0)
	devices := make([]models.DeviceItem, 0)

	fw["tensorflow"] = make([]string, 0)
	fw["tensorflow"] = append(fw["tensorflow"],"tf_withtools:1.15")
	devices = append(devices, models.DeviceItem{
		DeviceType: "npu",
		Avail:      1,
	})

	return fw, devices, nil
}

func doRequest(url, method string, headers map[string]string, rawBody interface {}) ([]byte, error) {

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

func DoRequest(url, method string, headers map[string]string, rawBody interface {}, output interface{}) error {

	rspData, err := doRequest(url, method, headers, rawBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rspData, output)
	if err != nil {
		return err
	}

	return nil
}

