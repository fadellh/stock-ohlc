package ohlc

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	pb "github.com/fadellh/stock-ohlc/summary-service/proto"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

func TestHandler_GetOhlcSummary(t *testing.T) {
	rd, mock := redismock.NewClientMock()

	req := &pb.SummaryRequest{
		Stockcode: "TLKM",
	}

	ctx := context.TODO()

	data := &OhlcStock{
		StockCode:     "TLKM",
		PreviousPrice: 1000,
		Value:         "0",
	}

	src, _ := json.Marshal(data)

	result := string(src[:])

	type fields struct {
		UnimplementedOhlcServer pb.UnimplementedOhlcServer
		redisClient             *redis.Client
	}
	type args struct {
		ctx context.Context
		in  *pb.SummaryRequest
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *pb.SummaryResponse
		expectedStatus *redismock.ExpectedString
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			fields: fields{
				UnimplementedOhlcServer: pb.UnimplementedOhlcServer{},
				redisClient:             rd,
			},
			args: args{
				ctx: ctx,
				in:  req,
			},
			expectedStatus: mock.ExpectGet("TLKM"),
			want: &pb.SummaryResponse{
				Prev:    1000,
				Open:    0,
				Highest: 0,
				Lowest:  0,
				Close:   0,
				Average: 0,
				Volume:  0,
				Value:   0,
				Code:    "TLKM",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				UnimplementedOhlcServer: tt.fields.UnimplementedOhlcServer,
				redisClient:             tt.fields.redisClient,
			}

			tt.expectedStatus.SetVal(result)

			got, err := h.GetOhlcSummary(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler.GetOhlcSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.GetRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}
