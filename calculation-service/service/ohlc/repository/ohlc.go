package ohlcRepository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcEntity "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/entity"
	"github.com/go-redis/redis/v8"
)

type OhlcRepository interface {
	StoreRedis(req ohlcEntity.OhlcStock) error
	GetRedis(req ohlcEntity.OhlcStock) (*ohlcEntity.OhlcStock, error)
}

type repository struct {
	redis *redis.Client
	cfg   *config.Config
}

func NewOhlcRepository(mgr manager.Manager) OhlcRepository {
	repo := new(repository)
	repo.redis = mgr.GetRedis()
	repo.cfg = mgr.GetConfig()

	return repo
}

func (repo *repository) StoreRedis(req ohlcEntity.OhlcStock) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx := context.TODO()
	expiredTime := time.Duration(repo.cfg.RedisExpired) * time.Minute
	err = repo.redis.Set(ctx, req.StockCode, data, expiredTime).Err()
	if err != nil {
		return err
	}

	return nil
}

func (repo *repository) GetRedis(req ohlcEntity.OhlcStock) (*ohlcEntity.OhlcStock, error) {
	key := req.StockCode
	ctx := context.TODO()

	value, err := repo.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var data *ohlcEntity.OhlcStock
	err = json.Unmarshal([]byte(value), &data)
	if err != nil {
		return nil, err
	}

	return data, nil

}
