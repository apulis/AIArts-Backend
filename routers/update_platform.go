package routers

import (
	"github.com/apulis/AIArtsBackend/models"
	"github.com/apulis/AIArtsBackend/services"
	"github.com/gin-gonic/gin"
)

func AddGroupUpdatePlatform(r *gin.Engine){
	group := r.Group("/ai_arts/api/updatePlatform")

	group.Use(Auth())

	group.Get("/latest", wrapper(updateToLatest))
}

// @Summary list datasets
// @Produce  json
// @Param pageNum query int true "page number, from 1"
// @Param pageSize query int true "count per page"
// @Success 200 {object} APISuccessRespGetDatasets "success"
// @Failure 400 {object} APIException "error"
// @Failure 404 {object} APIException "not found"
// @Router /ai_arts/api/datasets [get]
func updateToLatest(c *gin.Context) error{
	return 111
}