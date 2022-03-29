// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/domain/repository/db/user_collection_item/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	usercollectionitem "game-api/pkg/domain/model/user_collection_item"
	context "context"
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

// InsertWithLock mocks base method.
func (m *MockRepository) InsertWithLock(ctx context.Context, userCollectionItems []*usercollectionitem.UserCollectionItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertWithLock", ctx, userCollectionItems)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertWithLock indicates an expected call of InsertWithLock.
func (mr *MockRepositoryMockRecorder) InsertWithLock(ctx, userCollectionItems interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertWithLock", reflect.TypeOf((*MockRepository)(nil).InsertWithLock), ctx, userCollectionItems)
}

// SelectByUserID mocks base method.
func (m *MockRepository) SelectByUserID(userID string) ([]*usercollectionitem.UserCollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByUserID", userID)
	ret0, _ := ret[0].([]*usercollectionitem.UserCollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByUserID indicates an expected call of SelectByUserID.
func (mr *MockRepositoryMockRecorder) SelectByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByUserID", reflect.TypeOf((*MockRepository)(nil).SelectByUserID), userID)
}

// SelectByUserIDAndCollectionIDs mocks base method.
func (m *MockRepository) SelectByUserIDAndCollectionIDs(userID string, collectionItemIDs []string) ([]*usercollectionitem.UserCollectionItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByUserIDAndCollectionIDs", userID, collectionItemIDs)
	ret0, _ := ret[0].([]*usercollectionitem.UserCollectionItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByUserIDAndCollectionIDs indicates an expected call of SelectByUserIDAndCollectionIDs.
func (mr *MockRepositoryMockRecorder) SelectByUserIDAndCollectionIDs(userID, collectionItemIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByUserIDAndCollectionIDs", reflect.TypeOf((*MockRepository)(nil).SelectByUserIDAndCollectionIDs), userID, collectionItemIDs)
}