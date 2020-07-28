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
var DefaultVcName = "platform"

type UnixTime struct {
	time.Time
}

func init() {
	createTableIfNotExists(Dataset{})
	createTableIfNotExists(Modelset{})
	createTableIfNotExists(VersionInfoSet{})
	createTableIfNotExists(Templates{})
}

//驼峰转下划线形式
func CamelToCase(name string) string {
	var upperStr string
	for _, v := range name {
		if v >= 65 && v <= 90 {
			upperStr += "_" + string(v+32)
		} else {
			upperStr += string(v)
		}

	}
	return upperStr
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

// 以下结构体用于api/common
type DeviceItem struct {
	DeviceType string `json:"deviceType"`
	Avail      int    `json:"avail"`
}

type NodeInfo struct {
	TotalNodes        int            `json:"totalNodes"`
	CountByDeviceType map[string]int `json:"countByDeviceType"`
}

// 模板类别
const (
	TemplatePublic     int = 1
	TemplatePrivate    int = 2
	TemplateUserPublic int = 3 // 读取用户列表兼公共列表
)

// 以下结构体用于和DLTS平台交互
const (
	JobTypeTraining     string = "training"     // 老DLTS默认采用的jobType
	JobTypeArtsTraining string = "artsTraining" // 供电局项目：模型训练
	JobTypeCodeEnv      string = "codeEnv"      // 供电局项目：代码环境

	JobStatusAll string = "all"
)

const (
	TrainingTypeDist    string = "PSDistJob"
	TrainingTypeRegular string = "RegularJob"
)

type JobParams struct {
	Cmd             string `json:"cmd"`
	ContainerUserId int    `json:"containerUserId"`

	Enabledatapath    bool     `json:"enabledatapath"`
	Enablejobpath     bool     `json:"enablejobpath"`
	Enableworkpath    bool     `json:"enableworkpath"`
	Env               []string `json:"env"`
	FamilyToken       string   `json:"familyToken"`
	GpuType           string   `json:"gpuType"`
	HostNetwork       bool     `json:"hostNetwork"`
	Image             string   `json:"image"`
	InteractivePorts  bool     `json:"interactivePorts"`
	IsParent          int      `json:"isParent"`
	IsPrivileged      bool     `json:"isPrivileged"`
	JobId             string   `json:"jobId"`
	JobName           string   `json:"jobName"`
	JobPath           string   `json:"jobPath"`
	JobType           string   `json:"jobType"`
	Jobtrainingtype   string   `json:"jobtrainingtype"`
	Numps             int      `json:"numps"`
	Numpsworker       int      `json:"numpsworker"`
	PreemptionAllowed bool     `json:"preemptionAllowed"`
	Resourcegpu       int      `json:"resourcegpu"`
	Team              string   `json:"team"`
	UserId            int      `json:"userId"`
	UserName          string   `json:"userName"`
	VcName            string   `json:"vcName"`
	WorkPath          string   `json:"workPath"`

	CodePath    string `json:"codePath"`
	StartupFile string `json:"startupFile"`
	OutputPath  string `json:"outputPath"`
	DatasetPath string `json:"datasetPath"`
	Desc        string `json:"desc"`
}

type Job struct {
	JobId     string    `json:"jobId"`
	JobName   string    `json:"jobName"`
	JobParams JobParams `json:"jobParams"`
	JobStatus string    `json:"jobStatus"`
	JobTime   string    `json:"jobTime"`
	JobType   string    `json:"jobType"`
	Priority  int       `json:"priority"`
	UserName  string    `json:"userName"`
	VcName    string    `json:"vcName"`
}

type JobMeta struct {
	FinishedJobs      int `json:"finishedJobs"`
	QueuedJobs        int `json:"queuedJobs"`
	RunningJobs       int `json:"runningJobs"`
	VisualizationJobs int `json:"visualizationJobs"`
	TotalJobs         int `json:"totalJobs"`
}

type JobList struct {
	FinishedJobs []*Job  `json:"finishedJobs"`
	Meta         JobMeta `json:"meta"`
	QueuedJobs   []*Job  `json:"queuedJobs"`
	RunningJobs  []*Job  `json:"runningJobs"`
}

type NodeStatus struct {
	GPUType     string         `json:"gpuType"`
	Allocatable map[string]int `json:"gpu_allocatable"`
	Capacity    map[string]int `json:"gpu_capacity"`
}

// 接口：apis/GetVC?userName=&vcName=platform
type VcInfo struct {
	DeviceAvail map[string]int `json:"gpu_avaliable"`
	Nodes       []NodeStatus   `json:"node_status"`
}

type JobId struct {
	Id string `json:"jobId"`
}

type UriJobId struct {
	Id string `uri:"id" binding:"required"`
}

type TemplateId struct {
	Id int64 `uri:"id" binding:"required"`
}

type JobLog struct {
	Cursor string `json:"cursor"`
	Log    string `json:"log"`
}

// 创建endpoint
type CreateEndpointsReq struct {
	Endpoints []string `json:"endpoints"`
	JobId     string   `json:"jobId"`
}

// 返回值
type CreateEndpointsRsp struct {
}

type Endpoint struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Domain   string `json:"domain"`
	NodeName string `json:"nodeName,omitempty"`
	Port     string `json:"port"`
}

// 查询endpoints信息，返回
type GetEndpointsRsp struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type EndpointWrapper struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	AccessPoint string `json:"accessPoint"`
}

// 升级平台版本需要的信息
var UPGRADE_FILE_PATH = "/data/DLTSUpgrade"
var UPGRADE_CONFIG_FILE = "version.yaml"
var Upgrade_Progress = -1
var Log_Line_Point = 0
