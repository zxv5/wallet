package service

import (
	"context"
	"fmt"
	"wallet/internal/consts"
	"wallet/internal/model"
	"wallet/internal/types"
	"wallet/pkg/e"

	"github.com/gin-gonic/gin"
)

type walletService struct {
	ctx    *gin.Context
	svcCtx *ServiceContext
}

func NewWalletService(ctx *gin.Context, svcCtx *ServiceContext) *walletService {
	return &walletService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (svc *walletService) Deposit(in *types.DepositInput) e.Codes {
	if err := svc.svcCtx.Dao.Transaction(svc.ctx, func(ctx context.Context) e.Codes {
		walletInfo, err := svc.svcCtx.WalletDao.FindOneByUserIDForUpdate(ctx, in.UserID)
		if err != nil {
			return err
		}

		newBalance := walletInfo.Balance + in.Amount
		if err := svc.svcCtx.WalletDao.UpdateBalanceByUserID(ctx, in.UserID, newBalance); err != nil {
			return err
		}

		describe := fmt.Sprintf("Deposit %.2f", in.Amount)
		walletRecord := &model.WalletRecord{
			WalletID:        walletInfo.ID,
			Amount:          in.Amount,
			TransactionType: consts.WalletRecordTypeIncome,
			Describe:        &describe,
		}
		// Transaction records
		if err := svc.svcCtx.WalletRecordDao.Create(ctx, walletRecord); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (svc *walletService) Withdraw(in *types.WithdrawInput) e.Codes {
	if err := svc.svcCtx.Dao.Transaction(svc.ctx, func(ctx context.Context) e.Codes {
		walletInfo, err := svc.svcCtx.WalletDao.FindOneByUserIDForUpdate(ctx, in.UserID)
		if err != nil {
			return err
		}

		newBalance := walletInfo.Balance - in.Amount

		// Determine the balance is sufficient
		if newBalance < 0 {
			return e.InsufficientFunds
		}
		if err := svc.svcCtx.WalletDao.UpdateBalanceByUserID(ctx, in.UserID, newBalance); err != nil {
			return err
		}

		describe := fmt.Sprintf("Withdraw %.2f", in.Amount)
		walletRecord := &model.WalletRecord{
			WalletID:        walletInfo.ID,
			Amount:          in.Amount,
			TransactionType: consts.WalletRecordTypeExpend,
			Describe:        &describe,
		}
		// Transaction records
		if err := svc.svcCtx.WalletRecordDao.Create(ctx, walletRecord); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (svc *walletService) Transfer(in *types.TransferInput) e.Codes {
	if in.UserID == in.TargetUserID {
		return e.UserDuplicate
	}

	targetUserInfo, err := svc.svcCtx.UserDao.FindOneByID(svc.ctx, in.TargetUserID)
	if err != nil {
		return e.UserNotExist
	}

	if err := svc.svcCtx.Dao.Transaction(svc.ctx, func(ctx context.Context) e.Codes {
		// Get the wallet information for both users and wait for update
		walletInfo, err := svc.svcCtx.WalletDao.FindOneByUserIDForUpdate(ctx, in.UserID)
		if err != nil {
			return err
		}

		targetWalletInfo, err := svc.svcCtx.WalletDao.FindOneByUserIDForUpdate(ctx, in.TargetUserID)
		if err != nil {
			return err
		}

		newBalance := walletInfo.Balance - in.Amount
		// Determine the balance is sufficient
		if newBalance < 0 {
			return e.InsufficientFunds
		}
		// Update the user's balance
		if err := svc.svcCtx.WalletDao.UpdateBalanceByUserID(ctx, in.UserID, newBalance); err != nil {
			return err
		}

		targetNewBalance := targetWalletInfo.Balance + in.Amount
		// Update the target user's balance
		if err := svc.svcCtx.WalletDao.UpdateBalanceByUserID(ctx, in.TargetUserID, targetNewBalance); err != nil {
			return err
		}

		describe := fmt.Sprintf("Transfer %.2f to the user [%s %s]", in.Amount, targetUserInfo.FirstName, targetUserInfo.LastName)
		walletRecord := &model.WalletRecord{
			WalletID:        walletInfo.ID,
			Amount:          in.Amount,
			TransactionType: consts.WalletRecordTypeExpend,
			Describe:        &describe,
		}
		// Transaction records
		if err := svc.svcCtx.WalletRecordDao.Create(ctx, walletRecord); err != nil {
			return err
		}

		targetDescribe := fmt.Sprintf("Received a transfer of %.2f from [%s %s]", in.Amount, in.FirstName, in.LastName)
		targetWalletRecord := &model.WalletRecord{
			WalletID:        targetWalletInfo.ID,
			Amount:          in.Amount,
			TransactionType: consts.WalletRecordTypeIncome,
			Describe:        &targetDescribe,
		}
		// Transaction records
		if err := svc.svcCtx.WalletRecordDao.Create(ctx, targetWalletRecord); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (svc *walletService) Balance(in *types.BalanceInput) (*types.BalanceOutput, e.Codes) {
	walletInfo, err := svc.svcCtx.WalletDao.FindOneByUserID(svc.ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	output := &types.BalanceOutput{
		UserID:  in.UserID,
		Balance: walletInfo.Balance,
	}

	return output, nil
}

func (svc *walletService) Record(in *types.RecordInput) ([]*types.RecordOutput, int64, e.Codes) {
	walletInfo, err := svc.svcCtx.WalletDao.FindOneByUserID(svc.ctx, in.UserID)
	if err != nil {
		return nil, 0, err
	}

	walletRecordList, total, err := svc.svcCtx.WalletRecordDao.FindListAndCountByWalletID(svc.ctx, walletInfo.ID, in.Offset, in.Limit)
	if err != nil {
		return nil, 0, err
	}

	var output []*types.RecordOutput
	for _, record := range walletRecordList {
		output = append(output, &types.RecordOutput{
			ID:              record.ID,
			TransactionType: record.TransactionType,
			Amount:          record.Amount,
			Describe:        record.Describe,
			CreatedAt:       record.CreatedAt,
		})
	}

	return output, total, nil
}
