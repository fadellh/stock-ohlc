package manager

import (
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	kafkaPackage "github.com/fadellh/stock-ohlc/calculation-service/package/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
)

type FakeManager interface {
	Manager
}

type fakeManager struct {
	config *config.Config
	kafka  kafkaPackage.Kafka
	redis  *redis.Client
}

func NewFakeInit(ctrl *gomock.Controller) (FakeManager, error) {
	cfg := &config.Config{}
	kafka := kafkaPackage.NewMockKafka(ctrl)
	redis, _ := redismock.NewClientMock()

	return &fakeManager{
		config: cfg,
		kafka:  kafka,
		redis:  redis,
	}, nil
}
func (fm *fakeManager) GetConfig() *config.Config {
	return fm.config
}

func (fm *fakeManager) GetKafka() kafkaPackage.Kafka {
	return fm.kafka
}

func (fm *fakeManager) GetRedis() *redis.Client {
	return fm.redis
}
