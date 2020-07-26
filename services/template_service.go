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

// 将用户路径转换为host绝对路径
func ConvertPath(userName, path string) (string, error) {

	pathPrefix := fmt.Sprintf("/home/%s", userName)
	if !strings.HasPrefix(path, pathPrefix) {
		return "", errors.New("非法输出路径")
	}

	// todo： 从接口读取实际的存储路径
	newPathPrefix := fmt.Sprintf("/dlwsdata/work/%s/", userName)
	newPath := fmt.Sprintf("%s/%s", newPathPrefix, strings.TrimLeft(path, pathPrefix))

	if fileInfo, err := os.Stat(newPath); err != nil {
		return "", err
	} else if !fileInfo.IsDir() {
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
