package main

import (
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/routers"

	_ "github.com/apulis/AIArtsBackend/docs"
	_ "github.com/apulis/AIArtsBackend/loggers"
)

var logger = loggers.Log

func main() {
	port := configs.Config.Port
	router := routers.NewRouter()

	logger.Info("AIArtsBackend started, listening and serving HTTP on: ", 9000)
	router.Run(fmt.Sprintf(":%d", port))
}
