package services

import (
	"github.com/apulis/AIArtsBackend/models"
)

func ListDatasets(page, count int) ([]models.Dataset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListDatasets(offset, limit)
}

func CreateDataset(name, description, creator, version, path string) error {
	dataset := models.Dataset{
		Name:        name,
		Description: description,
		Creator:     creator,
		Version:     version,
		Path:        path,
	}
	return models.CreateDataset(dataset)
}

func UpdateDataset(id int, description string) error {
	dataset, err := models.GetDataSetById(id)
	if err != nil {
		return err
	}
	dataset.Description = description
	return models.UpdateDataset(&dataset)
}

func GetDataset(id int) (models.Dataset, error) {
	return models.GetDataSetById(id)
}
