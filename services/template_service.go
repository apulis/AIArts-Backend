package services

import (
	"fmt"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/models"
)

func GetAllTemplate(userName string, page, size, scope int, jobType string) ([]*models.TemplateItem, int, int, error) {

	query := ""
	provider := models.NewTemplateProvider(database.Db)

	var err error
	var templates []*models.Templates

	// 用户私有 + 公有
	if scope == models.TemplateUserPublic {

		query = "scope = ? or creator = ?"
		if templates, err = provider.FindPage("", (page-1)*size, size, query, scope, userName); err != nil {
			return nil, 0, 0, err
		}
		// public
	} else if scope == models.TemplatePublic {

		query = "scope = ? and job_type = ?"
		if templates, err = provider.FindPage("", (page-1)*size, size, query, scope, jobType); err != nil {
			return nil, 0, 0, err
		}

	} else {
		query = "creator = ? and job_type = ?"
		if templates, err = provider.FindPage("", (page-1)*size, size, query, userName, jobType); err != nil {
			return nil, 0, 0, err
		}
	}

	retItems := make([]*models.TemplateItem, 0)
	for _, v := range templates {
		if item := v.ToTemplateItem(); item != nil {
			retItems = append(retItems, item)
		}
	}

	return retItems, len(retItems), 1, nil
}

func CreateTemplate(userName string, scope int, jobType string, template models.TemplateParams) (int64, error) {

	provider := models.NewTemplateProvider(database.Db)

	record := &models.Templates{}
	record.Load(scope, userName, jobType, template)

	id, err := provider.Insert(record.ToMap())
	return id, err
}

func UpdateTemplate(id int64, userName string, scope int, jobType string, template models.TemplateParams) error {

	provider := models.NewTemplateProvider(database.Db)

	record := &models.Templates{}
	record.Load(scope, userName, jobType, template)

	var err error
	if id, err = provider.Update(id, record.ToMap()); id == 0 {
		fmt.Printf("update template err: %v", err)
	}

	return err
}

func DeleteTemplate(userName string, id int64) error {
	return database.Db.Raw(`DELETE FROM ai_arts.templates where id=?`, id).Error
}

func GetTemplate(userName string, id int64) (*models.Templates, error) {

	provider := models.NewTemplateProvider(database.Db)
	item, err := provider.FindById(id)

	return item, err
}
