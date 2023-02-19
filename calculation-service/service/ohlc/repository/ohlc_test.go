package ohlcRepository

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fadellh/stock-ohlc/calculation-service/package/config"
	"github.com/fadellh/stock-ohlc/calculation-service/package/manager"
	ohlcEntity "github.com/fadellh/stock-ohlc/calculation-service/service/ohlc/entity"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
)

func TestNewOhlcRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fmgr, err := manager.NewFakeInit(ctrl)
	if err != nil {
		t.Fatalf("Error create manager instance: err %+v", err)
	}

	h := new(repository)
	h.cfg = fmgr.GetConfig()
	h.redis = fmgr.GetRedis()

	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name string
		args args
		want OhlcRepository
	}{
		{
			name: "success",
			args: args{
				mgr: fmgr,
			},
			want: h,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOhlcRepository(tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOhlcRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_GetRedis(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	rd, mock := redismock.NewClientMock()

	req := ohlcEntity.OhlcStock{
		StockCode:     "TLKM",
		PreviousPrice: 1000,
	}

	data := &ohlcEntity.OhlcStock{
		StockCode:     "TLKM",
		PreviousPrice: 1000,
	}

	src, _ := json.Marshal(data)

	result := string(src[:])

	type fields struct {
		redis *redis.Client
		cfg   *config.Config
	}
	type args struct {
		req ohlcEntity.OhlcStock
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedStatus *redismock.ExpectedString
		want           *ohlcEntity.OhlcStock
		wantErr        bool
	}{
		{
			name: "success",
			fields: fields{
				redis: rd,
				cfg:   &config.Config{},
			},
			args: args{
				req: req,
			},
			expectedStatus: mock.ExpectGet(req.StockCode),
			want:           data,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository{
				redis: tt.fields.redis,
				cfg:   tt.fields.cfg,
			}
			tt.expectedStatus.SetVal(result)

			got, err := repo.GetRedis(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.GetRedis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.GetRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_StoreRedis(t *testing.T) {
	type fields struct {
		redis *redis.Client
		cfg   *config.Config
	}
	type args struct {
		req ohlcEntity.OhlcStock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &repository{
				redis: tt.fields.redis,
				cfg:   tt.fields.cfg,
			}
			if err := repo.StoreRedis(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("repository.StoreRedis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
