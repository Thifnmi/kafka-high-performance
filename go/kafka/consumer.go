package kafka

import (
	// "math/rand"
	"event-bus/router/logger"

	"crypto/sha256"
	"crypto/sha512"
	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}



type RoutingKey struct {
	Key    string   `json:"key"`
	Topics []string `json:"topics"`
}

type Consumer struct {
	Ready chan bool
	MessageChannel chan *sarama.ConsumerMessage
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger.Info("ConsumeClaim started")
	for message := range claim.Messages() {
		logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		consumer.MessageChannel <- message
		sess.MarkMessage(message, "")
	}
	logger.Info("ConsumeClaim finished")
	return nil
}
