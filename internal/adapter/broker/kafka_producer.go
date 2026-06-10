package broker

import (
	"github.com/segmentio/kafka-go"
	"context"
	"encoding/json"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

type Event struct {
    EventType string      `json:"event_type"`
    Payload   interface{} `json:"payload"`
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaProducer) Publish(ctx context.Context, eventType string, payload interface{}) error {
	data, err := json.Marshal(Event{
		EventType: eventType,
		Payload: payload,
	})
	if err != nil {
		return err
	}
	
	message := kafka.Message{
		Key:  []byte(eventType),
		Value: data,
	}
	return p.writer.WriteMessages(ctx, message)
}