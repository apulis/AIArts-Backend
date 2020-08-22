package models

type VisualJob struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	Name        string `json:"name"`
	Status      string `json:status`
	LogPath     string `json:"logPath"`
	Description string `json:"description"`
	RelateJobId string `json:relateJobId`
	CreateTime  string `json:"createTime"`
}

func CreateVisualJob(visualJob VisualJob) error {
	return db.Create(&visualJob).Error
}

func GetVisualJobBackroundJobId(visualJobId int) (string, error) {
	var visualJob VisualJob
	res := db.Where("id = ?", visualJobId).Find(&visualJob)
	if res.Error != nil {
		return "", res.Error
	}
	return visualJob.RelateJobId, nil
}

func GetVisualJobById(Id int) (VisualJob, error) {
	var visualJob VisualJob
	res := db.Where("id = ?", Id).Find(&visualJob)
	if res.Error != nil {
		return visualJob, res.Error
	}
	return visualJob, nil
}

func GetAllVisualJobByArguments(userName string, pageNum int, pageSize int, orderBy string, status string, jobName string, order string) ([]VisualJob, error) {
	var visualJobList []VisualJob
	temp := db.Where("userName =?", userName).Limit(pageSize).Offset((pageNum - 1) * pageSize)
	if orderBy != "" && order != "" {
		temp = temp.Order(orderBy + " " + order)
	}
	if jobName != "" {
		temp = temp.Where("jobName LIKE ?", jobName+"%")
	}
	if status != "" {
		temp = temp.Where("status =?", status)
	}
	res := temp.Find(&visualJobList)
	if res.Error != nil {
		return nil, res.Error
	}
	return visualJobList, nil
}

func GetVisualJobsSumCount() (int, error) {
	var count int
	res := db.Table("visualJob").Count(&count)
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
