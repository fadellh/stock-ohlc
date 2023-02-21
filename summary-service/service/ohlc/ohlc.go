package ohlc

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	pb "github.com/fadellh/stock-ohlc/summary-service/proto"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderType string

const (
	A OrderType = "A"
	P OrderType = "P"
	E OrderType = "E"
)

type OhlcStock struct {
	StockCode     string    `json:"stock_code"`
	PreviousPrice int       `json:"prev"`
	OpenPrice     int       `json:"open"`
	HighestPrice  int       `json:"highest"`
	LowestPrice   int       `json:"lowest"`
	ClosePrice    int       `json:"close"`
	Volume        int64     `json:"volume"`
	Value         string    `json:"value"`
	AveragePrice  int       `json:"average"`
	Type          OrderType `json:"type"`
}

func (h *Handler) GetOhlcSummary(ctx context.Context, in *pb.SummaryRequest) (*pb.SummaryResponse, error) {
	value, err := h.redisClient.Get(ctx, in.Stockcode).Result()
	if err == redis.Nil {
		log.Error().Err(err).Msgf("[GetOhlcSummary-1] %v", err)
		err = errors.New("no data entry")
		return nil, status.Errorf(
			codes.NotFound,
			err.Error(),
		)
	}

	if err != nil {
		log.Error().Err(err).Msgf("[GetOhlcSummary-2] %v", err)
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	var data *OhlcStock
	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		log.Error().Err(err).Msgf("[GetOhlcSummary-3] %v", err)
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	if data.Value == "" {
		data.Value = "0"
	}

	currentValue, err := strconv.ParseInt(data.Value, 10, 64)
	if err != nil {
		log.Error().Msgf("int err: %v", err.Error())
		return nil, status.Errorf(
			codes.Internal,
			err.Error(),
		)
	}

	log.Info().Msgf("Stock code: %v | Close Price: %d | value %s", data.StockCode, data.ClosePrice, data.Value)
	return &pb.SummaryResponse{
		Code:    data.StockCode,
		Prev:    int32(data.PreviousPrice),
		Open:    int32(data.OpenPrice),
		Highest: int32(data.HighestPrice),
		Lowest:  int32(data.LowestPrice),
		Close:   int32(data.ClosePrice),
		Average: int32(data.AveragePrice),
		Volume:  data.Volume,
		Value:   currentValue,
	}, nil
}
