package models

import (
	"fmt"
)
type VisualJob struct {
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`
	UserName    string `gorm: "userName" json:"username"`
	Name        string `json:"name"`
	Status      string `json:status`
	LogPath     string `gorm: "logPath" json:"logPath"`
	Description string `gorm:"type:text" json:"description"`
	RelateJobId string `gorm: "relateJobId" json:relateJobId`
}

func CreateVisualJob(visualJob VisualJob) error {
	return db.Create(&visualJob).Error
}

func GetVisualJobBackroundJobId(visualJobId int) (string, error) {
	var visualJob VisualJob
	res := db.First(&visualJob,visualJob)
	if res.Error != nil {
		return "", res.Error
	}
	return visualJob.RelateJobId, nil
}

func GetVisualJobById(Id int) (VisualJob, error) {
	var visualJob VisualJob
	res := db.First(&visualJob,Id)
	if res.Error != nil {
		return visualJob, res.Error
	}
	return visualJob, nil
}

func GetAllVisualJobByArguments(userName string, pageNum int, pageSize int, status string, jobName string, order string, orderBy string) ([]VisualJob, error) {
	var visualJobList []VisualJob
	temp := db.Where("user_name =?", userName)
	if orderBy != "" && order != "" {
		fmt.Println("search order %s",order)
		temp = temp.Order(orderBy + " " + order)
	}
	if jobName != "" {
		fmt.Println("search jobName %s",jobName)
		temp = temp.Where("name LIKE ?", jobName+"%")
	}
	if status != "" {
		fmt.Println("search status %s",status)
		temp = temp.Where("status =?", status)
	}
	res := temp.Find(&visualJobList)
	if res.Error != nil {
		return nil, res.Error
	}
	return visualJobList, nil
}

func GetVisualJobsSumCount(userName string ) (int, error) {
	var count int
	res := db.Table("visual_jobs").Where("deleted_at is NULL").Where("userName = ?",userName).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

func UpdateVisualJob(job *VisualJob) error {
	res := db.Save(job)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func DeleteVisualJob(job *VisualJob) error {
	return db.Delete(&job).Error
}
