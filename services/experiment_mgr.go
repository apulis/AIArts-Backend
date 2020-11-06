package services

import (
	_ "github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/models"
	"github.com/gin-gonic/gin"
)

type CommQueryParams struct {
	Offset  uint
	PageNum uint     `form:"pageNum" `
	Limit   uint     `form:"pageSize"`
	OrderBy string   `form:"orderBy" `
	Order   string   `form:"order" `
	All     uint     `form:"all"`
}

func (s *CommQueryParams) SetQueryParams(c *gin.Context) {
	 c.ShouldBind(s)
	 if s.Limit > 100 {
	 	s.Limit=100
	 } else if s.Limit < 10 {
	 	s.Limit=10
	}
	 s.Offset= s.PageNum * s.Limit
     if s.Order != "" && s.Order != "ASC" && s.Order != "DESC" {
     	s.Order=""
	 }
	 //@todo: check for order by ???
	 s.OrderBy=""

}

func ListExpProjects(queryParam CommQueryParams, username string) ([]models.ExpProject, uint, error) {
	return models.ListAllExpProjects(queryParam.Offset, queryParam.Limit, queryParam.All,
		queryParam.OrderBy, queryParam.Order)
}

func CreateExpProject(project* models.ExpProject) error {
	return models.CreateExpProject(project)
}

func  QueryExpProject(id uint64,project*models.ExpProject)error{
	return models.QueryExpProject(uint(id),project)
}

func  UpdateExpProject(id uint64,new_name string,project*models.RequestUpdates)error{
	if project == nil{
		return models.RenameExpProject(uint(id),new_name)
	}else{
		if len(new_name) != 0{
			(*project)["name"]=new_name
		}
		return models.UpdateExpProject(uint(id),project)
	}
}

func MarkExpProject(id uint64, hide bool) error {
	if hide {
		return models.HideExpProject(uint(id))
	}else{
		return models.RestoreExpProject(uint(id))
	}
}

func ListExperiments(queryParam CommQueryParams, projectID uint64) ([]models.Experiment, uint, error) {
	return models.ListAllExperiments(uint(projectID), queryParam.Offset, queryParam.Limit, queryParam.All,
		queryParam.OrderBy, queryParam.Order)
}

func CreateExperiment(experiment*models.Experiment)error{
	return models.CreateExperiment(experiment)
}

func  UpdateExperiment(id uint64,new_name string,experiment*models.RequestUpdates)error{
	if experiment == nil{
		return models.RenameExperiment(uint(id),new_name)
	}else{
		if len(new_name) != 0{
			(*experiment)["name"]=new_name
		}
		//@todo: how to store user define data ???
		return models.UpdateExperiment(uint(id),experiment)
	}
}
func MarkExperiment(id uint64, hide bool) error {
	if hide {
		return models.HideExperiment(uint(id))
	}else{
		return models.RestoreExperiment(uint(id))
	}
}

func  QueryExperiment(id uint64,experiment*models.Experiment)error{
	return models.QueryExperiment(uint(id),experiment)
}
