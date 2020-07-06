package services

import (
	"github.com/apulis/AIArtsBackend/models"
)

func ListDatasets() []models.Dataset {
	return models.ListDatasets()
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
