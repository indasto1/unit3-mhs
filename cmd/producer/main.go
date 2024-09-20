package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/handlers"
)

var (
	listen     = *flag.String("listen", ":8080", "router's port on which to listen")
	kafkaAddrs = *flag.String("kafka-addrs", "", "kafka's host:port list separated by coma")
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

	sigtermCh := make(chan os.Signal, 1)
	signal.Notify(sigtermCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigtermCh

	log.Info("Shutdown router...")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Error("Failed to finalize connections")
	}
}
