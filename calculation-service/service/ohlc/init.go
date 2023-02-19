package ohlc

import (
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcUsecase "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/usecase"
)

type CustomerApp interface{}

type Options struct {
	ohlcUsecase ohlcUsecase.OhlcUsecase
}

func NewOHLC(mgr manager.Manager) (CustomerApp, error) {
	ohlc, err := ohlcUsecase.NewOhlcUsecase(mgr)
	if err != nil {
		return nil, err
	}

	opt := new(Options)
	opt.ohlcUsecase = ohlc

	return opt, nil
}
