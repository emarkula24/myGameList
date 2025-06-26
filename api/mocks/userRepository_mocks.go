package mocks

import (
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) SelectUserByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *UserRepositoryMock) InsertUser(username, email, hashedPassword string) (int64, error) {
	args := m.Called(username, email, hashedPassword)
	return args.Get(0).(int64), args.Error(1)
}

func (m *UserRepositoryMock) PasswordByUsername(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func (m *UserRepositoryMock) SelectUserIdByUsername(username string) (int, error) {
	args := m.Called(username)
	return args.Int(0), args.Error(1)
}

func (m *UserRepositoryMock) InsertRefreshToken(userId int, refreshToken string, jti string) error {
	args := m.Called(userId, refreshToken, jti)
	return args.Error(0)
}

func (m *UserRepositoryMock) RefreshTokenById(userId int) (string, string, error) {
	args := m.Called(userId)
	return args.String(0), args.String(1), args.Error(2)
}
