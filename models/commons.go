package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/jinzhu/gorm"
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
	createTableIfNotExists(VisualJob{})
	createTableIfNotExists(SavedImage{})

	createTableIfNotExists(ExpProject{})
	createTableIfNotExists(Experiment{})

	db.Model(&Experiment{}).AddForeignKey("project_id", "exp_projects(id)", "RESTRICT", "RESTRICT")

	initVersionInfoTable()

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

// init version info of platform, make sure there is at least one record in versionInfo table
func initVersionInfoTable() {
	checkDataset := new(VersionInfoSet)
	err := db.Limit(1).Find(checkDataset).Error
	if err != gorm.ErrRecordNotFound {
		return
	}
	initVersion := VersionInfoSet{
		Description: "当前版本号：1.0.1",
		Version:     "1.0.1",
		Creator:     "admin",
	}
	err = UploadVersionInfoSet(initVersion)
	fmt.Println("upload ready")
	if err != nil {
		panic(err)
	}
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
	TemplatePublic        int = 1 // 公有
	TemplatePrivate       int = 2 // 用户私有
	TemplatePublicPrivate int = 3 // 公有（包括预置参数）和私有
	TemplatePredefined    int = 4 // 预置参数
)

// 以下结构体用于和DLTS平台交互
const (
	JobTypeTraining       string = "training"       // 老DLTS默认采用的jobType
	JobTypeArtsTraining   string = "artsTraining"   // 供电局项目：模型训练
	JobTypeArtsEvaluation string = "artsEvaluation" // 供电局项目：模型评估
	JobTypeCodeEnv        string = "codeEnv"        // 供电局项目：代码环境
	JobTypeVisualJob      string = "visualjob"      //供电局项目：可视化作业
	JobStatusAll          string = "all"
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

	ScriptParams map[string]string `json:"scriptParams"`
	JobGroup    string  `json:"jobGroup"`
	Track       int     `json:"track"`
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
	AllJobs      []*Job  `json:"allJobs"`
}

type NodeStatus struct {
	GPUType     string         `json:"gpuType"`
	Allocatable map[string]int `json:"gpu_allocatable"`
	Capacity    map[string]int `json:"gpu_capacity"`
	DeviceStr   string         `json:"deviceStr,omitempty"`
}

// 接口：apis/GetVC?userName=&vcName=platform
type VcInfo struct {
	DeviceAvail    map[string]int `json:"gpu_avaliable"`
	DeviceCapacity map[string]int `json:"gpu_capacity"`
	Nodes          []*NodeStatus  `json:"node_status"`
}

// 接口：apis/GetAllDevice?userName=
type DeviceItem2 struct {
	DeviceStr string `json:"deviceStr"`
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
	Cursor  string `json:"cursor,omitempty"`
	Log     string `json:"log,omitempty"`
	MaxPage int    `json:"maxPage"`
}

// 创建endpoint
type CreateEndpointsReq struct {
	Endpoints []string `json:"endpoints"`
	JobId     string   `json:"jobId"`
	Arguments string   `json:"arguments"` //启动endpoints的时候可能添加的命令行参数，在完成可视化作业需求时，因为需要更改tensorboard log路径是增加参数
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

/* Upgrade_Progress原是作为进度条百分比，现在作为升级过程的状态码，目前共有以下几种
* -1: not ready,系统刚进入时的准备状态（可以与success合并，但是仍未合并，作为健康态的表现）
* 0: upgrading,正在升级
* 100: success,升级完成，也是健康态
* 300: error,升级出错
 */
var Upgrade_Progress = -1
var Log_Line_Point = 0
