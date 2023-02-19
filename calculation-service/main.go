package main

import (
	"fmt"
	"os"

	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcService "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc"
)

func run() error {
	mgr, err := manager.NewInit()
	if err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	err = ohlcService.NewOHLC(mgr, signals)
	if err != nil {
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
