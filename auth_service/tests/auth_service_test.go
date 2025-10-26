package tests

import (
	"pos/auth_service/app/dto/request"
	"pos/auth_service/app/repositories"
	"pos/auth_service/app/services"
	"pos/auth_service/app/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Log("This is a placeholder test for login service.")

	authService := services.NewAuthService(repositories.NewUserRepository(testDB), testRedis)

	loginRequest := request.LoginRequest{
		Email:    "superadmin@warungacehbangari.com",
		Password: "Password",
	}

	loginResponse, err := authService.Login(loginRequest)
	assert.NoError(t, err)
	assert.Equal(t, loginRequest.Email, loginResponse.User.Email)
	t.Logf("Login Response: %+v\n", loginResponse)

	// Validate the token
	isValid, err := authService.ValidateToken(loginResponse.Token)
	assert.NoError(t, err)
	assert.True(t, isValid)

	// Test with invalid credentials
	invalidLoginRequest := request.LoginRequest{
		Email:    "superadmin@warungacehbangari.com",
		Password: "WrongPassword",
	}

	loginResponse, err = authService.Login(invalidLoginRequest)
	assert.Error(t, err)
	assert.Empty(t, loginResponse.Token)
	t.Logf("Login Response: %+v\n", loginResponse)

}

func TestHashPassword(t *testing.T) {
	password := "Password"
	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)
	t.Logf("Hashed Password: %s\n", hashedPassword)

	// Verify the password
	isValid := utils.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid)
}

func TestValidateToken(t *testing.T) {
	t.Log("This is a placeholder test for token validation.")
	authService := services.NewAuthService(repositories.NewUserRepository(testDB), testRedis)

	loginRequest := request.LoginRequest{
		Email:    "superadmin@warungacehbangari.com",
		Password: "Password",
	}

	loginResponse, err := authService.Login(loginRequest)
	assert.NoError(t, err)
	isValid, err := authService.ValidateToken(loginResponse.Token)
	assert.NoError(t, err)
	assert.True(t, isValid)

	// Test with an invalid token
	invalidToken := "invalid.token.string"
	isValid, err = authService.ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestLogout(t *testing.T) {
	t.Log("This is a placeholder test for logout service.")
	authService := services.NewAuthService(repositories.NewUserRepository(testDB), testRedis)

	loginRequest := request.LoginRequest{
		Email:    "superadmin@warungacehbangari.com",
		Password: "Password",
	}

	loginResponse, err := authService.Login(loginRequest)
	assert.NoError(t, err)
	// Validate token before logout
	isValid, err := authService.ValidateToken(loginResponse.Token)
	assert.NoError(t, err)
	assert.True(t, isValid)

	err = authService.Logout(loginResponse.Token)
	assert.NoError(t, err)
	t.Log("Logout successful.")

	// Validate token after logout
	isValid, err = authService.ValidateToken(loginResponse.Token)
	assert.NoError(t, err)
	assert.False(t, isValid)
}
