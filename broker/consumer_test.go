package broker

import (
	"bytes"
	"context"

	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"indasto1.com/unit3-mhs/configuration"
	"indasto1.com/unit3-mhs/models"
)

var memberId = "123id"

type mockConsumerGroupSession struct {
	ctx context.Context
}

func (mockConsumerGroupSession) Claims() map[string][]int32 { return nil }
func (mockConsumerGroupSession) MemberID() string           { return memberId }
func (mockConsumerGroupSession) GenerationID() int32        { return 0 }
func (mockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
}
func (mockConsumerGroupSession) Commit() {}
func (mockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
}
func (mockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {}
func (mSess mockConsumerGroupSession) Context() context.Context {
	return mSess.ctx
}

type mockConsumerGroupClaim struct {
	messages chan *sarama.ConsumerMessage
}

func (mockConsumerGroupClaim) Topic() string              { return "" }
func (mockConsumerGroupClaim) Partition() int32           { return 0 }
func (mockConsumerGroupClaim) InitialOffset() int64       { return 0 }
func (mockConsumerGroupClaim) HighWaterMarkOffset() int64 { return 0 }
func (mClaim mockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	return mClaim.messages
}

func TestSumsConsumeClaimLog(t *testing.T) {
	out := new(bytes.Buffer)
	logrus.SetOutput(out)

	sModel := models.Sum{
		Num1: 2,
		Num2: 5,
	}
	v, err := json.Marshal(sModel)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	msg := sarama.ConsumerMessage{
		Topic: configuration.SUM_TOPIC,
		Value: sarama.ByteEncoder(v),
	}

	mGroupClain := mockConsumerGroupClaim{
		messages: make(chan *sarama.ConsumerMessage),
	}

	mGroupSession := mockConsumerGroupSession{
		ctx: context.Background(),
	}

	pauseBeforeSend := time.Second * 2
	go func() {
		time.Sleep(pauseBeforeSend)
		select {
		case mGroupClain.messages <- &msg:
		default:
		}
	}()

	sumHandler := ConsumerSumGroupHandler{}
	go func() {
		err = sumHandler.ConsumeClaim(mGroupSession, mGroupClain)
		if err != nil {
			logrus.Errorf("failed to consume msg: %v", err)
		}
	}()

	time.Sleep(pauseBeforeSend + 1*time.Second)
	expectedLog := fmt.Sprintf(`msg="Sum of %v and %v is %v" memberId=%s`, sModel.Num1, sModel.Num2, sModel.Num1+sModel.Num2, memberId)
	result := out.String()
	if !strings.Contains(result, expectedLog) {
		t.Errorf("log doesn't contain expected string\nGot:\n%v\nExpected:\n%v", result, expectedLog)
	}
}
