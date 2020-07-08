package routers

import (
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/gin-contrib/cors"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var logger = loggers.Log

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(cors.Default())

	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	r.Use(loggers.GinLogger(logger))
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))
	AddGroupCode(r)
	AddGroupModel(r)
	AddGroupTraining(r)
	AddGroupDataset(r)
	AddGroupAnnotation(r)
	AddGroupInference(r)
	AddGroupFile(r)

	return r
}
