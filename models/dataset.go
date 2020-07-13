package models

type Dataset struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Name        string `gorm:"unique_index;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Creator     string `gorm:"not null" json:"creator"`
	Version     string `gorm:"not null" json:"version"`
	Path        string `gorm:"type:text" json:"path"`
	Status      string `json:"status"`
	Size        int    `json:"size"`
}

func ListDatasets(offset, limit int) ([]Dataset, int, error) {
	var datasets []Dataset
	db.Find(&datasets)

	total := 0
	res := db.Offset(offset).Limit(limit).Order("created_at desc").Find(&datasets)
	if res.Error != nil {
		return datasets, total, res.Error
	}

	db.Model(&Dataset{}).Count(&total)
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