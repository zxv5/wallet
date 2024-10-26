package dao

import (
	"context"
	"testing"
	"wallet/internal/consts"
	"wallet/internal/model"
	"wallet/pkg/e"

	"github.com/stretchr/testify/assert"
)

var describe = "describe"

func TestFindListAndCountByWalletIDSuccess(t *testing.T) {
	var (
		id    int64 = 20001
		email       = "TestFindListAndCountByWalletIDSuccess@email.com"
	)

	dao := setupTestDB(t)
	addDBData(t, dao, id, email)

	ctx := context.Background()
	walletRecordDao := NewWalletRecord(dao)
	records, _, err := walletRecordDao.FindListAndCountByWalletID(ctx, id, 0, 10)

	assert.NotNil(t, records)
	assert.NoError(t, err)
}

func TestFindListAndCountByWalletIDNotFound(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	walletRecordDao := NewWalletRecord(dao)
	records, totalCount, _ := walletRecordDao.FindListAndCountByWalletID(ctx, 999, 0, 10)

	assert.Nil(t, records)
	assert.Equal(t, int64(0), totalCount)
}

func TestCreateSuccess(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	walletRecordDao := NewWalletRecord(dao)
	walletRecord := &model.WalletRecord{
		WalletID:        1,
		Amount:          1,
		TransactionType: consts.WalletRecordTypeIncome,
		Describe:        &describe,
	}
	err := walletRecordDao.Create(ctx, walletRecord)

	assert.NoError(t, err)
}

func TestCreateError(t *testing.T) {
	dao := setupTestDB(t)

	ctx := context.Background()
	walletRecordDao := NewWalletRecord(dao)
	walletRecord := &model.WalletRecord{
		WalletID:        1,
		Amount:          99999999999999999,
		TransactionType: consts.WalletRecordTypeIncome,
		Describe:        &describe,
	}
	err := walletRecordDao.Create(ctx, walletRecord)

	assert.Equal(t, e.SQLErr, err)

}
