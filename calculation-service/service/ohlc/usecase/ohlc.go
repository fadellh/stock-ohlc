package ohlcUsecase

import (
	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	"github.com/rs/zerolog/log"
)

type OhlcUsecase interface {
	CalculateOHLC(*sarama.ConsumerMessage)
}

type Options struct {
	config    *config.Config
	consumer  sarama.Consumer
	OhlcTopic string
}

func NewOhlcUsecase(mgr manager.Manager) (OhlcUsecase, error) {
	opt := new(Options)
	opt.config = mgr.GetConfig()
	opt.consumer = mgr.GetKafka().Consumer()
	opt.OhlcTopic = opt.config.OhlcTopic

	return opt, nil
}

func (o *Options) CalculateOHLC(msg *sarama.ConsumerMessage) {
	log.Info().Msgf("New Message from kafka, message: %v", string(msg.Value))
	return
}
