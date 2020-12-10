package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	DATASET_UPLOAD_FROM_WEB 		int = 1
	DATASET_UPLOAD_FROM_OTHER 		int = 2
)

type Dataset struct {
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
	//存储绑定信息
	//plantform#id*plantform_id
	Binds string `json:"binds"`
	//是否是公开数据集
	IsPrivate    bool  `json:"isPrivate"`
	IsTranslated bool  `json:"isTranslated"`
	Size         int64 `gorm:"type bigint(20)" json:"size"`
}

func ListDatasets(offset, limit int, orderBy, order, name, status string, isTranslated bool, username string) ([]Dataset, int, error) {
	var datasets []Dataset
	total := 0
	//先查询该用户的所有数据集，再查询公开数据集
	whereQueryStr := fmt.Sprintf("creator = '%s' ", username)
	orQueryStr := fmt.Sprintf("is_private=0 ")
	orderQueryStr := fmt.Sprintf("%s %s ", CamelToCase(orderBy), order)
	if isTranslated {
		whereQueryStr += "and is_translated = 1 "
		orQueryStr += "and is_translated = 1 "
	}
	if name != "" {
		whereQueryStr += "and name like '%" + name + "%' "
		orQueryStr += "and name like '%" + name + "%' "
	}

	if status != "" && status != "all" {
		whereQueryStr += fmt.Sprintf("and status='%s' ", status)
		orQueryStr += fmt.Sprintf("and status='%s' ", status)
	}
	res := db.Offset(offset).Limit(limit).Order(orderQueryStr).Where(whereQueryStr).
		Or(orQueryStr).Find(&datasets)

	if gin.Mode() == "debug" {
		res = db.Offset(offset).Limit(limit).Order(orderQueryStr).Where(whereQueryStr).
			Or(orQueryStr).Find(&datasets)
	}
	if res.Error != nil {
		return datasets, total, res.Error
	}
	db.Model(&Dataset{}).Where(whereQueryStr).Or(orQueryStr).Count(&total)
	return datasets, total, nil
}

func GetDatasetById(id int) (Dataset, error) {
	dataset := Dataset{ID: id}
	res := db.First(&dataset)
	if res.Error != nil {
		return dataset, res.Error
	}
	return dataset, nil
}
func DatasetIsExist(name, username string) error {
	var dataset Dataset
	whereQueryStr := fmt.Sprintf(" name = '%s' and (creator = '%s' or is_private = 0 ) ", name, username)
	res := db.Where(whereQueryStr).Find(&dataset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func ListDataSetsByName(offset, limit int, name, username string) ([]Dataset, int, error) {
	var datasets []Dataset
	total := 0
	//展示指定用户的
	res := db.Offset(offset).Limit(limit).Order("created_at desc").
		Where("name=? and is_private=?", name, false).Or("name=? and creator=?", name, username).Find(&datasets)
	if res.Error != nil {
		return datasets, total, res.Error
	}
	db.Model(&Dataset{}).Where("name=? and is_private=?", name, false).Or("name=? and is_private=?", name, false).Count(&total)

	return datasets, total, nil
}
func CreateDataset(dataset Dataset) error {
	return db.Create(&dataset).Error
}

func UpdateDataset(dataset *Dataset) error {
	res := db.Save(dataset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteDataset(dataset *Dataset) error {
	res := db.Delete(&dataset)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func BindDatasetById(id int, platform, pid string) error {
	bind := platform + "#" + pid
	dataset := Dataset{ID: id}
	res := db.Select("binds").Find(&dataset)
	bindsString := dataset.Binds
	if res.Error != nil {
		return nil
	}
	bindsArr := strings.Split(bindsString, "*")
	for _, b := range bindsArr {
		if b == bind {
			return fmt.Errorf("already bind")
		}
	}
	//初次绑定去掉开头的*号
	if bindsString == "" {
		bindsString = bind
	} else {
		bindsString = bindsString + "*" + bind
	}
	res = db.Model(&dataset).Update("binds", bindsString)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func UnbindDatasetById(id int, platform, pid string) error {
	bind := platform + "#" + pid
	dataset := Dataset{ID: id}
	res := db.Select("binds").Find(&dataset)
	bindsString := dataset.Binds
	if res.Error != nil {
		return nil
	}
	isBind := false
	bindsArr := strings.Split(bindsString, "*")
	for index, b := range bindsArr {
		if b == bind {
			bindsArr = append(bindsArr[:index], bindsArr[index+1:]...)
			isBind = true
		}
	}
	if !isBind {
		return fmt.Errorf("no bind")
	}
	bindsString = strings.Join(bindsArr, "*")
	res = db.Model(&dataset).Update("binds", bindsString)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
