package routers

import (
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Resp struct {
	Code int
	Msg  string
	Data gin.H
}

var logger = loggers.Log

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(loggers.GinLogger(logger))
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))
	AddGroupCode(r)
	AddGroupModel(r)
	AddGroupTraining(r)
	AddGroupDataset(r)
	AddGroupAnnotation(r)
	AddGroupInference(r)

	return r
}
