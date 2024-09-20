package broker

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/models"
)

func CreateConsumerGroup(addresses []string) (sarama.ConsumerGroup, error) {
	return sarama.NewConsumerGroup(addresses, "sumConsumer", nil)
}

type ConsumerSumGroupHandler struct{}

func (c ConsumerSumGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c ConsumerSumGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c ConsumerSumGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			var sumModel models.Sum
			err := json.Unmarshal(msg.Value, &sumModel)
			if err != nil {
				log.WithError(err).Error("Failed to unmarshal on sum model claim")
				continue
			}

			log.Infof("Sum of %v and %v is %v\n", sumModel.Num1, sumModel.Num2, sumModel.Num1+sumModel.Num2)
		}
	}
}

func StartSumCalculator(consumerGroup sarama.ConsumerGroup, topic string, wg *sync.WaitGroup) context.CancelFunc {
	ctx, ctxCancelFn := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()

		groupHandler := ConsumerSumGroupHandler{}
		for {
			if ctx.Err() != nil {
				return
			}

			if err := consumerGroup.Consume(ctx, []string{topic}, groupHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}

				log.WithError(err).Fatal("Error from consumer")
			}
		}
	}()

	return ctxCancelFn
}
