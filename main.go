package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/loggers"
	"github.com/apulis/AIArtsBackend/routers"
)

var logger = loggers.Log

func main() {

	port := configs.Config.Port
	router := routers.NewRouter()

	logger.Info("AIArtsBackend started, listening and serving HTTP on: ", port)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go taskCleanTmpFiles(configs.Config.File.DatasetDir, configs.Config.File.CleanEverySeconds)
	go taskCleanTmpFiles(configs.Config.File.ModelDir, configs.Config.File.CleanEverySeconds)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func taskCleanTmpFiles(dir string, seconds int64) {
	for {
		cleanTmpFiles(dir)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func cleanTmpFiles(dir string) {
	tmpDir := dir + "/tmp"
	logger.Info("Checking dir: ", tmpDir)

	_, err := os.Stat(tmpDir)
	if err != nil {
		err = os.MkdirAll(tmpDir, os.ModeDir|os.ModePerm)
		if err != nil {
			logger.Info("Creating dir error: ", err.Error())
		}
	}

	reader, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		logger.Error("Read dir error: ", err.Error())
	}
	for _, fi := range reader {
		logger.Info("Checking file: ", tmpDir, fi.Name())
		duration := time.Now().Sub(fi.ModTime())
		if duration.Seconds() > float64(configs.Config.File.CleanBeforeSeconds) {
			logger.Info("Removing file: ", fi.Name())
			os.RemoveAll(tmpDir + fi.Name())
		}
	}
}
