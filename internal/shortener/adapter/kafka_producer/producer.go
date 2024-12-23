package kafka_producer

import "github.com/segmentio/kafka-go"

type Producer struct {
	writer *kafka.Writer
}

func New(writer *kafka.Writer) *Producer {
	return &Producer{writer: writer}
}
