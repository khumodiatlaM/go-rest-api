package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"go-rest-api/internal/core"
	"go-rest-api/internal/db"
	"go-rest-api/pkg/logger"
	"go-rest-api/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	userHandler *UserHandler
}

func (testSuite *UserHandlerTestSuite) SetupSuite() {
	ctx := context.Background()
	t := testSuite.T()
	dbPool, teardown := test.CreateDbTestContainer(ctx, t)
	defer teardown()

	mockLogger := logger.MockLogger{}
	userRepo := db.NewUserRepository(dbPool, &mockLogger)
	userServ := core.NewUserService(userRepo, &mockLogger)
	userHandler := NewUserHandler(userServ)
	testSuite.userHandler = userHandler
}

func (testSuite *UserHandlerTestSuite) TestCreateUser() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	router := httprouter.New()
	path := "/users"
	router.POST(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.CreateUser(w, r)
	})
	user := CreateUserRequest{
		Username: "testuser",
		Email:    "test@user.com",
		Password: "password123",
	}
	reqBody, err := json.Marshal(user)
	a.NoError(err)

	// when
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// then
	a.Equal(http.StatusCreated, res.Code)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
