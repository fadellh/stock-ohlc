package kafkaPackage

import (
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/rs/zerolog/log"
)

type Kafka interface {
	Connect() error
	Consume(topics []string, signals chan os.Signal)
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
	defer func() {
		if err := consumers.Close(); err != nil {
			log.Fatal().Msg(err.Error())
			return
		}
	}()

	o.consumer = consumers
	return nil
}

func (o *Options) Consume(topics []string, signals chan os.Signal) {
	chanMessage := make(chan *sarama.ConsumerMessage, 256)

	for _, topic := range topics {
		partitionList, err := o.consumer.Partitions(topic)
		if err != nil {
			log.Error().Err(err).Msgf("Unable to get partition got error %v", err)
			continue
		}
		for _, partition := range partitionList {
			go consumeMessage(o.consumer, topic, partition, chanMessage)
		}
	}
	log.Info().Msgf("Kafka is consuming....")

ConsumerLoop:
	for {
		select {
		case msg := <-chanMessage:
			log.Info().Msgf("New Message from kafka, message: %v", string(msg.Value))
		case sig := <-signals:
			if sig == os.Interrupt {
				break ConsumerLoop
			}
		}
	}
}

func consumeMessage(consumer sarama.Consumer, topic string, partition int32, c chan *sarama.ConsumerMessage) {
	msg, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to consume partition %v got error %v", partition, err)
		return
	}

	defer func() {
		if err := msg.Close(); err != nil {
			log.Error().Err(err).Msgf("Unable to close partition %v: %v", partition, err)
		}
	}()

	for {
		msg := <-msg.Messages()
		c <- msg
	}

}
