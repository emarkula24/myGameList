package service_test

import (
	"errors"
	"testing"

	"example.com/mygamelist/errorutils"
	"example.com/mygamelist/mocks"
	"example.com/mygamelist/service"
	"github.com/stretchr/testify/assert"
)

// MockHasher implements interfaces.PasswordHasher
type MockHasher struct {
	HashResult string
	HashError  error
}

func (m MockHasher) HashPassword(password string) (string, error) {
	return m.HashResult, m.HashError
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name             string
		username         string
		email            string
		password         string
		setupMocks       func(repo *mocks.UserRepositoryMock)
		mockHashResult   string
		mockHashError    error
		expectedUserID   int64
		expectedErrorStr string
	}{
		{
			name:     "user already exists",
			username: "existing_user",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("SelectUserByUsername", "existing_user").Return(true, nil).Once()
			},
			mockHashResult:   "should-not-be-called",
			expectedErrorStr: errorutils.ErrUserExists.Error(),
		},
		{
			name:     "error hashing password",
			username: "new_user",
			password: "badpass",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("SelectUserByUsername", "new_user").Return(false, nil).Once()
			},
			mockHashError:    errors.New("hash error"),
			expectedErrorStr: "failed to hash password: hash error",
		},
		{
			name:     "insert user fails",
			username: "new_user",
			email:    "email@example.com",
			password: "goodpass",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("SelectUserByUsername", "new_user").Return(false, nil).Once()
				repo.On("InsertUser", "new_user", "email@example.com", "hashed-password").
					Return(int64(0), errors.New("db error")).Once()
			},
			mockHashResult:   "hashed-password",
			expectedErrorStr: "failed to insert user: db error",
		},
		{
			name:     "successful registration",
			username: "new_user",
			email:    "email@example.com",
			password: "securepass",
			setupMocks: func(repo *mocks.UserRepositoryMock) {
				repo.On("SelectUserByUsername", "new_user").Return(false, nil).Once()
				repo.On("InsertUser", "new_user", "email@example.com", "hashed-password").
					Return(int64(42), nil).Once()
			},
			mockHashResult: "hashed-password",
			expectedUserID: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.UserRepositoryMock)

			if tt.setupMocks != nil {
				tt.setupMocks(mockRepo)
			}

			mockHasher := MockHasher{
				HashResult: tt.mockHashResult,
				HashError:  tt.mockHashError,
			}

			svc := service.NewUserService(mockRepo, mockHasher)

			userID, err := svc.RegisterUser(tt.username, tt.email, tt.password)

			if tt.expectedErrorStr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorStr)
				assert.Equal(t, int64(0), userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUserID, userID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
