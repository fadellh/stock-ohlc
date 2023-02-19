package ohlcUsecase

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcEntity "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/entity"
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

func NewOhlcUsecase(mgr manager.Manager) OhlcUsecase {
	opt := new(Options)
	opt.config = mgr.GetConfig()
	opt.OhlcTopic = opt.config.OhlcTopic

	return opt
}

func (o *Options) CalculateOHLC(msg *sarama.ConsumerMessage) {
	log.Info().Msgf("New Message from kafka, message: %v", string(msg.Key))
	var req ohlcEntity.OhlcMessage
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Error().Err(err).Msgf("[Usecase-1] %v", err)
		return
	}

	ohlcCalculation(req)

	return
}

func ohlcCalculation(req ohlcEntity.OhlcMessage) ohlcEntity.OhlcStock {
	result := ohlcEntity.OhlcStock{}
	result.StockCode = req.StockCode
	if req.Quantity == 0 && req.Type == ohlcEntity.A {
		result.OpenPrice = req.Price
	}
	return result
}
