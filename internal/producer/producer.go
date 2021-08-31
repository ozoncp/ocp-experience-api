package producer

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Send(msg ...EventMsg)
}

// NewProducer creates new kafka producer
func NewProducer(topic string, kafkaProducer sarama.SyncProducer) *producer {
	return &producer{topic: topic, kafkaProducer: kafkaProducer}
}

type producer struct {
	topic         string
	kafkaProducer sarama.SyncProducer
}

// Send sends a batch of message to Kafka broker
func (p *producer) Send(eventMessages ...EventMsg) {
	if len(eventMessages) == 0 {
		return
	}

	producerMessages := make([]*sarama.ProducerMessage, 0, len(eventMessages))

	for _, m := range eventMessages {
		producerMessages = append(producerMessages,
			&sarama.ProducerMessage{
				Topic:     p.topic,
				Partition: -1,
				Value:     m,
			},
		)
	}

	err := p.kafkaProducer.SendMessages(producerMessages)

	if err != nil {
		log.Error().Msgf("failed to send messages to Kafka: %v", err)
	}
}
