package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	a := assert.New(t)
	password := "password"
	hashedPassword, err := HashPassword(password)

	a.NoError(err)
	a.NotEmpty(hashedPassword)
	a.NotEqual(password, hashedPassword)
}

func TestVerifyPassword(t *testing.T) {
	a := assert.New(t)
	password := "password"
	hashedPassword, err := HashPassword(password)

	a.NoError(err)
	a.NoError(VerifyPassword(hashedPassword, password))
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	a := assert.New(t)
	password := "password"
	hashedPassword, err := HashPassword(password)

	a.NoError(err)
	a.Error(VerifyPassword(hashedPassword, "wrongpassword"))
}
