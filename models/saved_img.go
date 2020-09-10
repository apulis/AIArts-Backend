package models

import "fmt"

type SavedImage struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`

	Name        string `gorm:"not null" json:"name"`
	Version     string `gorm:"not null" json:"version"`
	Description string `gorm:"type:text" json:"description"`
	Creator     string `gorm:"not null" json:"creator"`
	FullName    string `gorm:"not null" json:"fullName"`

	IsPrivate   bool   `gorm:"not null" json:"isPrivate"`
	ContaindrId string `gorm:"not null" json:"containerId"`
	ImageId     string `gorm:"not null" json:"imageId"`
}

func ListSavedImages(offset, limit int, orderBy, order, name, username string) ([]SavedImage, int, error) {
	var savedImages []SavedImage
	total := 0

	whereQueryStr := fmt.Sprintf("creator = '%s' ", username)
	orQueryStr := fmt.Sprintf("is_private = 0 ")
	orderQueryStr := fmt.Sprintf("%s %s ", CamelToCase(orderBy), order)

	if name != "" {
		whereQueryStr += "and name like '%" + name + "%' "
		orQueryStr += "and name like '%" + name + "%'"
	}

	res := db.Offset(offset).Limit(limit).Order(orderQueryStr).Where(whereQueryStr).Or(orQueryStr).Find(&savedImages)
	if res.Error != nil {
		return savedImages, total, res.Error
	}
	db.Model(&SavedImage{}).Where(whereQueryStr).Or(orQueryStr).Count(&total)
	return savedImages, total, nil
}

func GetSavedImage(id int) (SavedImage, error) {
	savedImage := SavedImage{ID: id}
	res := db.First(&savedImage)
	if res.Error != nil {
		return savedImage, res.Error
	}
	return savedImage, nil
}

func CreateSavedImage(savedImage SavedImage) error {
	return db.Create(&savedImage).Error
}

func UpdateSavedImage(savedImage *SavedImage) error {
	res := db.Save(savedImage)
	return res.Error
}

func DeleteSavedImage(savedImage *SavedImage) error {
	res := db.Delete(&savedImage)
	return res.Error
}
