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

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	userHandler *UserHandler
	dbPool      *pgxpool.Pool
	tearDown    func()
}

func (testSuite *UserHandlerTestSuite) SetupSuite() {
	ctx := context.Background()
	t := testSuite.T()
	dbPool, teardown := test.CreateDbTestContainer(ctx, t)
	testSuite.dbPool = dbPool
	testSuite.tearDown = teardown

	mockLogger := logger.MockLogger{}
	mockLogger.On("Error", mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything).Return()
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	userRepo := db.NewUserRepository(dbPool, &mockLogger)
	mockUserEvent := core.MockUserEventService{}
	mockUserEvent.On("PublishUserCreatedEvent", mock.Anything, mock.Anything).Return(nil)
	userServ := core.NewUserService(userRepo, &mockLogger, &mockUserEvent)
	userHandler := NewUserHandler(userServ, &mockLogger, "testsecret")
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

func (testSuite *UserHandlerTestSuite) TestCreateUser_InvalidPayload() {

	t := testSuite.T()
	a := assert.New(t)

	// given
	router := httprouter.New()
	path := "/users"
	router.POST(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.CreateUser(w, r)
	})

	reqBody, err := json.Marshal("invalid payload")
	a.NoError(err)

	// when
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// then
	a.Equal(http.StatusBadRequest, res.Code)
}

func (testSuite *UserHandlerTestSuite) TestCreateUser_InvalidRequestBodyData() {
	testScenarios := []struct {
		name string
		user CreateUserRequest
	}{
		{
			name: "missing username",
			user: CreateUserRequest{
				Username: "",
				Email:    "test@gmail.com",
				Password: "password123",
			},
		},
		{
			name: "invalid email",
			user: CreateUserRequest{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
		},
		{
			name: "short password",
			user: CreateUserRequest{
				Username: "testuser",
				Email:    "test@gmail.com",
				Password: "123",
			},
		},
	}

	t := testSuite.T()

	// given
	router := httprouter.New()
	path := "/users"
	router.POST(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.CreateUser(w, r)
	})

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			a := assert.New(t)
			reqBody, err := json.Marshal(scenario.user)
			a.NoError(err)

			// when
			req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			// then
			a.Equal(http.StatusBadRequest, res.Code)
		})
	}
}

func (testSuite *UserHandlerTestSuite) TestGetUser() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	Id := uuid.New()
	router := httprouter.New()
	path := "/users/" + Id.String()
	router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.GetUser(w, r)
	})

	// First, create a user to ensure there is one to retrieve

	username := "testuser3435"
	email := "testuse36373r@gmail.com"
	password := "password123"
	ctx := context.Background()

	const query = `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err := testSuite.dbPool.Exec(ctx, query,
		Id,
		username,
		email,
		password,
	)
	a.NoError(err)

	// when ... we get the user by ID
	getReq := httptest.NewRequest(http.MethodGet, path, nil)
	getReq.Header.Set("Content-Type", "application/json")
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)

	// then
	a.Equal(http.StatusOK, getRes.Code)
	var resultBody UserResponse
	err = json.Unmarshal(getRes.Body.Bytes(), &resultBody)
	expectedResult := UserResponse{
		Id:       Id,
		Username: username,
		Email:    email,
	}
	a.NoError(err)
	if diff := cmp.Diff(expectedResult, resultBody, cmpopts.IgnoreFields(UserResponse{}, "CreatedAt", "UpdatedAt")); diff != "" {
		t.Error(diff)
	}
}

func (testSuite *UserHandlerTestSuite) TestGetUser_NotFound() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	Id := uuid.New()
	router := httprouter.New()
	path := "/users/" + Id.String()
	router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.GetUser(w, r)
	})
	// when ... we get the user by ID
	getReq := httptest.NewRequest(http.MethodGet, path, nil)
	getReq.Header.Set("Content-Type", "application/json")
	getRes := httptest.NewRecorder()
	router.ServeHTTP(getRes, getReq)

	// then
	a.Equal(http.StatusNotFound, getRes.Code)
}

func (testSuite *UserHandlerTestSuite) TestLoginUser() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	Id := uuid.New()
	router := httprouter.New()
	path := "/users/" + Id.String()
	router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.GetUser(w, r)
	})

	// First, create a user to ensure there is one to login

	username := "testuser35"
	email := "testuse373r@gmail.com"
	password := "password123"
	hashedPassword, err := test.HashPassword(password)
	a.NoError(err)
	ctx := context.Background()

	const query = `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err = testSuite.dbPool.Exec(ctx, query,
		Id,
		username,
		email,
		hashedPassword,
	)
	a.NoError(err)

	// when ... we login the user
	loginPath := "/users/login"
	router.POST(loginPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.LoginUser(w, r)
	})
	loginReqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	a.NoError(err)

	loginReq := httptest.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)

	// then
	a.Equal(http.StatusOK, loginRes.Code)
	var resultBody LoginUserResponse
	err = json.Unmarshal(loginRes.Body.Bytes(), &resultBody)
	a.NoError(err)
	a.NotEmpty(resultBody.Token)
}

func (testSuite *UserHandlerTestSuite) TestLoginUser_ReturnsError() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	Id := uuid.New()
	router := httprouter.New()
	path := "/users/" + Id.String()
	router.GET(path, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.GetUser(w, r)
	})

	// First, create a user to ensure there is one to login

	username := "testuser305"
	email := "testuse3093r@gmail.com"
	password := "password123"
	hashedPassword, err := test.HashPassword(password)
	a.NoError(err)
	ctx := context.Background()

	const query = `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err = testSuite.dbPool.Exec(ctx, query,
		Id,
		username,
		email,
		hashedPassword,
	)
	a.NoError(err)

	// when ... we login the user
	loginPath := "/users/login"
	router.POST(loginPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.LoginUser(w, r)
	})
	loginReqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": "wrongpassword",
	})
	a.NoError(err)

	loginReq := httptest.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)

	// then
	a.Equal(http.StatusInternalServerError, loginRes.Code)
}

func (testSuite *UserHandlerTestSuite) TestLoginUser_NotFound() {
	t := testSuite.T()
	a := assert.New(t)

	// given

	router := httprouter.New()
	loginPath := "/users/login"
	router.POST(loginPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.LoginUser(w, r)
	})

	loginReqBody, err := json.Marshal(map[string]string{
		"email":    "test@gmail.com",
		"password": "password123",
	})
	a.NoError(err)

	// when ... we login the user
	loginReq := httptest.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)

	// then
	a.Equal(http.StatusUnauthorized, loginRes.Code)
}

func (testSuite *UserHandlerTestSuite) TestLoginUser_InvalidPayload() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	router := httprouter.New()
	loginPath := "/users/login"
	router.POST(loginPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.LoginUser(w, r)
	})

	loginReqBody, err := json.Marshal("invalid payload")
	a.NoError(err)

	// when ... we login the user
	loginReq := httptest.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)

	// then
	a.Equal(http.StatusBadRequest, loginRes.Code)
}

func (testSuite *UserHandlerTestSuite) TestLoginUser_MissingEmail() {
	t := testSuite.T()
	a := assert.New(t)

	// given
	router := httprouter.New()
	loginPath := "/users/login"
	router.POST(loginPath, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		testSuite.userHandler.LoginUser(w, r)
	})

	loginReqBody, err := json.Marshal(map[string]string{
		"email":    "",
		"password": "password123",
	})
	a.NoError(err)

	// when ... we login the user
	loginReq := httptest.NewRequest(http.MethodPost, loginPath, bytes.NewBuffer(loginReqBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)

	// then
	a.Equal(http.StatusBadRequest, loginRes.Code)
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
