package services

import (
	"github.com/apulis/AIArtsBackend/models"
)

const (
	MODELSET_STATUS_NORMAL   = "normal"
	MODELSET_STATUS_DELETING = "deleting"
)

func ListModelSets(page, count int,username string) ([]models.Modelset, int, error) {

	offset := count * (page - 1)
	limit := count
	return models.ListModelSets(offset, limit,username)
}

func ListModelSetsByName(page, count int, name string) ([]models.Modelset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListModelSetsByName(offset, limit, name)
}

func CreateModelset(name, description, creator, version, path, jobId string) error {
	size, err := GetDirSize(path)
	if err != nil {
		return err
	}
	modelset := models.Modelset{
		Name:        name,
		Description: description,
		Creator:     creator,
		Version:     version,
		Path:        path,
		Size:        size,
		JobId:       jobId,
		Status:      MODELSET_STATUS_NORMAL,
	}
	return models.CreateModelset(modelset)
}

func UpdateModelset(id int, description string) error {
	modelset, err := models.GetModelsetById(id)
	if err != nil {
		return err
	}
	modelset.Description = description
	return models.UpdateModelset(&modelset)
}

func GetModelset(id int) (models.Modelset, error) {
	return models.GetModelsetById(id)
}

func DeleteModelset(id int) error {
	modelset, err := models.GetModelsetById(id)
	if err != nil {
		return err
	}

	modelset.Status = MODELSET_STATUS_DELETING
	err = models.UpdateModelset(&modelset)
	if err != nil {
		return err
	}

	//err = os.RemoveAll(modelset.Path)
	//if err != nil {
	//	return err
	//}
	return models.DeleteModelset(&modelset)
}
