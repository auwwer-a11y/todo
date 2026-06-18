package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type KafkaConsumer struct {
	logger *slog.Logger
	reader *kafka.Reader
}

func NewKafkaConsumer(logger *slog.Logger, brokers []string, topic string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		logger: logger,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic: topic,
			GroupID: groupID,
		}),
	}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}

