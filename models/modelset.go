package models

import (
	"fmt"
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
	Path        string `gorm:"type:text" json:"path"`
	Status      string `json:"status"`
	Size        int    `json:"size"`
	//模型类型 计算机视觉
	Use        string `json:"use"`
	JobId      string `json:"jobId"`
	DataFormat string `json:"dataFormat"`
	Arguments  string `gorm:"type:text" json:"arguments,omitempty"`
	EngineType string `json:"engineType"`
	Precision  string `json:"precision"`
	IsAdvance  bool   `json:"isAdvance"`
}

func ListModelSets(offset, limit int, isAdvance bool, name, status, username string) ([]Modelset, int, error) {
	var modelsets []Modelset
	total := 0
	whereQueryStr := fmt.Sprintf("creator='%s' and is_advance = '%s' ", username, isAdvance)
	if name != "" {
		whereQueryStr += fmt.Sprintf("and name='%s' ", name)
	}
	if status != "" {
		whereQueryStr += fmt.Sprintf("and status='%s' ", status)
	}

	res := db.Debug().Offset(offset).Limit(limit).Order("created_at desc").Where(whereQueryStr).Find(&modelsets)

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
