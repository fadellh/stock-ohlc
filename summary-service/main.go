package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/fadellh/stock-ohlc/summary-service/config"
	pb "github.com/fadellh/stock-ohlc/summary-service/proto"
	redisPackage "github.com/fadellh/stock-ohlc/summary-service/redis"
	handleOhlc "github.com/fadellh/stock-ohlc/summary-service/service/ohlc"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	UnimplementedOhlcServer pb.UnimplementedOhlcServer
	redisClient             redis.Client
}

func run() error {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
		return err
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
		return err
	}

	redisClient, err := redisPackage.NewRedis(cfg).Connect()
	if err != nil {
		log.Fatal().Err(err).Msgf(err.Error())
		return err
	}

	handlerOhlc := handleOhlc.New(redisClient)

	s := grpc.NewServer()
	pb.RegisterOhlcServer(s, handlerOhlc)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
