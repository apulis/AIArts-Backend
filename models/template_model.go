package models

import "github.com/jinzhu/gorm"

const (
	PublicTemplate  = 1
	PrivateTemplate = 2
)

// 以下结构体工程通用
type Template struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Engine      string            `json:"engine"`
	DeviceType  string            `json:"deviceType"`
	DeviceNum   int               `json:"deviceNum"`
	CodePath    string            `json:"codePath"`
	StartupFile string            `json:"startupFile"`
	OutputPath  string            `json:"outputPath"`
	DatasetPath string            `json:"datasetPath"`
	Params      map[string]string `json:"params"`
	CreateTime  string            `json:"createTime"`
}

// 与数据库一一对应
type TemplateDbItem struct {
	ID int `gorm:"primary_key" json:"id"`

	Name    string `gorm:"not null" json:"name"`
	Scope   int    `gorm:"not null" json:"scope"`
	Data    string `gorm:"not null" json:"data"`
	JobType string `gorm:"not null" json:"jobType"`
	Creator string `json:"creator"`

	CreatedAt UnixTime `json:"createdAt"`
	UpdatedAt UnixTime `json:"updatedAt"`
}

type TemplateProvider struct {
	gormDb *gorm.DB
}

func NewTemplateProvider(gormDb *gorm.DB) *TemplateProvider {
	return &TemplateProvider{gormDb: gormDb}
}

// 分页查询
func (this *TemplateProvider) FindPage(order string, offset, limit int, query string, args ...interface{}) ([]*TemplateDbItem, error) {

	var tmp []*TemplateDbItem
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
