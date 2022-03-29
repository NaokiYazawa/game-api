// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/domain/repository/db/ranking/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	ranking "game-api/pkg/domain/model/ranking"
	user "game-api/pkg/domain/model/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// AddRankingList mocks base method.
func (m *MockRepository) AddRankingList(user *user.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRankingList", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRankingList indicates an expected call of AddRankingList.
func (mr *MockRepositoryMockRecorder) AddRankingList(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRankingList", reflect.TypeOf((*MockRepository)(nil).AddRankingList), user)
}

// SelectRankingList mocks base method.
func (m *MockRepository) SelectRankingList(start int64) ([]*ranking.Ranking, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectRankingList", start)
	ret0, _ := ret[0].([]*ranking.Ranking)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectRankingList indicates an expected call of SelectRankingList.
func (mr *MockRepositoryMockRecorder) SelectRankingList(start interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectRankingList", reflect.TypeOf((*MockRepository)(nil).SelectRankingList), start)
}
