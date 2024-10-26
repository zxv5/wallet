package dao

import (
	"context"
	"testing"
	"wallet/pkg/e"

	"github.com/stretchr/testify/assert"
)

func TestFindOneByIDSuccess(t *testing.T) {
	var (
		id    int64 = 10001
		email       = "TestFindOneByIDSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	userDao := NewUser(dao)
	user, err := userDao.FindOneByID(ctx, id)

	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.NoError(t, err)

}

func TestFindOneByIDNotFound(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	userDao := NewUser(dao)
	user, code := userDao.FindOneByID(ctx, 999)

	assert.Nil(t, user)
	assert.Equal(t, e.NotFound, code)
}

func TestFindOneByEmailSuccess(t *testing.T) {
	var (
		id    int64 = 10002
		email       = "TestFindOneByEmailSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	userDao := NewUser(dao)
	user, code := userDao.FindOneByEmail(ctx, email)

	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.NoError(t, code)
}

func TestFindOneByEmailNotFound(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	userDao := NewUser(dao)
	user, code := userDao.FindOneByEmail(ctx, "notfound@email.com")

	assert.Nil(t, user)
	assert.Equal(t, e.NotFound, code)

}
