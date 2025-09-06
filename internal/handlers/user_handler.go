package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-rest-api/internal/core"
	"go-rest-api/pkg/logger"
	"net/http"
	"strings"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *core.User) (*core.User, error)
	GetUserByID(ctx context.Context, id string) (*core.User, error)
	LoginUser(ctx context.Context, email, password, jwtSecret string) (string, error)
}

type UserHandler struct {
	userService UserService
	Logger      logger.CustomLogger
	JwtSecret   string
}

func NewUserHandler(userService UserService, logger logger.CustomLogger, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: userService,
		Logger:      logger,
		JwtSecret:   jwtSecret,
	}
}

func HeathCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"available"}`))
}

func (u *CreateUserRequest) ToUser() core.User {
	return core.User{
		Username: u.Username,
		Email:    strings.ToLower(u.Email),
		Password: u.Password,
	}
}

func ToUserResponse(u core.User) UserResponse {
	return UserResponse{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return errors.New("Username is required")
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return errors.New("Valid email is required")
	}
	if len(req.Password) < 6 {
		return errors.New("Password must be at least 6 characters long")
	}
	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var userReq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		writeJSONErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := userReq.Validate(); err != nil {
		writeJSONErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user := userReq.ToUser()
	result, err := h.userService.CreateUser(ctx, &user)
	if err != nil {
		writeJSONErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToUserResponse(*result))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	id := r.URL.Path[len("/users/"):]
	if id == "" {
		writeJSONErrorResponse(w, http.StatusBadRequest, "User Id is required")
		return
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		writeJSONErrorResponse(w, http.StatusInternalServerError, "User not found")
		return
	}
	if user == nil {
		writeJSONErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToUserResponse(*user))
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var userReq LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		writeJSONErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// ... validate request data
	if userReq.Email == "" || userReq.Password == "" {
		writeJSONErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	token, err := h.userService.LoginUser(r.Context(), strings.ToLower(userReq.Email), userReq.Password, h.JwtSecret)
	if err != nil {
		writeJSONErrorResponse(w, http.StatusInternalServerError, "Failed to login user")
		return
	}

	if token == "" {
		writeJSONErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginUserResponse{token})
}

func writeJSONErrorResponse(w http.ResponseWriter, status int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{errorMsg})
}
