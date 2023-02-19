package ohlcUsecase

import (
	"encoding/json"
	"math"

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
	var req ohlcEntity.OhlcMessage
	if err := json.Unmarshal(msg.Value, &req); err != nil {
		log.Error().Err(err).Msgf("[CalculateOHLC-Usecase-1] %v", err)
		return
	}
	log.Info().Msgf("New Message from kafka, stock code: %v, price %d", req.StockCode, req.Price)

	if req.Price == 0 && req.ExecutionPrice > 0 {
		req.Price = req.ExecutionPrice
	}

	if req.Quantity == 0 && req.ExecutedQuantity > 0 {
		req.Quantity = req.ExecutedQuantity
	}

	result := ohlcEntity.OhlcStock{}
	result.StockCode = req.StockCode
	if req.Quantity == 0 {
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

	result = ohlcTrxCalculation(req, *ohlcData)
	if result.StockCode == "" {
		return
	}

	err = o.repo.StoreRedis(result)
	if err != nil {
		log.Error().Err(err).Msgf("[CalculateOHLC-Usecase-4] %v", err)
		return
	}

	return
}

func ohlcTrxCalculation(req ohlcEntity.OhlcMessage, currentOhlc ohlcEntity.OhlcStock) ohlcEntity.OhlcStock {
	if req.StockCode != currentOhlc.StockCode {
		log.Error().Msgf("Stock code different")
		return ohlcEntity.OhlcStock{}
	}

	if currentOhlc.Volume == 0 {
		currentOhlc.OpenPrice = req.Price
		currentOhlc.LowestPrice = req.Price
	}

	if req.Price >= currentOhlc.HighestPrice {
		currentOhlc.HighestPrice = req.Price
	}

	if req.Price <= currentOhlc.LowestPrice {
		currentOhlc.LowestPrice = req.Price
	}

	currentOhlc.ClosePrice = req.Price
	currentOhlc.Volume += req.Quantity
	currentOhlc.Value += (req.Quantity * req.Price)

	volume := float64(currentOhlc.Volume)
	value := float64(currentOhlc.Value)
	averagePrice := value / volume

	currentOhlc.AveragePrice = int(math.Round(averagePrice))
	return currentOhlc
}
