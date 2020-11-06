package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

//ExpProject AIArts项目表示,类似于文件夹概念 .
type ExpProject struct {
	ID   uint   `json:"id"   gorm:"primary_key;auto_increment"`
	Name string `json:"name" gorm:"unique;not null"`

	Description string    `json:"description" gorm:"type:text"`
	Creator     string    `json:"creator"`
	CreatedAt   UnixTime  `json:"createdAt"`
	UpdatedAt   UnixTime  `json:"updatedAt"`
	DeletedAt   *UnixTime `json:"deletedAt,omitempty"`
	UserGroup   string    `json:"userGroup"`
}

type JsonMetaData struct{
	//data map[string]interface{}
	data_str string
}
func (d*JsonMetaData)MarshalJSON()([]byte,error){
	if len(d.data_str) == 0 {
		return []byte("null"),nil
	}
	return []byte(d.data_str),nil
}
func (d*JsonMetaData)UnmarshalJSON(b[]byte)error{
	if len(b) >= 2 && b[0] == '{'{
		d.data_str=string(b)
	}else{
		d.data_str=""
	}
	return nil
}

func (d JsonMetaData) Value() (driver.Value, error) {
	if len(d.data_str) == 0 {
		return nil,nil
	}else{
		return d.data_str,nil
	}
}

func (d *JsonMetaData) Scan(v interface{}) error {
	switch ty:= v.(type){
	case string: d.data_str=ty
	case []byte: d.data_str=string(ty)
	default:     d.data_str=""
	}
	return nil
}
//Experiment 实验管理实例,分组一批runs集合 .
type Experiment struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	ProjectID uint   `json:"projectID" gorm:"unique_index:proj_idx"`
	Name      string `json:"name" gorm:"unique_index:proj_idx;not null"`

	Description string       `json:"description" gorm:"type:text"`
	Creator     string       `json:"creator"`
	CreatedAt   UnixTime     `json:"createdAt"`
	UpdatedAt   UnixTime     `json:"updatedAt"`
	DeletedAt   *UnixTime    `json:"deletedAt,omitempty"`
	Meta        *JsonMetaData `json:"meta,omitempty" gorm:"type:text"`
}

type RequestUpdates map[string]interface{}
type RespRowData    map[string]interface{}

var list_projects_fields    =  "name,id,description,creator,created_at,updated_at,deleted_at,user_group"
var list_experiments_fields =  "name,id,description,creator,created_at,updated_at,deleted_at"

func (reqMap*RequestUpdates)TranslateJsonMeta(args...string)int{
	cnt := 0
	for _,key := range(args){
         value,ok   := (*reqMap)[key]
         if ok {
         	switch ty := value.(type){
			case map[string]interface{}:
				  json_str,_ :=json.Marshal(ty)
				  (*reqMap)[key] = JsonMetaData{string(json_str)}
				  cnt++
			case []byte: if len(ty) >= 2 && ty[0] == '{' {
				  (*reqMap)[key] = JsonMetaData{string(ty)}
				  cnt++
			}
			case string:
				  if len(ty) >= 2 && ty[0] == '{'{
				   (*reqMap)[key] = JsonMetaData{ty}
				   cnt++
				  }
			}
		 }
	}
    return cnt
}


func wrapDBOpUpdate(db * gorm.DB , changes int64 ) error{
	if db.Error!=nil {
		return db.Error
	}
    if changes > 0 && db.RowsAffected != changes {
    	return errors.New("unexpected update DB")
	}
	return nil
}

func ListAllExpProjects(offset, limit, isAll uint, orderBy, order string) ([]ExpProject, uint, error) {

	var datasets []ExpProject
	var total uint
	var orderQueryStr string
	if len(orderBy) != 0 {
		orderQueryStr=fmt.Sprintf("order by %s %s ", CamelToCase(orderBy), order)
	}
	var res *gorm.DB
	if isAll == 0 {
		res = db.Offset(offset).Limit(limit).Order(orderQueryStr).Select(list_projects_fields).Find(&datasets)
		if res.Error != nil {
			return datasets, total, res.Error
		}
		db.Model(&ExpProject{}).Count(&total)
	}else{
        res = db.Unscoped().Offset(offset).Limit(limit).Order(orderQueryStr).Select(list_projects_fields).Find(&datasets)
		if res.Error != nil {
			return datasets, total, res.Error
		}
		db.Unscoped().Model(&ExpProject{}).Count(&total)
	}

	return datasets, total, nil
}

func CreateExpProject(project *ExpProject) error {
     return db.Create(project).Error
}

func QueryExpProject(id uint , project*ExpProject) error{
	project.ID=id
	return db.Unscoped().First(project).Error;
}

func RenameExpProject(id uint,name string)error{
	project := ExpProject{ID:id}
	return wrapDBOpUpdate(db.Unscoped().Model(&project).Update("name",name),1)
}

func UpdateExpProject(id uint,project*RequestUpdates)error{
    return wrapDBOpUpdate(db.Model(&ExpProject{ID:id}).
    	Select("name","description","user_group").
    	Updates(*project),1)
}

func HideExpProject(id uint)error{
    return wrapDBOpUpdate(db.Delete(&ExpProject{ID:id}),1)
}

func RestoreExpProject(id uint)error{
    return wrapDBOpUpdate(db.Unscoped().Model(&ExpProject{ID:id}).Update("deleted_at",nil),1)
}

func ListAllExperiments(projectID ,offset, limit, isAll uint, orderBy, order string) ([]Experiment, uint, error) {

	var datasets []Experiment
	var total uint
	var orderQueryStr string
	if len(orderBy) != 0 {
		orderQueryStr=fmt.Sprintf("order by %s %s ", CamelToCase(orderBy), order)
	}
	var res *gorm.DB
	if isAll == 0 {
		res = db.Debug().Offset(offset).Limit(limit).Order(orderQueryStr).Select(list_experiments_fields).
			Find(&datasets,"project_id=?",projectID)
		if res.Error != nil {
			return datasets, total, res.Error
		}
		db.Model(&Experiment{}).Where("project_id=?",projectID).Count(&total)
	}else{
		res = db.Unscoped().Offset(offset).Limit(limit).Order(orderQueryStr).Select(list_experiments_fields).
			Find(&datasets,"project_id=?",projectID)
		if res.Error != nil {
			return datasets, total, res.Error
		}
		db.Unscoped().Model(&Experiment{}).Where("project_id=?",projectID).Count(&total)
	}

	return datasets, total, nil
}
func CreateExperiment(experiment *Experiment) error {
	return db.Create(experiment).Error
}

func RenameExperiment(id uint,name string)error{
	experiment := Experiment{ID:id}
	return wrapDBOpUpdate(db.Unscoped().Model(&experiment).Update("name",name),1)
}

func UpdateExperiment(id uint,experiment*RequestUpdates)error{
	experiment.TranslateJsonMeta("meta")
	return wrapDBOpUpdate(db.Model(&Experiment{ID:id}).
		Select("name","description","meta").
		Updates(*experiment),1)
}
func HideExperiment(id uint)error{
	return wrapDBOpUpdate(db.Delete(&Experiment{ID:id}),1)
}

func RestoreExperiment(id uint)error{
	return wrapDBOpUpdate(db.Unscoped().Model(&Experiment{ID:id}).Update("deleted_at",nil),1)
}
func QueryExperiment(id uint , experiment*Experiment) error{
	experiment.ID=id
	return db.Unscoped().First(experiment).Error;
}