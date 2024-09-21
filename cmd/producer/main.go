// Producer service, that has has only one POST /data endpoint which when called will push message to kafka.
// Always exposed at 8080 port
// cli args:
//
//	-kafka-addrs 	- (string) kafka's host:port list separated by coma
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/broker"
	"indasto1.com/unit3-mhs/handlers"
)

const (
	listen = ":8080"
)

var (
	kafkaAddrs = flag.String("kafka-addrs", "", "kafka's host:port list separated by coma")
)

func CreateRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), handlers.HandleError)

	router.POST("/data", handlers.ProduceData)

	return router
}

func main() {
	flag.Parse()

	err := broker.InitProducer(strings.Split(*kafkaAddrs, ","))
	if err != nil {
		log.WithError(err).WithField("addrs", strings.Split(*kafkaAddrs, ",")).Fatal("Failed to create producer")
	}
	defer broker.CloseProducer()

	server := &http.Server{
		Addr:    listen,
		Handler: CreateRouter(),
	}

	go func() {
		log.Infof("Listening on %s", listen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Router stopped")
		}
	}()

	log.Info("Service started")
	sigtermCh := make(chan os.Signal, 1)
	signal.Notify(sigtermCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigtermCh

	log.Info("Shutdown router...")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer ctxCancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error("Failed to finalize connections")
	}
}
