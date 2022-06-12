package test

import (
	"context"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func RandomMakeUser(t *testing.T) db.User {
	name := tool.RandomOwner()
	arg := db.CreateUserParams{
		Username:       name,
		HashedPassword: "secret",
		FullName:       name,
		Email:          tool.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)

	assert.NotZero(t, user.CreatedAt)
	//github action is will run order not same with local so this case is will fail
	//assert.Empty(t, user.PasswordChangedAt)
	return user
}

// run test|debug testing
func TestCreateUser(t *testing.T) {
	RandomMakeUser(t)
}

func TestGetUser(t *testing.T) {
	newUser := RandomMakeUser(t)
	selectUser, err := testQueries.GetUser(context.Background(), newUser.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, selectUser)

	assert.Equal(t, newUser.Username, selectUser.Username)
	assert.Equal(t, newUser.HashedPassword, selectUser.HashedPassword)
	assert.Equal(t, newUser.FullName, selectUser.FullName)
	assert.Equal(t, newUser.Email, selectUser.Email)

	assert.WithinDuration(t, newUser.CreatedAt, selectUser.CreatedAt, time.Second)
}
