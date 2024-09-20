package broker

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"indasto1.com/unit3-mhs/models"
)

var producer sarama.SyncProducer

func InitProducer(addresses []string) error {
	var err error
	producer, err = sarama.NewSyncProducer(addresses, nil)
	if err != nil {
		return err
	}
	defer producer.Close()

	return nil
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

	_, _, err = producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	return nil
}
