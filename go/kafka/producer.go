package kafka

// SIGUSR1 toggle the pause/resume consumption
import (
	"event-bus/router/config"
	"event-bus/router/logger"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	// "encoding/json"
	"github.com/IBM/sarama"
)


func RunAsyncProducer(conf *config.AppConfig, messages <-chan *sarama.ConsumerMessage) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = false
	config.Producer.Flush.Frequency = conf.FlushTime
	config.Producer.Flush.Messages = conf.FlushMessageSize

	producer, err := sarama.NewAsyncProducer(strings.Split(conf.EventTarget.KafkaBroker, ","), config)
	if err != nil {
		logger.Fatalf("Failed to start Sarama producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Fatalf("Failed to close Sarama producer: %v", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			logger.Debugf("Failed to produce message: %v", err)
		}
	}()

	for message := range messages {
		producerMessage := &sarama.ProducerMessage{
			Topic: conf.EventTarget.KafkaTopic,
			Value: sarama.StringEncoder(message.Value),
		}
		select {
		case producer.Input() <- producerMessage:
		case <-time.After(time.Second):
			logger.Info("Timeout on sending message to producer")
		}
	}

	wg.Wait()
}
