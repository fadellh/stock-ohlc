run-rg:
	@go run request-generator/main.go

run-cs:
	@go run calculation-service/main.go

test:
	@go test ./...

test-coverage:
	@go test -v -coverprofile cover.out ./...
	@go tool cover -html=cover.out -o cover.html
	@open cover.html

generate-mock-cs:
	@echo "GENERATING ..."
	@echo "[Package]"
	@echo "- Kafka"
	@mockgen -destination=calculation-service/package/kafka/kafka_mock.go -package=kafkaPackage -source=calculation-service/package/kafka/kafka.go
	@echo "[Service]"
	@echo "- ohlc"
	@mockgen -destination=calculation-service/service/ohlc/usecase/ohlc_mock.go -package=ohlcUsecase -source=calculation-service/service/ohlc/usecase/ohlc.go
	@mockgen -destination=calculation-service/service/ohlc/repository/ohlc_mock.go -package=ohlcRepository -source=calculation-service/service/ohlc/repository/ohlc.go
