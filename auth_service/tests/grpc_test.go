package tests

import (
	"context"
	"pos/auth_service/app/dto/request"
	"pos/auth_service/app/repositories"
	"pos/auth_service/app/services"
	"pos/auth_service/pb"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTokenGrpc(t *testing.T) {
	//First login
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

	loginToken := loginResponse.Token

	//Call ValidateToken via gRPC
	grpcClient := pb.NewAuthServiceClient(grpcConn)

	res, err := grpcClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: loginToken,
	})

	assert.NoError(t, err)
	assert.True(t, res.IsValid)
	t.Logf("gRPC ValidateToken Response: %+v\n", res)

	//Logout and test invalid token
	err = authService.Logout(loginToken)
	assert.NoError(t, err)
	res, err = grpcClient.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: loginToken,
	})
	assert.NoError(t, err)
	assert.False(t, res.IsValid)
	t.Logf("gRPC ValidateToken Response after logout: %+v\n", res)
}
