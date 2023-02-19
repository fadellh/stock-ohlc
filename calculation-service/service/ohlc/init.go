package ohlc

import (
	"os"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcUsecase "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/usecase"
	"github.com/rs/zerolog/log"
)

func NewOHLC(mgr manager.Manager, signals chan os.Signal) error {
	ohlc := ohlcUsecase.NewOhlcUsecase(mgr)

	cfg := mgr.GetConfig()
	consumer := mgr.GetKafka().Consumer()

	chanMessage := make(chan *sarama.ConsumerMessage, 256)

	partitionList, err := consumer.Partitions(cfg.OhlcTopic)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get partition got error %v", err)
		return err
	}

	for _, partition := range partitionList {
		go consumeMessage(consumer, cfg.OhlcTopic, partition, chanMessage)
	}

	log.Info().Msgf("Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			ohlc.CalculateOHLC(msg)
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}

	return nil
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			log.Error().Err(err).Msgf("Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}
}
