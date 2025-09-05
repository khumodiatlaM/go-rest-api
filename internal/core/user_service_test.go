package core

import (
	"context"
	"go-rest-api/pkg/logger"
	"testing"

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
