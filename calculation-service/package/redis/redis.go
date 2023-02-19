package redisPackage

import (
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/go-redis/redis/v8"
)

type Redis interface {
	Connect() (*redis.Client, error)
}

type Options struct {
	address         string
	username        string
	password        string
	database        int
	redisConnection string
}

func NewRedis(cfg *config.Config) Redis {
	opt := new(Options)
	opt.address = cfg.RedisAddress
	opt.username = cfg.RedisUsername
	opt.password = cfg.RedisPassword
	opt.database = cfg.RedisDatabase
	opt.redisConnection = cfg.RedisConnection

	return opt
}

func (o *Options) Connect() (*redis.Client, error) {
	redisClient, err := redis.ParseURL(o.redisConnection)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(redisClient)

	return rdb, nil
}
