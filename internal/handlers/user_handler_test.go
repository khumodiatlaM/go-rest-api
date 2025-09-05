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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	userHandler *UserHandler
	tearDown    func()
}

func (testSuite *UserHandlerTestSuite) SetupSuite() {
	ctx := context.Background()
	t := testSuite.T()
	dbPool, teardown := test.CreateDbTestContainer(ctx, t)
	testSuite.tearDown = teardown

	mockLogger := logger.MockLogger{}
	mockLogger.On("Error", mock.Anything).Return()
	userRepo := db.NewUserRepository(dbPool, &mockLogger)
	userServ := core.NewUserService(userRepo, &mockLogger)
	userHandler := NewUserHandler(userServ)
	testSuite.userHandler = userHandler
}

func (testSuite *UserHandlerTestSuite) TearDownSuite() {
	if testSuite.tearDown != nil {
		testSuite.tearDown()
	}
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

//func (testSuite *UserHandlerTestSuite) TestGetUser() {
//	t := testSuite.T()
//	a := assert.New(t)
//
//	// given
//	router := httprouter.New()
//	path := "/users/:id"
//	router.GET(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//		testSuite.userHandler.GetUser(w, r, ps)
//	})
//	// First, create a user to ensure there is one to retrieve
//	user := CreateUserRequest{
//		Username: "testuser2",
//		Email:    "test2@user.com",
//		Password: "password123",
//	}
//	reqBody, err := json.Marshal(user)
//	a.NoError(err)
//
//	// ... create user
//	createReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
//	createReq.Header.Set("Content-Type", "application/json")
//	createRes := httptest.NewRecorder()
//	router.POST("/users", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//		testSuite.userHandler.CreateUser(w, r)
//	})
//
//	// when ... we get the user by ID
//	getReq := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
//	getReq.Header.Set("Content-Type", "application/json")
//	getRes := httptest.NewRecorder()
//	router.ServeHTTP(createRes, createReq)
//	router.ServeHTTP(getRes, getReq)
//
//	// then
//	a.Equal(http.StatusOK, getRes.Code)
//    var resultBody GetUserResponse
//	err = json.Unmarshal(getRes.Body.Bytes(), &resultBody)
//	a.NoError(err)
//	if diff := cmp.Diff(resultBody.Username, user.Username); diff != "" {
//		t.Error(diff)
//	}
//}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
