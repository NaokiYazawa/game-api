// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/domain/repository/db/collection_item/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	collectionitem "game-api/pkg/domain/model/collection_item"
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

// SelectAll mocks base method.
func (m *MockRepository) SelectAll() ([]*collectionitem.CollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectAll")
	ret0, _ := ret[0].([]*collectionitem.CollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectAll indicates an expected call of SelectAll.
func (mr *MockRepositoryMockRecorder) SelectAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectAll", reflect.TypeOf((*MockRepository)(nil).SelectAll))
}

// SelectByCollectionIDs mocks base method.
func (m *MockRepository) SelectByCollectionIDs(collectionItemIDs []string) ([]*collectionitem.CollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByCollectionIDs", collectionItemIDs)
	ret0, _ := ret[0].([]*collectionitem.CollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByCollectionIDs indicates an expected call of SelectByCollectionIDs.
func (mr *MockRepositoryMockRecorder) SelectByCollectionIDs(collectionItemIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByCollectionIDs", reflect.TypeOf((*MockRepository)(nil).SelectByCollectionIDs), collectionItemIDs)
}