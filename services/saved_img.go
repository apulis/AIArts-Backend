package services

import "github.com/apulis/AIArtsBackend/models"

func ListSavedImages(page, count int, orderBy, order, name, username string) ([]models.SavedImage, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListSavedImages(offset, limit, orderBy, order, name, username)
}

func CreateSavedImage(name, version, description, jobId string, isPrivate bool) error {
	return nil
}

func UpdateSavedImage(id int, description string) error {
	savedImage, err := models.GetSavedImage(id)
	if err != nil {
		return err
	}
	savedImage.Description = description
	return models.UpdateSavedImage(&savedImage)
}

func GetSavedImage(id int) (models.SavedImage, error) {
	return models.GetSavedImage(id)
}

func DeleteSavedImage(id int) error {
	savedImage, err := models.GetSavedImage(id)
	if err != nil {
		return err
	}
	return models.DeleteSavedImage(&savedImage)
}
