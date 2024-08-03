package config

import (
	"event-bus/router/logger"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type EventSource struct {
	KafkaBroker    string `mapstructure:"server_addr"`
	KafkaTopic     string `mapstructure:"kafka_topic"`
	KafkaAuth      bool   `mapstructure:"kafka_auth"`
	KafkaMechanism string `mapstructure:"kafka_mechanism"`
	KafkaUser      string `mapstructure:"kafka_user"`
	KafkaPassword  string `mapstructure:"kafka_password"`
}

type AppConfig struct {
	Env                  string        `mapstructure:"env"`
	EventSource          EventSource   `mapstructure:"event_source"`
	EventTarget          EventSource   `mapstructure:"event_target"`
	FlushTime            time.Duration `mapstructure:"flush_time"`
	FlushMessageSize     int           `mapstructure:"flush_message_size"`
	BufferMessageChannel int           `mapstructure:"buffer_message_channel"`
}

var (
	instance *AppConfig = nil
	once     sync.Once
)

func getEnv(key string, fallback interface{}) interface{} {
	var rValue interface{}
	value, exists := os.LookupEnv(key)
	if !exists {
		rValue = fallback
	} else {
		rValue = value
	}
	switch fallback.(type) {
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			rValue = intValue
		} else {
			rValue = fallback
		}
	case time.Duration:
		if durationValue, err := time.ParseDuration(value); err == nil {
			rValue = durationValue
		} else {
			rValue = fallback
		}

	case bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			rValue = boolValue
		} else {
			rValue = fallback
		}
	}
	return rValue
}

func InitConfig() *AppConfig {
	once.Do(
		func() {
			err := godotenv.Load(".env")
			if err != nil {
				logger.Infof("Error: ", err)
			}
			eventSource := &EventSource{
				KafkaBroker:    getEnv("EVENT_SOURCE.KAFKA_BROKER", "").(string),
				KafkaTopic:     getEnv("EVENT_SOURCE.KAFKA_TOPIC", "").(string),
				KafkaAuth:      getEnv("EVENT_SOURCE.KAFKA_AUTH", false).(bool),
				KafkaMechanism: getEnv("EVENT_SOURCE.KAFKA_MECHANISM", "").(string),
				KafkaUser:      getEnv("EVENT_SOURCE.KAFKA_USER", "").(string),
				KafkaPassword:  getEnv("EVENT_SOURCE.KAFKA_PASSWORD", "").(string),
			}
			eventTarget := &EventSource{
				KafkaBroker:    getEnv("EVENT_TARGET.KAFKA_BROKER", "").(string),
				KafkaTopic:     getEnv("EVENT_TARGET.KAFKA_TOPIC", "").(string),
				KafkaAuth:      getEnv("EVENT_TARGET.KAFKA_AUTH", false).(bool),
				KafkaMechanism: getEnv("EVENT_TARGET.KAFKA_MECHANISM", "").(string),
				KafkaUser:      getEnv("EVENT_TARGET.KAFKA_USER", "").(string),
				KafkaPassword:  getEnv("EVENT_TARGET.KAFKA_PASSWORD", "").(string),
			}
			instance = &AppConfig{
				Env:                  getEnv("ENV", "development").(string),
				EventSource:          *eventSource,
				EventTarget:          *eventTarget,
				FlushTime:            getEnv("FLUSH_TIME", 500*time.Millisecond).(time.Duration),
				FlushMessageSize:     getEnv("FLUSH_MESSAGE_SIZE", 1000).(int),
				BufferMessageChannel: getEnv("BUFFER_MESSAGE_CHANNEL", 512).(int),
			}
		},
	)

	return instance
}
