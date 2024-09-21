// Consumer service, which listents on one topic (sum) in kafka
// cli args:
//   -kafka-addrs 	- (string) kafka's host:port list separated by coma
//	 -num-consumers - (int) number of partitions for sum topic in kafka
package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/broker"
	"indasto1.com/unit3-mhs/configuration"
)

var (
	kafkaAddrs     = flag.String("kafka-addrs", "", "kafka's host:port list separated by coma")
	numOfConsumers = flag.Int("num-consumers", 1, "number of partitions for sum topic")
)

func main() {
	flag.Parse()

	if *kafkaAddrs == "" {
		log.Fatal("Kafka's addresses are not passed")
	}

	if *numOfConsumers < 1 {
		log.Fatal("Number of partitions can't b less then 1")
	}

	kafkaAddrsArr := strings.Split(*kafkaAddrs, ",")
	consumers := make([]sarama.ConsumerGroup, 0, *numOfConsumers)
	for i := 0; i < *numOfConsumers; i++ {
		consumer, err := broker.CreateConsumerGroup(kafkaAddrsArr)
		if err != nil {
			log.WithError(err).Fatal("Failed to create consumer")
		}
		defer consumer.Close()

		consumers = append(consumers, consumer)
	}

	wg := sync.WaitGroup{}
	cancelFn := broker.StartSumCalculator(consumers, configuration.SUM_TOPIC, &wg)

	log.WithField("Consumers num", *numOfConsumers).Info("Service started")
	sigtermCh := make(chan os.Signal, 1)
	signal.Notify(sigtermCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigtermCh

	log.Info("Shutdown consumer service...")

	cancelFn()

	wg.Wait()
}
