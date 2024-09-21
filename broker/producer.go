package broker

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"indasto1.com/unit3-mhs/models"
)

var Producer sarama.SyncProducer

func InitProducer(addresses []string) error {
	var err error

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Producer.Return.Successes = true
	cfg.Net.TLS.Enable = false

	Producer, err = sarama.NewSyncProducer(addresses, cfg)
	if err != nil {
		return err
	}

	return nil
}

func CloseProducer() error {
	return Producer.Close()
}

func SendMessage(s models.Sum, topic string) error {
	v, err := json.Marshal(s)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(uuid.New().String()),
		Value: sarama.ByteEncoder(v),
	}

	_, _, err = Producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	return nil
}
