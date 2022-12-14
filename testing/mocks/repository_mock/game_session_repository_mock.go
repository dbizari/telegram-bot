// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/game_session.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	domain "tdl/internal/domain/game_session"

	gomock "github.com/golang/mock/gomock"
)

// MockGameSessionRepositoryAPI is a mock of GameSessionRepositoryAPI interface.
type MockGameSessionRepositoryAPI struct {
	ctrl     *gomock.Controller
	recorder *MockGameSessionRepositoryAPIMockRecorder
}

// MockGameSessionRepositoryAPIMockRecorder is the mock recorder for MockGameSessionRepositoryAPI.
type MockGameSessionRepositoryAPIMockRecorder struct {
	mock *MockGameSessionRepositoryAPI
}

// NewMockGameSessionRepositoryAPI creates a new mock instance.
func NewMockGameSessionRepositoryAPI(ctrl *gomock.Controller) *MockGameSessionRepositoryAPI {
	mock := &MockGameSessionRepositoryAPI{ctrl: ctrl}
	mock.recorder = &MockGameSessionRepositoryAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameSessionRepositoryAPI) EXPECT() *MockGameSessionRepositoryAPIMockRecorder {
	return m.recorder
}

// CreateGame mocks base method.
func (m *MockGameSessionRepositoryAPI) CreateGame(ctx context.Context, gameSession *domain.GameSession) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGame", ctx, gameSession)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGame indicates an expected call of CreateGame.
func (mr *MockGameSessionRepositoryAPIMockRecorder) CreateGame(ctx, gameSession interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGame", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).CreateGame), ctx, gameSession)
}

// ExitGame mocks base method.
func (m *MockGameSessionRepositoryAPI) ExitGame(ctx context.Context, userName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExitGame", ctx, userName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExitGame indicates an expected call of ExitGame.
func (mr *MockGameSessionRepositoryAPIMockRecorder) ExitGame(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExitGame", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).ExitGame), ctx, userName)
}

// Get mocks base method.
func (m *MockGameSessionRepositoryAPI) Get(ctx context.Context, gameSessionID string) (*domain.GameSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, gameSessionID)
	ret0, _ := ret[0].(*domain.GameSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockGameSessionRepositoryAPIMockRecorder) Get(ctx, gameSessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).Get), ctx, gameSessionID)
}

// GetByMember mocks base method.
func (m *MockGameSessionRepositoryAPI) GetByMember(ctx context.Context, username string) (*domain.GameSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByMember", ctx, username)
	ret0, _ := ret[0].(*domain.GameSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByMember indicates an expected call of GetByMember.
func (mr *MockGameSessionRepositoryAPIMockRecorder) GetByMember(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByMember", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).GetByMember), ctx, username)
}

// GetNotFinishedGameByMember mocks base method.
func (m *MockGameSessionRepositoryAPI) GetNotFinishedGameByMember(ctx context.Context, userID string) (*domain.GameSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotFinishedGameByMember", ctx, userID)
	ret0, _ := ret[0].(*domain.GameSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotFinishedGameByMember indicates an expected call of GetNotFinishedGameByMember.
func (mr *MockGameSessionRepositoryAPIMockRecorder) GetNotFinishedGameByMember(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotFinishedGameByMember", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).GetNotFinishedGameByMember), ctx, userID)
}

// Update mocks base method.
func (m *MockGameSessionRepositoryAPI) Update(ctx context.Context, gameSession *domain.GameSession) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, gameSession)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockGameSessionRepositoryAPIMockRecorder) Update(ctx, gameSession interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGameSessionRepositoryAPI)(nil).Update), ctx, gameSession)
}
