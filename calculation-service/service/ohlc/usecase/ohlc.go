package ohlcUsecase

import (
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
)

type OhlcUsecase interface {
}

type Options struct {
	config *config.Config
}

func NewOhlcUsecase(mgr manager.Manager) (OhlcUsecase, error) {
	opt := new(Options)
	opt.config = mgr.GetConfig()

	return opt, nil
}
