package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
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
	Size        int    `json:"size"`
	//模型类型 图像分类
	Use        string         `json:"use"`
	JobId      string         `json:"jobId"`
	DataFormat string         `json:"dataFormat"`
	Dataset    string         `json:"dataset,omitempty"`
	//omitempty 值为空，不编码
	Arguments  *ArgumentsItem `gorm:"type:text" json:"arguments,omitempty"`
	EngineType string         `json:"engineType"`
	Precision  string         `json:"precision"`
	IsAdvance  bool           `json:"isAdvance"`
	//模型路径
	ModelPath  string           `json:"modelPath"`
	//模型参数路径
	ArgumentPath  string           `json:"argumentPath"`


}

type ArgumentsItem map[string]string

func ListModelSets(offset, limit int,orderBy,order string, isAdvance bool, name, status, username string) ([]Modelset, int, error) {
	var modelsets []Modelset
	total := 0
	is_advance := 0
	if isAdvance {
		is_advance = 1
	}
	whereQueryStr := fmt.Sprintf("creator='%s' and is_advance = %d ", username, is_advance)
	if name != "" {
		whereQueryStr += fmt.Sprintf("and name='%s' ", name)
	}
	if status != "" {
		whereQueryStr += fmt.Sprintf("and status='%s' ", status)
	}

	orderQueryStr := fmt.Sprintf("%s %s ",orderBy, order)
	res := db.Debug().Offset(offset).Limit(limit).Order(orderQueryStr).Where(whereQueryStr).Find(&modelsets)

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

func (this *ArgumentsItem) Value() (driver.Value, error) {
	binData, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}
	return string(binData), nil
}

func (this *ArgumentsItem) Scan(v interface{}) error {
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
