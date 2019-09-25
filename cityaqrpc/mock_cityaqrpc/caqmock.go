// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ctessum/cityaq/cityaqrpc (interfaces: CityAQClient)

// Package mock_cityaqrpc is a generated GoMock package.
package mock_cityaqrpc

import (
	context "context"
	cityaqrpc "github.com/ctessum/cityaq/cityaqrpc"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockCityAQClient is a mock of CityAQClient interface
type MockCityAQClient struct {
	ctrl     *gomock.Controller
	recorder *MockCityAQClientMockRecorder
}

// MockCityAQClientMockRecorder is the mock recorder for MockCityAQClient
type MockCityAQClientMockRecorder struct {
	mock *MockCityAQClient
}

// NewMockCityAQClient creates a new mock instance
func NewMockCityAQClient(ctrl *gomock.Controller) *MockCityAQClient {
	mock := &MockCityAQClient{ctrl: ctrl}
	mock.recorder = &MockCityAQClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCityAQClient) EXPECT() *MockCityAQClientMockRecorder {
	return m.recorder
}

// Cities mocks base method
func (m *MockCityAQClient) Cities(arg0 context.Context, arg1 *cityaqrpc.CitiesRequest, arg2 ...grpc.CallOption) (*cityaqrpc.CitiesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Cities", varargs...)
	ret0, _ := ret[0].(*cityaqrpc.CitiesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cities indicates an expected call of Cities
func (mr *MockCityAQClientMockRecorder) Cities(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cities", reflect.TypeOf((*MockCityAQClient)(nil).Cities), varargs...)
}

// CityGeometry mocks base method
func (m *MockCityAQClient) CityGeometry(arg0 context.Context, arg1 *cityaqrpc.CityGeometryRequest, arg2 ...grpc.CallOption) (*cityaqrpc.CityGeometryResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CityGeometry", varargs...)
	ret0, _ := ret[0].(*cityaqrpc.CityGeometryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CityGeometry indicates an expected call of CityGeometry
func (mr *MockCityAQClientMockRecorder) CityGeometry(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CityGeometry", reflect.TypeOf((*MockCityAQClient)(nil).CityGeometry), varargs...)
}

// EmissionsGrid mocks base method
func (m *MockCityAQClient) EmissionsGrid(arg0 context.Context, arg1 *cityaqrpc.EmissionsGridRequest, arg2 ...grpc.CallOption) (*cityaqrpc.EmissionsGridResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EmissionsGrid", varargs...)
	ret0, _ := ret[0].(*cityaqrpc.EmissionsGridResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EmissionsGrid indicates an expected call of EmissionsGrid
func (mr *MockCityAQClientMockRecorder) EmissionsGrid(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmissionsGrid", reflect.TypeOf((*MockCityAQClient)(nil).EmissionsGrid), varargs...)
}

// EmissionsMap mocks base method
func (m *MockCityAQClient) EmissionsMap(arg0 context.Context, arg1 *cityaqrpc.EmissionsMapRequest, arg2 ...grpc.CallOption) (*cityaqrpc.EmissionsMapResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EmissionsMap", varargs...)
	ret0, _ := ret[0].(*cityaqrpc.EmissionsMapResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EmissionsMap indicates an expected call of EmissionsMap
func (mr *MockCityAQClientMockRecorder) EmissionsMap(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmissionsMap", reflect.TypeOf((*MockCityAQClient)(nil).EmissionsMap), varargs...)
}
