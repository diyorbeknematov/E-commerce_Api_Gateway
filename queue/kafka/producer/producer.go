package producer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	ProducerMessage(topic string, msg []byte) error
	Close()
}

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) KafkaProducer {
	return &kafkaProducer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *kafkaProducer) ProducerMessage(topic string, msg []byte) error {
	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Value: msg,
	})
}

func (p *kafkaProducer) Close() {
	p.writer.Close()
}
