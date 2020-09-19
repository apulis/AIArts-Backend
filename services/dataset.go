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

func ListDatasets(page, count int, orderBy, order, name, status string,isTranslated bool , username string) ([]models.Dataset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListDatasets(offset, limit, orderBy, order, name, status,isTranslated, username)
}
func ListDatasetsByName(page, count int, name, username string) ([]models.Dataset, int, error) {
	offset := count * (page - 1)
	limit := count
	return models.ListDataSetsByName(offset, limit, name, username)
}
func AppendAnnoDataset(datasets []models.Dataset,total,pageNum, pageSize int, orderBy, order string)(string,int,[]models.Dataset){
	var annoDatasets []models.DataSet
	queryStringParameters := models.QueryStringParametersV2{
		PageNum:  pageNum,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Order:    order,
	}
	annoDatasets, _, err := ListAllDatasets(queryStringParameters)
	message := "success"
	if err != nil {
		message = "label image platform is error"
		//return AppError(FAILED_FETCH_ANNOTATION_CODE, "label plantform is error")
	} else {
		for _, v := range annoDatasets {
			if v.ConvertStatus == "finished" {
				dataset := models.Dataset{
					Name:        v.Name,
					Description: v.Info,
					Path:        v.ConvertOutPath,
					Status:      v.Name,
					//是否是公开数据集
					IsPrivate:    v.IsPrivate,
					IsTranslated: true,
				}
				datasets = append(datasets, dataset)
				total += 1
			}
		}
	}
	return message,total,datasets
}
func CreateDataset(name, description, creator, version, path string, isPrivate bool,isTranslated bool) error {
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
		IsTranslated:   isTranslated,
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
