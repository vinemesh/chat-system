// Code for interacting with the kafka cluster
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	Producer *kafka.Producer
}

func NewProducer(brokers string) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "all",
	})
	if err != nil {
		return nil, err
	}
	return &Producer{Producer: p}, nil
}

// ProduceMessage produces a message for the given topic with the specified key and value.
//
// Parameters:
//
//	topic string - the topic to produce the message to
//	key string - the key of the message
//	value string - the value of the message
//
// Return type: error
func (p *Producer) ProduceMessage(topic, key, value string) error {
	return p.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}, nil)
}
