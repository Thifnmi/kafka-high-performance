package main

import (
    "log"
    "github.com/IBM/sarama"
)

func main() {
    config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
    producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
    if err != nil {
        log.Fatalf("Failed to start producer: %v", err)
    }
    defer producer.Close()
    for true {

        message := &sarama.ProducerMessage{
            Topic: "IAM",
            Value: sarama.StringEncoder("Hello Kafka"),
        }
        _, _, err = producer.SendMessage(message)
        if err != nil {
            log.Fatalf("Failed to send message: %v", err)
        }
        log.Println("Message sent successfully")
    }
}
