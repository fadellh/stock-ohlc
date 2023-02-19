package manager

import (
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	kafkaPackage "github.com/fadellh/stock-ohlc/calculation-service/package/kafka"
	redisPackage "github.com/fadellh/stock-ohlc/calculation-service/package/redis"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type Manager interface {
	GetConfig() *config.Config
	GetKafka() kafkaPackage.Kafka
	GetRedis() *redis.Client
}

type manager struct {
	config *config.Config
	kafka  kafkaPackage.Kafka
	redis  *redis.Client
}

func NewInit() (Manager, error) {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error().Err(err).Msgf(err.Error())
		return nil, err
	}

	kafka := kafkaPackage.NewKafka(cfg)
	err = kafka.Connect()
	if err != nil {
		log.Error().Err(err).Msgf(err.Error())
		return nil, err
	}

	redis, err := redisPackage.NewRedis(cfg).Connect()
	if err != nil {
		log.Error().Err(err).Msgf(err.Error())
		return nil, err
	}

	return &manager{
		config: cfg,
		kafka:  kafka,
		redis:  redis,
	}, nil
}

func (sm *manager) GetConfig() *config.Config {
	return sm.config
}

func (sm *manager) GetKafka() kafkaPackage.Kafka {
	return sm.kafka
}
func (sm *manager) GetRedis() *redis.Client {
	return sm.redis
}
