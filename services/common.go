package services

import (
	"bufio"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"go/scanner"
	"net/http"
)

var db = database.Db
var logger = loggers.Log


func GetResource() (map[string][]string, []models.DeviceItem, error) {

	fw := make(map[string]string, 0)
	devices := make([]models.DeviceItem, 0)

	fw["tensorflow"] = make([]string, 0)
		"tf_withtools:1.15"
	devices = append(devices, models.DeviceItem{
		DeviceType: "npu",
		Avail:      1,
	})

	return fw, devices, nil
}


func RequestDLTS(obj interface{}) {

	resp, err := http.Get(configs.Config.DltsUrl)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
		fmt.Println("Response status:", resp.Status)
	Print the first 5 lines of the response body.

		scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	this.Request()
	if this.err == nil && obj != nil {
		err := jsonUtil.Unmarshal(this.responseData, obj)
		if err != nil {
			this.err = fmt.Errorf("json decode error, responseText: %s, err=%v", this.ResponseText(), err)
		}
	}
	return this
}