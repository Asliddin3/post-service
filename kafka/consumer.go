package kafka

import (
	"context"
	"fmt"

	"github.com/Asliddin3/post-servise/config"
	handler "github.com/Asliddin3/post-servise/kafka/handler"
	"github.com/Asliddin3/post-servise/pkg/logger"
	"github.com/Asliddin3/post-servise/pkg/messagebroker"
	"github.com/Asliddin3/post-servise/storage"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	KafkaConsumer *kafka.Reader
	kafkaHandler  *handler.KafkaHandler
	log           logger.Logger
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, log logger.Logger, topic string) messagebroker.Consumer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)
	return &KafkaConsumer{
		KafkaConsumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{connString},
			Topic:    topic,
			MinBytes: 10e3, //10KB
			MaxBytes: 10e6, //10MB

		}),
		kafkaHandler: handler.NewKafkaHandlerFunc(*conf, storage.NewStoragePg(db), log),
		log:          log,
	}
}

func (k KafkaConsumer) Start() {
	fmt.Println(">>> Consumer started.")
	for {
		m, err := k.KafkaConsumer.ReadMessage(context.Background())
		fmt.Println("master of copy paste here", err, m)
		if err != nil {
			k.log.Error("Error on consuming a message:", logger.Error(err))
			break
		}
		err = k.kafkaHandler.Handle(m.Value)
		if err != nil {
			k.log.Error("failed to handle consumed topic:",
				logger.String("on topic", m.Topic), logger.Error(err))
		} else {
			fmt.Println()
			k.log.Info("Successfully consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message", "success"))
			fmt.Println()
		}
	}
	err := k.KafkaConsumer.Close()
	if err != nil {
		k.log.Error("Error on closing consumer:", logger.Error(err))
	}
}
