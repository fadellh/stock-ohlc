package publiser

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

// SendMessage function to send message into kafka
func (p *KafkaProducer) SendMessage(topic, msg string) error {

	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
	}

	partition, offset, err := p.Producer.SendMessage(kafkaMsg)
	if err != nil {
		log.Error().Err(err).Msgf("Send message error: %v", err)
		return err
	}

	log.Info().Msgf("Send message success, Topic %v, Partition %v, Offset %d", topic, partition, offset)
	return nil
}
