package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
)

var db = database.Db
var logger = loggers.Log
var DefaultVcName = "atlas"

type UnixTime struct {
	time.Time
}

func init() {
	createTableIfNotExists(Dataset{})
	createTableIfNotExists(Modelset{})
}

func createTableIfNotExists(modelType interface{}) {
	val := reflect.Indirect(reflect.ValueOf(modelType))
	modelName := val.Type().Name()

	hasTable := db.HasTable(modelType)
	if !hasTable {
		logger.Info(fmt.Sprintf("Table of %s not exists, create it.", modelName))
		db.CreateTable(modelType)
	} else {
		logger.Info(fmt.Sprintf("Table of %s already exists.", modelName))
	}
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	microSec := t.Unix() * 1000
	return []byte(strconv.FormatInt(microSec, 10)), nil
}

func (t UnixTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *UnixTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = UnixTime{Time: value}
		return nil
	}
	return fmt.Errorf("cannot convert %v to timestamp", v)
}

type JobParams struct {
	Cmd                   string `json:"cmd"`
	ContainerUserId       int 	`json:"containerUserId"`
	DataPath              string `json:"dataPath"`
	Enabledatapath        bool `json:"enabledatapath"`
	Enablejobpath         bool `json:"enablejobpath"`
	Enableworkpath        bool `json:"enableworkpath"`
	Env                   []string `json:"env"`
	FamilyToken           string `json:"familyToken"`
	GpuType               string `json:"gpuType"`
	HostNetwork           bool `json:"hostNetwork"`
	Image                 string `json:"image"`
	InteractivePorts      bool `json:"interactivePorts"`
	IsParent              int `json:"isParent"`
	IsPrivileged          bool `json:"isPrivileged"`
	JobId                 string `json:"jobId"`
	JobName               string `json:"jobName"`
	JobPath               string `json:"jobPath"`
	JobType               string `json:"jobType"`
	Jobtrainingtype       string `json:"jobtrainingtype"`
	Numps                 int `json:"numps"`
	Numpsworker			  int `json:"numpsworker"`
	PreemptionAllowed     bool `json:"preemptionAllowed"`
	Resourcegpu           int `json:"resourcegpu"`
	Team                  string `json:"team"`
	UserId                int `json:"userId"`
	UserName              string `json:"userName"`
	VcName                string `json:"vcName"`
	WorkPath              string `json:"workPath"`
}

type Job struct {
	JobId 		string `json:"jobId"`
	JobName		string `json:"jobName"`
	JobParams   JobParams `json:"jobParams"`
	JobStatus   string `json:"jobStatus"`
	JobTime		string `json:"jobTime"`
	JobType     string `json:"jobType"`
	Priority    int `json:"priority"`
	UserName	string `json:"userName"`
	VcName		string `json:"vcName"`
}

type JobMeta struct {
	FinishedJobs int `json:"finishedJobs"`
	QueuedJobs   int `json:"queuedJobs"`
	RunningJobs  int `json:"runningJobs"`
	VisualizationJobs int `json:"visualizationJobs"`
}

type JobList struct {
	FinishedJobs []*Job `json:"finishedJobs"`
	Meta JobMeta `json:"meta"`
	QueuedJobs []*Job `json:"queuedJobs"`
	RunningJobs []*Job `json:"runningJobs"`
}


type VcInfo struct {
	DeviceAvail 	map[string]int `json:"gpu_avaliable"`
}

type JobId struct {
	Id 		string `json:"jobId"`
}

type UriJobId struct {
	Id string `uri:"id" binding:"required"`
}