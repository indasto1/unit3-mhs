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
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_6_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	return sarama.NewConsumerGroup(addresses, "sumConsumer", cfg)
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

			log.WithField("memberId", session.MemberID()).Infof("Sum of %v and %v is %v", sumModel.Num1, sumModel.Num2, sumModel.Num1+sumModel.Num2)
		}
	}
}

func StartSumCalculator(consumerGroups []sarama.ConsumerGroup, topic string, wg *sync.WaitGroup) context.CancelFunc {
	ctx, ctxCancelFn := context.WithCancel(context.Background())

	for _, consumerGroup := range consumerGroups {
		wg.Add(1)
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
	}

	return ctxCancelFn
}
