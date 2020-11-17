package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/configs"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
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

func checkDBError(err error)error{
	if err != nil && strings.Contains(err.Error(),"Duplicate entry") {
		return configs.NewAPIException(http.StatusBadRequest,configs.NAME_ALREADY_EXIST_CODE,err.Error())
	}else{
		return err
	}
}

//Experiment 实验管理实例,分组一批runs集合 .
type Experiment struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	ProjectID uint   `json:"projectId" gorm:"unique_index:proj_idx"`
	Name      string `json:"name" gorm:"unique_index:proj_idx;not null"`

	Description string       `json:"description" gorm:"type:text"`
	Creator     string       `json:"creator"`
	CreatedAt   UnixTime     `json:"createdAt"`
	UpdatedAt   UnixTime     `json:"updatedAt"`
	DeletedAt   *UnixTime    `json:"deletedAt,omitempty"`
	Meta        *JsonMetaData `json:"meta,omitempty" gorm:"type:text"`

	ProjectName        string  `json:"projectName" gorm:"-"`
	TrackId            string  `json:"trackId"     gorm:"-"`
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

func checkUserOrderBy(orderBy,order string) func (*gorm.DB) *gorm.DB{

	 return func(db*gorm.DB)*gorm.DB{
	 	 if len(orderBy) > 0 {
	 	 	return db.Order( fmt.Sprintf("%s %s",CamelToCase(orderBy),order) )
		 }
		 return db
	 }
}
func checkFilterName(name string)func (*gorm.DB)*gorm.DB{
	 return func(db*gorm.DB)*gorm.DB{
	 	if len(name) > 0 {
	 		return db.Where("name like ? ","%"+name+"%")
		}
		return db
	 }
}
func checkListUnScope(isAll uint) func(*gorm.DB)*gorm.DB{
	 return func(db*gorm.DB)*gorm.DB{
	 	switch isAll{
		case 1:   return db.Unscoped()
		case 2:   return db.Unscoped().Where("deleted_at is not null")
		default:  return db
		}
	 }
}



func ListAllExpProjects(offset, limit, isAll uint, name,orderBy, order string) ([]ExpProject, uint, error) {

	var datasets []ExpProject
	var total uint

	err := db.Scopes(checkListUnScope(isAll),checkUserOrderBy(orderBy,order),checkFilterName(name)).
		    Offset(offset).Limit(limit).Select(list_projects_fields).Find(&datasets).Error
	if err == nil{
		err = db.Model(&ExpProject{}).Scopes(checkListUnScope(isAll),checkFilterName(name)).Count(&total).Error
	}

	return datasets,total,err
}

func CreateExpProject(project *ExpProject) error {
	 return  checkDBError(db.Create(project).Error)
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

func ListAllExperiments(projectID ,offset, limit, isAll uint,name, orderBy, order string) ([]Experiment, uint, error) {

	var datasets []Experiment
	var total uint

	err := db.Scopes(checkListUnScope(isAll),checkUserOrderBy(orderBy,order),checkFilterName(name)).
		         Offset(offset).Limit(limit).Select(list_experiments_fields).
		        Where("project_id=?",projectID).Find(&datasets).Error
	if err == nil{
		err = db.Model(&Experiment{}).Scopes(checkListUnScope(isAll),checkFilterName(name)).
			    Where("project_id=?",projectID).Count(&total).Error
	}

	return datasets,total,err
}
func CreateExperiment(experiment *Experiment) error {
	return checkDBError(db.Create(experiment).Error)
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
	return db.Raw(`select experiments.*,exp_projects.name as project_name from experiments left join exp_projects on
               experiments.project_id=exp_projects.id where experiments.id=` + strconv.Itoa(int(id))).
		        Scan(experiment).Error
}
