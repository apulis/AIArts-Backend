package main

import (
	"fmt"

	"github.com/apulis/AIArtsBackend/config"
	"github.com/apulis/AIArtsBackend/router"

	_ "github.com/apulis/AIArtsBackend/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	port := config.Config.Port
	r := router.NewRouter()

	r.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER"))

	r.Run(fmt.Sprintf(":%d", port))
}
