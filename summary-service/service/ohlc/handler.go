package ohlc

import (
	pb "github.com/fadellh/stock-ohlc/summary-service/proto"
	"github.com/go-redis/redis/v8"
)

type Handler struct {
	pb.UnimplementedOhlcServer
	Message     string
	redisClient *redis.Client
}

func New(redis *redis.Client) *Handler {
	return &Handler{redisClient: redis}
}
