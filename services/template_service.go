package services

import (
	"errors"
	"fmt"
	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/models"
	"log"
	"os"
	"strings"
)

func GetAllTemplate(userName string, page, size, scope int, jobType, searchWord, orderBy, order string) ([]*models.TemplateItem, int, int, error) {

	query := ""
	provider := models.NewTemplateProvider(database.Db)

	var err error
	var total int
	var templates []*models.Templates
	var extQuery, orderQuery string

	if len(searchWord) > 0 {
		extQuery += " and name like '%" + searchWord + "%'"
	}

	if len(orderBy) > 0 {
		if strings.ToLower(order) == "asc" {
			orderQuery = orderBy + " asc"
		} else {
			orderQuery = orderBy + " desc"
		}
	} else {
		orderQuery = "created_at desc"
	}

	// 用户私有 + 公有
	if scope == models.TemplatePublicPrivate {
		query = "(scope in (?, ?) or creator = ?) and job_type = ?"
		if templates, err = provider.FindPage("", (page-1)*size, size, query,
			models.TemplatePublic, models.TemplatePredefined, userName, jobType); err != nil {
			return nil, 0, 0, err
		}
	} else if scope == models.TemplatePublic {

		query = "scope in (?, ?) and job_type = ?"
		if len(extQuery) > 0 {
			query += extQuery
		}

		if templates, err = provider.FindPage("", (page-1)*size, size, query, models.TemplatePublic, models.TemplatePredefined, jobType); err != nil {
			return nil, 0, 0, err
		}
	} else if scope == models.TemplatePredefined {

		query = "scope = ? and job_type = ?"
		if len(extQuery) > 0 {
			query += extQuery
		}

		if templates, err = provider.FindPage("", (page-1)*size, size, query, scope, jobType); err != nil {
			return nil, 0, 0, err
		}

	} else if scope == models.TemplatePrivate {

		query = "creator = ? and job_type = ? "
		if len(extQuery) > 0 {
			query += extQuery
		}

		if templates, err = provider.FindPage(orderQuery, (page-1)*size, size, query, userName, jobType); err != nil {
			return nil, 0, 0, err
		}

		if total, err = provider.Count(query, userName, jobType); err != nil {
			return nil, 0, 0, err
		}
	}

	retItems := make([]*models.TemplateItem, 0)
	for _, v := range templates {
		if item := v.ToTemplateItem(); item != nil {
			retItems = append(retItems, item)
		}
	}

	totalPages := total / size
	if (total % size) != 0 {
		totalPages += 1
	}

	return retItems, total, totalPages, nil
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
	//return database.Db.Raw(`DELETE FROM ai_arts.templates where id=?`, id).Error
	db := database.Db.Exec(`DELETE FROM templates where id=?`, id)
	return db.Error
}

func GetTemplate(userName string, id int64) (*models.Templates, error) {

	provider := models.NewTemplateProvider(database.Db)
	item, err := provider.FindById(id)

	return item, err
}

// 将用户路径转换为host绝对路径
func ConvertPath(userName, path string) (string, error) {

	pathPrefix := fmt.Sprintf("/home/%s", userName)
	if !strings.HasPrefix(path, pathPrefix) {
		return "", errors.New("非法输出路径")
	}

	// todo： 从接口读取实际的存储路径
	newPathPrefix := fmt.Sprintf("/dlwsdata/work/%s/", userName)
	newPath := fmt.Sprintf("%s/%s", newPathPrefix, strings.TrimPrefix(path, pathPrefix))

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return "", fmt.Errorf("路径非合法目录：%s", newPath)
	}

	return newPath, nil
}

// 上传结束后的处理工作
func UploadDone(userName, filePath string) error {

	if fileInfo, err := os.Stat(filePath); err != nil {
		return err
	} else if fileInfo.IsDir() {
		return fmt.Errorf("非法文件：%s", filePath)
	}

	// Change permissions Linux.
	if err := os.Chmod(filePath, 0755); err != nil {
		log.Println("更改文件模式报错", filePath, err)
		return err
	}

	return nil
}
