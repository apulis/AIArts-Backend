package models

import "time"

type Dataset struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Name        string     `gorm:"unique_index;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Creator     string     `gorm:"not null" json:"creator"`
	Version     string     `gorm:"not null" json:"version"`
	Path        string     `gorm:"type:text" json:"path"`
}

func ListDatasets() []Dataset {
	var datasets []Dataset
	db.Find(&datasets)
	return datasets
}

func CreateDataset(dataset Dataset) error {
	return db.Create(&dataset).Error
}
