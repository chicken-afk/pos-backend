package tests

import (
	"pos/auth_service/app/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	email := "superadmin@warungacehbangari.com"
	userRepo := repositories.NewUserRepository(testDB)

	user, err := userRepo.FindByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.NotEqual(t, user.Role.ID, 0)
	//print user details
	t.Logf("User: %+v\n", user)
	t.Logf("Role: %+v\n", user.Role)
	t.Logf("Outlet: %+v\n", user.Outlet)
}
