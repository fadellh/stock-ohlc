package ohlcUsecase

import (
	"reflect"
	"testing"

	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	"github.com/golang/mock/gomock"
)

func TestNewOhlcUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fmgr, _ := manager.NewFakeInit(ctrl)

	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name string
		args args
		want OhlcUsecase
	}{
		{
			name: "success",
			args: args{
				mgr: fmgr,
			},
			want: NewOhlcUsecase(fmgr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				if got := NewOhlcUsecase(tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewOhlcUsecase() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

// func Test_ohlcCalculation(t *testing.T) {
// 	type args struct {
// 		req ohlcEntity.OhlcMessage
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want ohlcEntity.OhlcStock
// 	}{
// 		{
// 			name: "type A",
// 			args: args{
// 				req: ohlcEntity.OhlcMessage{
// 					OrderBook: "23",
// 					Price:     7000,
// 					StockCode: "TLKM",
// 					Type:      "A",
// 					Quantity:  0,
// 				},
// 			},
// 			want: ohlcEntity.OhlcStock{
// 				StockCode:     "TLKM",
// 				PreviousPrice: 7000,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := ohlcCalculation(tt.args.req); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("ohlcCalculation() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
