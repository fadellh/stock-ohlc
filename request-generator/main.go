package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/request-generator/publiser"
	"github.com/rs/zerolog/log"
)

type StockTrx struct {
	OrderBook string `json:"order_book"`
	Price     int    `json:"price"`
	StockCode string `json:"stock_code"`
	Type      string `json:"type"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {

	kafkaConfig := getKafkaConfig("", "")
	producers, err := sarama.NewSyncProducer([]string{"kafka:9092"}, kafkaConfig)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to create kafka producer got error %v", err)
		return err
	}
	defer func() {
		if err := producers.Close(); err != nil {
			log.Error().Err(err).Msgf("Unable to stop kafka producer: %v", err)
			return
		}
	}()

	log.Info().Msgf("Success create kafka sync-producer")

	kafka := &publiser.KafkaProducer{
		Producer: producers,
	}

	dirName := "./request-generator/subsetdata/"
	filenames := []string{}
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	for _, filename := range filenames {
		decodeJsonPubliser(dirName+filename, kafka)
	}

	return nil
}

func decodeJsonPubliser(filename string, kafka *publiser.KafkaProducer) error {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msgf("failed to open json file: %s, error: %v", filename, err)
		return err
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Error().Err(err).Msg("failed to read json file, error")
		return err
	}

	d := json.NewDecoder(strings.NewReader(string(jsonData)))
	for {
		var v map[string]interface{}
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg(err.Error())
			}
			break
		}

		payload, _ := json.Marshal(v)
		kafka.SendMessage("stock-ohlc", payload)
	}

	return nil
}

func getKafkaConfig(username, password string) *sarama.Config {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Net.WriteTimeout = 5 * time.Second
	kafkaConfig.Producer.Retry.Max = 0

	if username != "" {
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	}
	return kafkaConfig
}
