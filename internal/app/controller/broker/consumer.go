package broker

import (
	"context"
	"encoding/json"
	"github.com/matster07/user-balance-service/internal/app/configs"
	"github.com/matster07/user-balance-service/internal/app/controller/server"
	"github.com/matster07/user-balance-service/internal/app/data/dto"
	"github.com/matster07/user-balance-service/internal/pkg/logging"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type Consumer struct {
	Reader *kafka.Reader
}

var instance *Consumer
var once sync.Once

func GetConsumer() *Consumer {
	once.Do(func() {
		config := configs.GetConfig()

		kafkaConfig := kafka.ReaderConfig{
			Brokers:         config.Brokers,
			GroupID:         "readers",
			Topic:           config.Topic,
			MinBytes:        10e3,            // 10KB
			MaxBytes:        10e6,            // 10MB
			MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
			ReadLagInterval: -1,
		}

		instance = &Consumer{
			Reader: kafka.NewReader(kafkaConfig),
		}

		logging.GetLogger().Infof("consumer is ready")
	})

	return instance
}

// Подписываемся на топик с обновлениями статусов заказов
func (c *Consumer) Read(handler server.Handler) {
	logger := logging.GetLogger()

	go func() {
		for {
			m, err := c.Reader.ReadMessage(context.Background())
			if err != nil {
				continue
			}

			var deliverStatusDto dto.DeliverStatusDTO
			err = json.Unmarshal(m.Value, &deliverStatusDto)
			if err != nil {
				logger.Error("error while unmarshal message: %s", err.Error())
				continue
			}

			err = handler.Process(deliverStatusDto)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
		}
	}()
}
