package ohlcUsecase

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcEntity "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/entity"
	ohlcRepository "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/repository"
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

func Test_ohlcTrxCalculation(t *testing.T) {

	currentData := ohlcEntity.OhlcStock{
		StockCode:     "TLKM",
		PreviousPrice: 500,
		OpenPrice:     0,
		HighestPrice:  0,
		LowestPrice:   0,
		ClosePrice:    1000,
		Volume:        0,
		Value:         "0",
		Type:          "P",
	}

	currentDataMass := ohlcEntity.OhlcStock{
		StockCode:     "BBRI",
		PreviousPrice: 8000,
		OpenPrice:     8050,
		HighestPrice:  8100,
		LowestPrice:   7950,
		ClosePrice:    8100,
		Volume:        900,
		Value:         "7210000",
		Type:          "P",
		AveragePrice:  8011,
	}

	currentValue := 7210000
	var amount_of_two_value int64 = int64(currentValue + (4333 * 225))
	strValue := strconv.FormatInt(amount_of_two_value, 10)

	type args struct {
		req         ohlcEntity.OhlcMessage
		currentOhlc ohlcEntity.OhlcStock
	}
	tests := []struct {
		name string
		args args
		want ohlcEntity.OhlcStock
	}{
		{
			name: "different stock code",
			args: args{
				req: ohlcEntity.OhlcMessage{
					OrderBook: "",
					Price:     0,
					StockCode: "TLKM",
					Type:      "",
					Quantity:  0,
				},
				currentOhlc: ohlcEntity.OhlcStock{
					StockCode:     "BBRI",
					PreviousPrice: 0,
					OpenPrice:     0,
					HighestPrice:  0,
					LowestPrice:   0,
					ClosePrice:    0,
					Volume:        0,
					Value:         "0",
					AveragePrice:  0,
					Type:          "",
				},
			},
			want: ohlcEntity.OhlcStock{},
		},
		{
			name: "open price",
			args: args{
				req: ohlcEntity.OhlcMessage{
					OrderBook: "32",
					Price:     1000,
					StockCode: "TLKM",
					Quantity:  5,
					Type:      "P",
				},
				currentOhlc: currentData,
			},
			want: ohlcEntity.OhlcStock{
				StockCode:     "TLKM",
				PreviousPrice: 500,
				OpenPrice:     1000,
				HighestPrice:  1000,
				LowestPrice:   1000,
				ClosePrice:    1000,
				Volume:        5,
				Value:         "5000",
				AveragePrice:  1000,
				Type:          "P",
			},
		},
		{
			name: "amount of 2 value and lowest price",
			args: args{
				req: ohlcEntity.OhlcMessage{
					OrderBook: "33",
					Price:     4333,
					StockCode: "BBRI",
					Type:      "P",
					Quantity:  225,
				},
				currentOhlc: currentDataMass,
			},
			want: ohlcEntity.OhlcStock{
				StockCode:     "BBRI",
				PreviousPrice: currentDataMass.PreviousPrice,
				OpenPrice:     currentDataMass.OpenPrice,
				HighestPrice:  currentDataMass.HighestPrice,
				LowestPrice:   4333,
				ClosePrice:    4333,
				Volume:        1125,
				Value:         strValue,
				AveragePrice:  7275,
				Type:          "P",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ohlcTrxCalculation(tt.args.req, tt.args.currentOhlc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ohlcTrxCalculation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_CalculateOHLC(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := ohlcRepository.NewMockOhlcRepository(ctrl)

	messageTypeA := []byte(
		`
		{"order_book":"35","price":"4540","stock_code":"UNVR","type":"A"}
		`)
	reqTypeA := ohlcEntity.OhlcStock{
		StockCode:     "UNVR",
		PreviousPrice: 4540,
		OpenPrice:     0,
		HighestPrice:  0,
		LowestPrice:   0,
		ClosePrice:    0,
		Volume:        0,
		Value:         "0",
		AveragePrice:  0,
	}

	messageTypeP := []byte(
		`
		{"type":"P","executed_quantity":"33","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
		`)

	reqTypeP := ohlcEntity.OhlcStock{
		StockCode:     "UNVR",
		PreviousPrice: 0,
		OpenPrice:     0,
		HighestPrice:  0,
		LowestPrice:   0,
		ClosePrice:    0,
		Volume:        0,
		Value:         "0",
		AveragePrice:  0,
		Type:          "",
	}

	errMsg := []byte(
		`
			{"type:"P","executed_quantity":"33","order_book":"35","execution_price":"4530","stock_code":"UNVR"}
			`)

	resOhlcRedis := &ohlcEntity.OhlcStock{
		StockCode:     "UNVR",
		PreviousPrice: 5000,
		OpenPrice:     3000,
		HighestPrice:  3500,
		LowestPrice:   2000,
		ClosePrice:    2000,
		Volume:        100,
		Value:         "500000",
		AveragePrice:  5000,
		Type:          "P",
	}

	type fields struct {
		config    *config.Config
		consumer  sarama.Consumer
		OhlcTopic string
		repo      ohlcRepository.OhlcRepository
	}
	type args struct {
		msg *sarama.ConsumerMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		call1   *gomock.Call
		call2   *gomock.Call
		call3   *gomock.Call
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:   "success opening stock",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: messageTypeA,
			}},
			call1: mockRepo.EXPECT().GetRedis(reqTypeA).Return(nil, nil),
			call2: mockRepo.EXPECT().StoreRedis(reqTypeA).Return(nil),
		},
		{
			name:   "success redis exist",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: messageTypeP,
			}},
			call1: mockRepo.EXPECT().GetRedis(reqTypeP).Return(resOhlcRedis, nil),
			call2: mockRepo.EXPECT().StoreRedis(gomock.Any()).Return(nil),
		},
		{
			name:   "err unmarshal",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: errMsg,
			}},
		},
		{
			name:   "err get redis",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: messageTypeP,
			}},
			call1: mockRepo.EXPECT().GetRedis(reqTypeP).Return(nil, errors.New("err")),
		},
		{
			name:   "err store redis when not exist",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: messageTypeA,
			}},
			call1: mockRepo.EXPECT().GetRedis(reqTypeA).Return(nil, nil),
			call2: mockRepo.EXPECT().StoreRedis(reqTypeA).Return(errors.New("err")),
		},
		{
			name:   "err store redis when exist",
			fields: fields{config: &config.Config{}, consumer: nil, OhlcTopic: "", repo: mockRepo},
			args: args{msg: &sarama.ConsumerMessage{
				Value: messageTypeP,
			}},
			call1: mockRepo.EXPECT().GetRedis(reqTypeP).Return(resOhlcRedis, nil),
			call2: mockRepo.EXPECT().StoreRedis(gomock.Any()).Return(errors.New("err")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{
				config:    tt.fields.config,
				consumer:  tt.fields.consumer,
				OhlcTopic: tt.fields.OhlcTopic,
				repo:      tt.fields.repo,
			}
			o.CalculateOHLC(tt.args.msg)
		})
	}
}
