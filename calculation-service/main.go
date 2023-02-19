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

	_, err = ohlcService.NewOHLC(mgr)
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
