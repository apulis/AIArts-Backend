package services

import (
	"fmt"

	"github.com/apulis/AIArtsBackend/models"
)

const (
	DATASET_STATUS_NORMAL   = "normal"
	DATASET_STATUS_OCCUPIED = "occupied"
	DATASET_STATUS_DELETING = "deleting"
)

func ListDatasets(page, count int, name, status, username string) ([]models.Dataset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListDatasets(offset, limit, name, status, username)
}
func ListDatasetsByName(page, count int, name, username string) ([]models.Dataset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListDataSetsByName(offset, limit, name, username)
}
func CreateDataset(name, description, creator, version, path string, isPrivate bool) error {
	size, err := GetDirSize(path)
	if err != nil {
		return err
	}
	dataset := models.Dataset{
		Name:        name,
		Description: description,
		Creator:     creator,
		Version:     version,
		Path:        path,
		Size:        size,
		IsPrivate:   isPrivate,
		Status:      DATASET_STATUS_NORMAL,
	}
	return models.CreateDataset(dataset)
}

func UpdateDataset(id int, description string) error {
	dataset, err := models.GetDatasetById(id)
	if err != nil {
		return err
	}
	dataset.Description = description
	return models.UpdateDataset(&dataset)
}

func GetDataset(id int) (models.Dataset, error) {
	return models.GetDatasetById(id)
}

func DeleteDataset(id int) error {
	dataset, err := models.GetDatasetById(id)
	if err != nil {
		return err
	}
	dataset.Status = DATASET_STATUS_DELETING
	if dataset.Binds == "" {
		err = models.DeleteDataset(&dataset)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("still using")

	}
	//err = os.RemoveAll(dataset.Path)
	//if err != nil {
	//	return err
	//}
	return nil

}
func BindDataset(id int, platform, pid string) error {
	err := models.BindDatasetById(id, platform, pid)
	if err != nil {
		return err
	}
	return nil
}
func UnbindDataset(id int, platform, pid string) error {
	err := models.UnbindDatasetById(id, platform, pid)
	if err != nil {
		return err
	}
	return nil
}
