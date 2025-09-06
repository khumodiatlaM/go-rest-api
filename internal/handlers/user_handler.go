package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-rest-api/internal/core"
	"net/http"
	"strings"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *core.User) (*core.User, error)
	GetUserByID(ctx context.Context, id string) (*core.User, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
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
		return errors.New("username is required")
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		return errors.New("valid email is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}
	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var userReq CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid create user request payload", http.StatusBadRequest)
		return
	}

	if err := userReq.Validate(); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	user := userReq.ToUser()
	result, err := h.userService.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToUserResponse(*result))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	id := r.URL.Path[len("/users/"):]
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ToUserResponse(*user))
}
