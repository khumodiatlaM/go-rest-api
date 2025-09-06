package db

import (
	"context"
	"go-rest-api/internal/core"
	"go-rest-api/pkg/logger"
	"go-rest-api/test"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	userRepo core.UserRepository
	tearDown func()
}

func (testSuite *UserRepositoryTestSuite) SetupSuite() {
	t := testSuite.T()
	dbPool, tear := test.CreateDbTestContainer(context.Background(), t)
	testSuite.tearDown = tear
	mockLogger := logger.MockLogger{}
	mockLogger.On("Error", mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	testSuite.userRepo = NewUserRepository(dbPool, &mockLogger)
}

func (testSuite *UserRepositoryTestSuite) TearDownSuite() {
	if testSuite.tearDown != nil {
		testSuite.tearDown()
	}
}

func (testSuite *UserRepositoryTestSuite) TestUserRepository_CreateUser() {
	t := testSuite.T()
	a := assert.New(t)
	// given
	testUser := core.User{
		ID:        uuid.New(),
		Username:  "JohnDoe123",
		Email:     "johndoe@gmail.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// when
	createdUser, err := testSuite.userRepo.CreateUser(context.Background(), &testUser)

	// then
	expectedUser := core.User{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Email:     testUser.Email,
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: testUser.UpdatedAt,
	}
	a.NoError(err)
	if diff := cmp.Diff(&expectedUser, createdUser, cmpopts.IgnoreFields(core.User{}, "Password", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Error(diff)
	}
}

func (testSuite *UserRepositoryTestSuite) TestUserRepository_CreateUser_ReturnsError() {
	t := testSuite.T()
	a := assert.New(t)
	// given
	testUser := core.User{
		ID:        uuid.New(),
		Username:  "JohnDoe12536",
		Email:     "johndoe3534@gmail.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// when
	_, err := testSuite.userRepo.CreateUser(context.Background(), &testUser)
	a.NoError(err)

	// ... try to create same user again to trigger unique constraint violation
	_, err = testSuite.userRepo.CreateUser(context.Background(), &testUser)

	// then
	a.Error(err)
}

func (testSuite *UserRepositoryTestSuite) TestUserRepository_GetUserByID() {
	t := testSuite.T()
	a := assert.New(t)
	// given
	testUser := core.User{
		ID:        uuid.New(),
		Username:  "JohnDoe1253633",
		Email:     "johndoe6373@gmail.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ... first create the user
	createdUser, err := testSuite.userRepo.CreateUser(context.Background(), &testUser)
	a.NoError(err)

	// when
	user, err := testSuite.userRepo.GetUserByID(context.Background(), createdUser.ID.String())

	// then
	expectedUser := core.User{
		ID:        testUser.ID,
		Username:  testUser.Username,
		Email:     testUser.Email,
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: testUser.UpdatedAt,
	}
	a.NoError(err)
	if diff := cmp.Diff(&expectedUser, user, cmpopts.IgnoreFields(core.User{}, "Password", "CreatedAt", "UpdatedAt")); diff != "" {
		t.Error(diff)
	}
}

func (testSuite *UserRepositoryTestSuite) TestUserRepository_GetUserByID_UserNotFound() {
	t := testSuite.T()
	a := assert.New(t)
	// when
	user, err := testSuite.userRepo.GetUserByID(context.Background(), uuid.New().String())

	// then
	a.NoError(err)
	a.Nil(user)
}

func (testSuite *UserRepositoryTestSuite) TestUserRepository_GetUserByID_InvalidUUID() {
	t := testSuite.T()
	a := assert.New(t)
	// when
	user, err := testSuite.userRepo.GetUserByID(context.Background(), "invalid-uuid")

	// then
	a.Error(err)
	a.Nil(user)
}

func TestNewUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
