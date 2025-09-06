package core

import (
	"context"
	"go-rest-api/pkg/logger"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	testUser := User{
		ID:       uuid.New(),
		Username: "JohnDoe123",
		Email:    "johndoe@gmail.com",
		Password: "hashedpassword",
	}
	mockUserRepo.On("CreateUser", mock.Anything, &testUser).Return(&testUser, nil)

	// when
	user, err := userService.CreateUser(context.Background(), &testUser)

	// then
	a.NoError(err)
	if diff := cmp.Diff(testUser, *user); diff != "" {
		t.Error(diff)
	}
}

func TestUserService_CreateUser_ReturnsError(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	testUser := User{
		ID:       uuid.New(),
		Username: "JohnDoe123",
		Email:    "johndoe@gmail.com",
		Password: "hashedpassword",
	}
	mockUserRepo.On("CreateUser", mock.Anything, &testUser).Return(&User{}, assert.AnError)

	// when
	_, err := userService.CreateUser(context.Background(), &testUser)

	// then
	a.Error(err)
}

func TestUserService_GetUser(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	testUser := User{
		ID:        uuid.New(),
		Username:  "JohnDoe123",
		Email:     "johnDoe@gmail.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserRepo.On("GetUserByID", mock.Anything, testUser.ID.String()).Return(&testUser, nil)

	// when
	user, err := userService.GetUserByID(context.Background(), testUser.ID.String())

	// then
	a.NoError(err)
	a.Equal(&testUser, user)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	mockUserRepo.On("GetUserByID", mock.Anything, "non-existent-id").Return(nil, nil)

	// when
	user, err := userService.GetUserByID(context.Background(), "non-existent-id")

	// then
	a.NoError(err)
	a.Nil(user)
}

func TestUserService_GetUser_ReturnAnError(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	mockUserRepo.On("GetUserByID", mock.Anything, "some-id").Return(nil, assert.AnError)

	// when
	user, err := userService.GetUserByID(context.Background(), "some-id")

	// then
	a.Error(err)
	a.Nil(user)
}

func TestUserService_AuthenticateUser(t *testing.T) {
	a := assert.New(t)

	// given
	mockLogger := logger.MockLogger{}
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	hashedPassword, err := HashPassword("password")
	testUser := User{
		ID:        uuid.New(),
		Username:  "JohnDoe13",
		Email:     "JohnDoe13@gamil.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserRepo.On("GetUserByEmail", mock.Anything, testUser.Email).Return(&testUser, nil)

	// when
	user, err := userService.AuthenticateUser(context.Background(), testUser.Email, "password")

	// then
	a.NoError(err)
	if diff := cmp.Diff(testUser, *user); diff != "" {
		t.Error(diff)
	}
}

func TestUserService_AuthenticateUser_WrongPassword(t *testing.T) {
	a := assert.New(t)

	// given
	mockLogger := logger.MockLogger{}
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	hashedPassword, err := HashPassword("password")
	testUser := User{
		ID:        uuid.New(),
		Username:  "JohnDoe13",
		Email:     "JohnDoe13@gamil.com",
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockUserRepo.On("GetUserByEmail", mock.Anything, testUser.Email).Return(&testUser, nil)

	// when
	result, err := userService.AuthenticateUser(context.Background(), testUser.Email, "wrongpassword")

	// then
	a.Error(err)
	a.Nil(result)
}

func TestUserService_AuthenticateUser_ReturnsError(t *testing.T) {
	a := assert.New(t)
	// given
	mockLogger := logger.MockLogger{}
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockUserRepo := MockUserRepository{}
	userService := NewUserService(&mockUserRepo, &mockLogger)

	mockUserRepo.On("GetUserByEmail", mock.Anything, "non-existent-email").Return(&User{}, assert.AnError)

	// when
	result, err := userService.AuthenticateUser(context.Background(), "non-existent-email", "password")

	// then
	a.Error(err)
	a.Nil(result)
}
