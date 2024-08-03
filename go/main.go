package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"event-bus/router/config"
	"event-bus/router/kafka"

	"github.com/IBM/sarama"
)


func main() {
	conf := config.InitConfig()
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := kafka.Consumer{
		Ready:          make(chan bool),
		MessageChannel: make(chan *sarama.ConsumerMessage, conf.BufferMessageChannel), // Buffered channel
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(conf.EventSource.KafkaBroker, ","), fmt.Sprintf("%s-group-%s", conf.EventSource.KafkaTopic, conf.Env), config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{conf.EventSource.KafkaTopic}, &consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
				// Sleep a bit before trying to reconnect
				select {
				case <-time.After(time.Second):
				case <-ctx.Done():
					return
				}
			}
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		kafka.RunAsyncProducer(conf, consumer.MessageChannel)
	}()

	<-sigterm
	log.Println("Termination signal received. Shutting down...")
	cancel()
	close(consumer.MessageChannel)
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Fatalf("Error closing client: %v", err)
	}
	log.Println("Consumer group closed")
}