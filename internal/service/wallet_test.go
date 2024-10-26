package service

import (
	"testing"
	"wallet/internal/types"
	"wallet/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWalletServiceDepositSuccess(t *testing.T) {
	var (
		id    int64 = 200001
		email       = "TestWalletServiceDepositSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: id,
		Amount: 100.00,
	})

	assert.NoError(t, err)

	info, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, float64(100), info.Balance)
}

func TestWalletServiceDepositUserNotFound(t *testing.T) {
	svcCtx := setupSvcCtx(t)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: 99999,
		Amount: 100,
	})

	assert.Equal(t, e.NotFound, err)
}

func TestWalletServiceWithdrawSuccess(t *testing.T) {
	var (
		id    int64 = 200002
		email       = "TestWalletServiceWithdrawSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: id,
		Amount: 100,
	})
	assert.NoError(t, err)

	err = walletService.Withdraw(&types.WithdrawInput{
		UserID: id,
		Amount: 99,
	})
	assert.NoError(t, err)

	info, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, float64(1), info.Balance)
}

func TestWalletServiceWithdrawInsufficientFunds(t *testing.T) {
	var (
		id    int64 = 200003
		email       = "TestWalletServiceWithdrawInsufficientFunds@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Withdraw(&types.WithdrawInput{
		UserID: id,
		Amount: 111,
	})

	assert.Equal(t, e.InsufficientFunds, err)
}

func TestWalletServiceTransferSuccess(t *testing.T) {
	var (
		id     int64 = 200004
		id2    int64 = 200005
		email        = "TestWalletServiceTransferSuccess@email.com"
		email2       = "TestWalletServiceTransferSuccess2@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	addDBData(t, svcCtx, id2, email2)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: id,
		Amount: 100,
	})
	assert.NoError(t, err)

	err = walletService.Transfer(&types.TransferInput{
		UserID:       id,
		TargetUserID: id2,
		Amount:       1,
	})
	assert.NoError(t, err)

	info, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, float64(99), info.Balance)

	info2, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id2)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), info2.Balance)
}

func TestWalletServiceTransferFailure(t *testing.T) {
	var (
		id     int64 = 200006
		id2    int64 = 200007
		email        = "TestWalletServiceTransferFailure@email.com"
		email2       = "TestWalletServiceTransferFailure2@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	addDBData(t, svcCtx, id2, email2)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: id,
		Amount: 100,
	})
	assert.NoError(t, err)

	err = walletService.Deposit(&types.DepositInput{
		UserID: id2,
		Amount: 9999999999.99,
	})
	assert.NoError(t, err)

	err = walletService.Transfer(&types.TransferInput{
		UserID:       id,
		TargetUserID: id2,
		Amount:       100,
	})
	assert.NotNil(t, err)

	info, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, float64(100), info.Balance)

	info2, err := walletService.svcCtx.WalletDao.FindOneByUserID(ctx, id2)
	assert.NoError(t, err)
	assert.Equal(t, float64(9999999999.99), info2.Balance)
}

func TestWalletServiceBalanceSuccess(t *testing.T) {
	var (
		id    int64 = 200008
		email       = "TestWalletServiceBalanceSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)

	err := walletService.Deposit(&types.DepositInput{
		UserID: id,
		Amount: 100,
	})
	assert.NoError(t, err)

	output, err := walletService.Balance(&types.BalanceInput{UserID: id})
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, float64(100), output.Balance)
}

func TestWalletServiceRecordSuccess(t *testing.T) {
	var (
		id    int64 = 200009
		email       = "TestWalletServiceRecordSuccess@email.com"
	)

	svcCtx := setupSvcCtx(t)
	addDBData(t, svcCtx, id, email)
	ctx, _ := gin.CreateTestContext(nil)

	walletService := NewWalletService(ctx, svcCtx)
	output, _, err := walletService.Record(&types.RecordInput{
		UserID: id,
		Offset: 0,
		Limit:  10,
	})

	assert.NoError(t, err)
	assert.NotNil(t, output)
}
