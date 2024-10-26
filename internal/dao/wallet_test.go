package dao

import (
	"context"
	"testing"
	"wallet/pkg/e"

	"github.com/stretchr/testify/assert"
)

func TestFindOneByUserIDSuccess(t *testing.T) {
	var (
		id    int64 = 30001
		email       = "TestFindOneByUserIDSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	walletDao := NewWallet(dao)
	wallet, err := walletDao.FindOneByUserID(ctx, id)

	assert.NotNil(t, wallet)
	assert.NoError(t, err)
}

func TestFindOneByUserIDNotFound(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	walletDao := NewWallet(dao)
	wallet, err := walletDao.FindOneByUserID(ctx, 999)

	assert.Nil(t, wallet)
	assert.Equal(t, e.NotFound, err)
}

func TestFindOneByUserIDForUpdateSuccess(t *testing.T) {
	var (
		id    int64 = 30002
		email       = "TestFindOneByUserIDForUpdateSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	walletDao := NewWallet(dao)
	wallet, err := walletDao.FindOneByUserIDForUpdate(ctx, id)

	assert.NotNil(t, wallet)

	assert.NoError(t, err)
}

func TestUpdateBalanceByUserIDSuccess(t *testing.T) {
	var (
		id    int64 = 30003
		email       = "TestUpdateBalanceByUserIDSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	walletDao := NewWallet(dao)
	err := walletDao.UpdateBalanceByUserID(ctx, id, 150)
	assert.NoError(t, err)

	wallet, err := walletDao.FindOneByUserID(ctx, id)
	assert.Equal(t, float64(150), wallet.Balance)
}

func TestUpdateBalanceByUserIDNotFound(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	walletDao := NewWallet(dao)
	err := walletDao.UpdateBalanceByUserID(ctx, 999, 150.0)

	assert.Error(t, err)
}
