package core

import (
	"context"
	"go-rest-api/pkg/logger"
	"testing"
	"time"

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
	mockUserRepo.On("CreateUser", mock.Anything, &testUser).Return(nil)

	// when
	err := userService.CreateUser(context.Background(), &testUser)

	// then
	a.NoError(err)
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
