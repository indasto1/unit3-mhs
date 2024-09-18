package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/handlers"
)

func CreateRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), handlers.HandleError)

	router.POST("/data", handlers.ProduceData)

	return router
}

func main() {
	// TODO: parse flags

	server := &http.Server{
		// TODO: replace by value from config
		Addr:    ":8080",
		Handler: CreateRouter(),
	}

	go func() {
		log.Infof("Listening on %s", ":8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal()
		}
	}()

	quitSignalCh := make(chan os.Signal, 1)
	signal.Notify(quitSignalCh, syscall.SIGTERM, syscall.SIGINT)
	<-quitSignalCh

	log.Info("Shutdown router...")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Error("Failed to finalize connections")
	}
}
