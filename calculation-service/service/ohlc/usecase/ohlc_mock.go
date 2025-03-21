// Code generated by MockGen. DO NOT EDIT.
// Source: calculation-service/service/ohlc/usecase/ohlc.go

// Package ohlcUsecase is a generated GoMock package.
package ohlcUsecase

import (
	reflect "reflect"

	sarama "github.com/Shopify/sarama"
	gomock "github.com/golang/mock/gomock"
)

// MockOhlcUsecase is a mock of OhlcUsecase interface.
type MockOhlcUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockOhlcUsecaseMockRecorder
}

// MockOhlcUsecaseMockRecorder is the mock recorder for MockOhlcUsecase.
type MockOhlcUsecaseMockRecorder struct {
	mock *MockOhlcUsecase
}

// NewMockOhlcUsecase creates a new mock instance.
func NewMockOhlcUsecase(ctrl *gomock.Controller) *MockOhlcUsecase {
	mock := &MockOhlcUsecase{ctrl: ctrl}
	mock.recorder = &MockOhlcUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOhlcUsecase) EXPECT() *MockOhlcUsecaseMockRecorder {
	return m.recorder
}

// CalculateOHLC mocks base method.
func (m *MockOhlcUsecase) CalculateOHLC(arg0 *sarama.ConsumerMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CalculateOHLC", arg0)
}

// CalculateOHLC indicates an expected call of CalculateOHLC.
func (mr *MockOhlcUsecaseMockRecorder) CalculateOHLC(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateOHLC", reflect.TypeOf((*MockOhlcUsecase)(nil).CalculateOHLC), arg0)
}
