package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/broker"
	"indasto1.com/unit3-mhs/configuration"
)

var (
	kafkaAddrs = *flag.String("kafka-addrs", "", "kafka's host:port list separated by coma")
)

func main() {
	flag.Parse()

	if kafkaAddrs == "" {
		log.Fatal("Kafka's addresses are not passed")
	}

	consumerGroup, err := broker.CreateConsumerGroup(strings.Split(kafkaAddrs, ","))
	if err != nil {
		log.WithError(err).Fatal("Failed to create consumer")
	}
	defer consumerGroup.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	cancelFn := broker.StartSumCalculator(consumerGroup, configuration.SUM_TOPIC, &wg)

	sigtermCh := make(chan os.Signal, 1)
	signal.Notify(sigtermCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigtermCh

	log.Info("Shutdown consumer service...")

	cancelFn()

	wg.Wait()
}
