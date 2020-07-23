package models

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
	Type        string `json:"type"`
	JobId       string `json:"jobId"`
}

func ListModelSets(offset, limit int,username string) ([]Modelset, int, error) {
	var modelsets []Modelset
	db.Find(&modelsets)
	total := 0
	res := db.Offset(offset).Limit(limit).Order("created_at desc").Where(&Modelset{Creator: username}).Find(&modelsets)
	if res.Error != nil {
		return modelsets, total, res.Error
	}

	db.Model(&Modelset{}).Where(&Modelset{Creator: username}).Count(&total)
	return modelsets, total, nil
}

func ListModelSetsByName(offset, limit int, name ,username string) ([]Modelset, int, error) {
	var modelsets []Modelset

	total := 0
	res := db.Offset(offset).Limit(limit).Order("created_at desc").Where(&Modelset{Name: name,Creator: username}).Find(&modelsets)
	if res.Error != nil {
		return modelsets, total, res.Error
	}

	db.Model(&Modelset{}).Where(&Modelset{Name: name,Creator: username}).Count(&total)
	return modelsets, total, nil
}

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
