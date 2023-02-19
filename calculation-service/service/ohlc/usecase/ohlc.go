package ohlcUsecase

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcEntity "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/entity"
	ohlcRepository "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/repository"
	"github.com/rs/zerolog/log"
)

type OhlcUsecase interface {
	CalculateOHLC(*sarama.ConsumerMessage)
}

type Options struct {
	config    *config.Config
	consumer  sarama.Consumer
	OhlcTopic string
	repo      ohlcRepository.OhlcRepository
}

func NewOhlcUsecase(mgr manager.Manager) OhlcUsecase {
	opt := new(Options)
	opt.config = mgr.GetConfig()
	opt.OhlcTopic = opt.config.OhlcTopic
	opt.repo = ohlcRepository.NewOhlcRepository(mgr)

	return opt
}

func (o *Options) CalculateOHLC(msg *sarama.ConsumerMessage) {
	log.Info().Msgf("New Message from kafka, message: %v", string(msg.Key))
	var req ohlcEntity.OhlcMessage
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Error().Err(err).Msgf("[CalculateOHLC-Usecase-1] %v", err)
		return
	}

	result := ohlcEntity.OhlcStock{}
	result.StockCode = req.StockCode
	if req.Quantity == 0 && req.Type == ohlcEntity.A {
		result.PreviousPrice = req.Price
	}

	ohlcData, err := o.repo.GetRedis(result)
	if err != nil {
		log.Error().Err(err).Msgf("[CalculateOHLC-Usecase-2] %v", err)
		return
	}

	if ohlcData == nil {
		err := o.repo.StoreRedis(result)
		if err != nil {
			log.Error().Err(err).Msgf("[CalculateOHLC-Usecase-3] %v", err)
			return
		}
		return
	}

	result = ohlcTrxCalculation(*ohlcData)
	return
}

func ohlcTrxCalculation(data ohlcEntity.OhlcStock) ohlcEntity.OhlcStock {
	return data
}
