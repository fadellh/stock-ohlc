package kafkaPackage

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/rs/zerolog/log"
)

type Kafka interface {
	Connect() error
	Consumer() sarama.Consumer
}

type Options struct {
	kafkaAddr    string
	writeTimeOut int64
	maxRetry     int
	username     string
	password     string
	consumer     sarama.Consumer
}

func NewKafka(cfg *config.Config) Kafka {
	opt := new(Options)
	opt.kafkaAddr = cfg.KafkaAddress
	opt.writeTimeOut = cfg.WriteTimeout
	opt.maxRetry = int(cfg.MaxRetry)

	return opt
}

func (o *Options) Connect() error {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = time.Duration(o.writeTimeOut) * time.Second
	kafkaConfig.Producer.Retry.Max = o.maxRetry

	if o.username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = o.username
		kafkaConfig.Net.SASL.Password = o.password
	}

	consumers, err := sarama.NewConsumer([]string{o.kafkaAddr}, kafkaConfig)
	if err != nil {
		log.Error().Err(err).Msgf("Error create kakfa consumer got error %v", err)
		return err
	}

	o.consumer = consumers
	return nil
}

func (o *Options) Consumer() sarama.Consumer {
	return o.consumer
}
