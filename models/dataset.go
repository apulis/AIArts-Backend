package models

import (
	"github.com/jinzhu/gorm"
)

type Dataset struct {
	gorm.Model
	Name        string `gorm:"unique_index;not null"`
	Description string `gorm:"type:text"`
	Creator     string `gorm:"not null"`
	Version     string `gorm:"not null"`
	Path        string `gorm:"type:text"`
}
