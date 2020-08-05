package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	PublicTemplate  = 1
	PrivateTemplate = 2
)

// 此结构体用于前后端通信
type TemplateMeta struct {
	Id      int    `json:"id,omitempty"`
	Name    string `json:"name"`
	Scope   int    `json:"scope"`
	JobType string `json:"jobType"`
	Creator string `json:"creator"`

	CreatedAt UnixTime `json:"createdAt"`
	UpdatedAt UnixTime `json:"updatedAt"`
}

type TemplateParams struct {
	Name        string            `json:"name"`
	Engine      string            `json:"engine"`
	DeviceType  string            `json:"deviceType"`
	DeviceNum   int               `json:"deviceNum"`
	CodePath    string            `json:"codePath,omitempty"`
	StartupFile string            `json:"startupFile"`
	OutputPath  string            `json:"outputPath,omitempty"`
	DatasetPath string            `json:"datasetPath"`
	Params      map[string]string `json:"params,omitempty"`
	Desc        string            `json:"desc,omitempty"`
	CreateTime  string            `json:"createTime,omitempty"`
}

type TemplateItem struct {
	MetaData TemplateMeta   `json:"metaData"`
	Params   TemplateParams `json:"params"`
}

// 与数据库一一对应
type Templates struct {
	ID int `gorm:"primary_key" json:"id"`

	Name    string         `gorm:"not null" json:"name"`
	Scope   int            `gorm:"not null" json:"scope"`
	Data    TemplateParams `gorm:"type:text;not null" json:"data"` // TemplateParams转换为json的结果
	JobType string         `gorm:"not null" json:"jobType"`
	Creator string         `json:"creator"`

	CreatedAt UnixTime  `json:"createdAt"`
	UpdatedAt UnixTime  `json:"updatedAt"`
	DeletedAt *UnixTime `json:"deletedAt"`
}

var (
	escapedScopePattern = regexp.MustCompile("^\\s*[-_\\w\\d\\.`]+\\s*$")

	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// MySql INSERT ******
func Insert(db *sql.DB, table string, data map[string]interface{}) (lastInsertId int64, err error) {

	if len(data) == 0 {
		return 0, errors.New("empty data")
	}

	if !escapedScopePattern.MatchString(table) {
		return 0, fmt.Errorf("invalid table: %s", table)
	}

	idx, size := 0, len(data)
	columns, placeholder, args := make([]string, size), make([]string, size), make([]interface{}, size)

	for key, val := range data {
		if !escapedScopePattern.MatchString(key) {
			return 0, fmt.Errorf("invalid column: %s", key)
		}
		columns[idx] = quote(ToSnakeCase(key))
		placeholder[idx] = "?"
		args[idx] = val
		idx++
	}

	query := fmt.Sprintf("INSERT INTO %s(%s) values(%s)", quote(table), strings.Join(columns, ","), strings.Join(placeholder, ","))
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// MySql UPDATE ****** WHERE id = ?
func Update(db *sql.DB, table string, id int64, data map[string]interface{}) (rowsAffected int64, err error) {
	if len(data) == 0 {
		return 0, nil
	}

	if !escapedScopePattern.MatchString(table) {
		return 0, fmt.Errorf("invalid table: %s", table)
	}

	idx, size := 0, len(data)
	querySegs, args := make([]string, size), make([]interface{}, size)
	for key, val := range data {
		if !escapedScopePattern.MatchString(key) {
			return 0, fmt.Errorf("invalid column: %s", key)
		}
		querySegs[idx] = quote(ToSnakeCase(key)) + "=?"
		args[idx] = val
		idx++
	}

	query := fmt.Sprintf("UPDATE %s set %s WHERE id=%d", quote(table), strings.Join(querySegs, ","), id)
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

type TemplateProvider struct {
	gormDb *gorm.DB
}

func NewTemplateProvider(gormDb *gorm.DB) *TemplateProvider {
	return &TemplateProvider{gormDb: gormDb}
}

func (this *TemplateProvider) TableName() string { return "ai_arts.templates" }

func (this *TemplateProvider) GetDB() *gorm.DB { return this.gormDb }

// 分页查询
func (this *TemplateProvider) FindPage(order string, offset, limit int, query string, args ...interface{}) ([]*Templates, error) {

	var tmp []*Templates
	db := this.gormDb.Offset(offset).Limit(limit)

	if query != "" {
		db = db.Where(query, args...)
	}

	if order != "" {
		db = db.Order(order)
	}

	if err := db.Find(&tmp).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tmp, nil
}

func (this *TemplateProvider) Insert(data map[string]interface{}) (lastInsertId int64, err error) {
	n, err := Insert(this.gormDb.DB(), this.TableName(), data)
	return int64(n), err
}

func (this *TemplateProvider) Update(id int64, data map[string]interface{}) (rowsAffected int64, err error) {
	return Update(this.gormDb.DB(), this.TableName(), int64(id), data)
}

func (this *TemplateProvider) FindById(id int64) (*Templates, error) {

	tmp := &Templates{}
	db := this.gormDb.Where("id=?", id).First(&tmp)

	if db.RecordNotFound() {
		return nil, nil
	} else if db.Error != nil {
		return nil, db.Error
	} else {
		return tmp, nil
	}
}

// 字符转换
func quote(value string) string {

	value = strings.TrimSpace(value)
	if value[0] == '`' {
		return value
	}

	return "`" + strings.Replace(value, ".", "`.`", -1) + "`"
}

// 数据库记录
func (this *Templates) Load(scope int, creator, jobType string, item TemplateParams, isNew bool) {

	this.Scope = scope
	this.JobType = jobType

	this.Name = item.Name
	this.Creator = creator
	this.Data = item

	if isNew {
		this.CreatedAt = UnixTime{
			Time: time.Now(),
		}
		this.UpdatedAt = this.CreatedAt
	} else {
		this.UpdatedAt = UnixTime{
			Time: time.Now(),
		}
	}
}

func (this *Templates) ToMap() map[string]interface{} {

	data := make(map[string]interface{})
	data["name"] = this.Name
	data["scope"] = this.Scope

	binData, err := json.Marshal(this.Data)
	if err == nil {
		data["data"] = string(binData)
	}

	data["jobType"] = this.JobType
	data["creator"] = this.Creator

	if !this.CreatedAt.IsZero() {
		data["createdAt"] = this.CreatedAt
	}

	if !this.UpdatedAt.IsZero() {
		data["updatedAt"] = this.UpdatedAt
	}

	return data
}

func (this *Templates) ToTemplateItem() *TemplateItem {

	item := &TemplateItem{
		MetaData: TemplateMeta{
			Id:        this.ID,
			Name:      this.Name,
			Scope:     this.Scope,
			JobType:   this.JobType,
			Creator:   this.Creator,
			CreatedAt: this.CreatedAt,
			UpdatedAt: this.UpdatedAt,
		},

		Params: this.Data,
	}

	return item
}

func (this *TemplateParams) Value() (driver.Value, error) {

	binData, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	return string(binData), nil
}

func (this *TemplateParams) Scan(v interface{}) error {

	switch t := v.(type) {
	case string:
		if t != "" {
			err := json.Unmarshal([]byte(t), this)
			if err != nil {
				return err
			}
		}

	case []byte:
		if len(t) != 0 {
			err := json.Unmarshal(t, this)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("无法将[%v] 反序列化为TemplateParams类型", reflect.TypeOf(v).Name())
	}

	return nil
}
