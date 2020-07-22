package routers

import (
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var logger = loggers.Log

func NewRouter() *gin.Engine {
	r := gin.New()

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Use(cors.Default())
	//r.Use(Auth())

	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	r.Use(loggers.GinLogger(logger))
	r.Use(gin.Recovery())

	AddGroupCode(r)
	AddGroupModel(r)
	AddGroupTraining(r)
	AddGroupDataset(r)
	AddGroupAnnotation(r)
	AddGroupInference(r)
	AddGroupFile(r)
	AddGroupGeneral(r)
	AddGroupUpdatePlatform(r)

	return r
}
