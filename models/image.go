package models

import (
	"database/sql/driver"
)

// 数据库结构
type Image struct {

	ID        int       	`gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  	`json:"createdAt"`
	UpdatedAt UnixTime  	`json:"updatedAt"`
	DeletedAt *UnixTime 	`json:"deletedAt"`

	// 镜像类别：
	// aiframework: tensorflow, mindspore
	// os: ubuntu, debian, etc.
	ImageType     string      `gorm:"not null" json:"imageType"`

	// 镜像全称
	ImageFullName string      `gorm:"not null" json:"imageFullName"`

	// 镜像参数
	Details       ImageParams `json:"details"`
}

// 数据库字段
type ImageParams struct {
	Desc          		 string  `json:"desc"`
	Category          	 string  `json:"category"`
	Brand	             string  `json:"brand"`
	CPUArchType          string  `json:"cpuArchType"`
	DeviceType           string  `json:"deviceType"`
}

func ListImages(offset, limit int) ([]Image, int, error) {

	var images []Image
	total := 0

	res := db.Offset(offset).Limit(limit).Order("created_at desc").Find(&images)
	if res.Error != nil {
		return images, total, res.Error
	}

	db.Model(&Image{}).Count(&total)
	return images, total, nil
}

// Scan Scanner
func (details*ImageParams) Scan(value interface{}) error {
	return scan(details, value)
}

// Value Valuer
func (details *ImageParams) Value() (driver.Value, error) {
	return value(details)
}