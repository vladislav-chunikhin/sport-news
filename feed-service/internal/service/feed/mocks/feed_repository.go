// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	feed "github.com/vladislav-chunikhin/feed-service/internal/repository/feed"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockRepository) GetByID(ctx context.Context, ID primitive.ObjectID) (*feed.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, ID)
	ret0, _ := ret[0].(*feed.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockRepositoryMockRecorder) GetByID(ctx, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ctx, ID)
}

// GetLatestArticles mocks base method.
func (m *MockRepository) GetLatestArticles(ctx context.Context, cursor string, limit int) ([]*feed.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestArticles", ctx, cursor, limit)
	ret0, _ := ret[0].([]*feed.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestArticles indicates an expected call of GetLatestArticles.
func (mr *MockRepositoryMockRecorder) GetLatestArticles(ctx, cursor, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestArticles", reflect.TypeOf((*MockRepository)(nil).GetLatestArticles), ctx, cursor, limit)
}