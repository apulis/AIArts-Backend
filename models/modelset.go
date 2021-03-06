package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Modelset struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Creator     string `gorm:"not null" json:"creator"`
	Version     string `gorm:"not null" json:"version"`
	Status      string `json:"status"`
	Size        int64  `gorm:"type bigint(20)" json:"size"`
	//模型类型 图像分类
	Use        string `json:"use"`
	JobId      string `json:"jobId"`
	DataFormat string `json:"dataFormat"`
	//Dataset     string `json:"dataset"`
	DatasetName string `json:"datasetName"`
	DatasetPath string `json:"datasetPath"`
	//omitempty 值为空，不编码
	Params    *ParamsItem `gorm:"type:text" json:"params"`
	Engine    string      `json:"engine"`
	Precision string      `json:"precision"`
	IsAdvance bool        `json:"isAdvance"`
	//模型路径
	CodePath string `json:"codePath"`
	//指定的模型参数路径
	ParamPath string `json:"paramPath"`
	// 输出文件路径
	OutputPath string `json:"outputPath"`
	//启动文件路径
	StartupFile string `json:"startupFile"`
	//模型路径
	VisualPath string `json:"visualPath"`
	//评估训练任务id
	EvaluationId string `json:"evaluationId"`
	// 评估设备类型
	DeviceType string `json:"deviceType"`
	DeviceNum  int    `json:"deviceNum"`
}

type ParamsItem map[string]string

func ListModelSets(offset, limit int, orderBy, order string, isAdvance bool, name, status, use, username string) ([]Modelset, int, error) {
	var modelsets []Modelset
	total := 0

	whereQueryStr := fmt.Sprintf("creator='%s' and is_advance = 0 ", username)
	if isAdvance {
		whereQueryStr = fmt.Sprintf(" is_advance = 1 ")
	}
	if name != "" {
		whereQueryStr += "and name like '%" + name + "%' "
	}
	if status != "" {
		whereQueryStr += fmt.Sprintf("and status='%s' ", status)
	}
	if strings.HasPrefix(use,`Avisualis`) {
		whereQueryStr += "and `use` like '" + use + "%' "
	} else {
		whereQueryStr += "and `use` not like 'Avisualis%' "
	}
	orderQueryStr := fmt.Sprintf("%s %s ", CamelToCase(orderBy), order)
	res := db.Offset(offset).Limit(limit).Order(orderQueryStr).Where(whereQueryStr).Find(&modelsets)

	if res.Error != nil {
		return modelsets, total, res.Error
	}
	db.Model(&Modelset{}).Where(whereQueryStr).Count(&total)
	return modelsets, total, nil
}

//alter table modelsets add column use  varchar(255) ;
//alter table modelsets add column precision varchar(255);alter table modelsets add column is_advance  varchar(255) ;
func GetModelsetById(id int) (Modelset, error) {
	modelset := Modelset{ID: id}
	res := db.First(&modelset)
	if res.Error != nil {
		return modelset, res.Error
	}
	return modelset, nil
}
func GetModelsetByName(name string) (Modelset, error) {
	var modelset Modelset
	whereQueryStr := fmt.Sprintf(" name = '%s' ", name)
	res := db.Where(whereQueryStr).Find(&modelset)
	if res.Error != nil {
		return modelset, res.Error
	}
	return modelset, nil
}

func CreateModelset(modelset Modelset) error {
	return db.Create(&modelset).Error
}

func UpdateModelset(modelset *Modelset) error {
	res := db.Save(modelset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteModelset(modelset *Modelset) error {
	res := db.Delete(&modelset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (this *ParamsItem) Value() (driver.Value, error) {
	binData, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return string(binData), nil
}

func (this *ParamsItem) Scan(v interface{}) error {
	switch t := v.(type) {
	case string:
		if t != "" {
			err := json.Unmarshal([]byte(t), this)
			if err != nil {
				return err
			}
		}
	case []byte:
		if len(t) != 0 {
			err := json.Unmarshal(t, this)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("无法将[%v] 反序列化为Modelset类型", reflect.TypeOf(v).Name())
	}

	return nil
}
