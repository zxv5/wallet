package dao

import (
	"context"
	"database/sql"
	"wallet/internal/model"
	"wallet/pkg/e"
	"wallet/pkg/logger"
)

type Wallet struct {
	dao *Dao
}

func NewWallet(dao *Dao) *Wallet {
	return &Wallet{dao: dao}
}

func (d *Wallet) FindOneByUserID(ctx context.Context, userID int64) (*model.Wallet, e.Codes) {
	var wallet model.Wallet

	query := `SELECT "id", "balance" FROM "wallet" WHERE "user_id" = $1;`
	row := d.dao.WithContext(ctx).QueryRow(query, userID)

	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NotFound
		}
		logger.Errorf("failed to find wallet: %w", err)
		return nil, e.SQLErr
	}

	return &wallet, nil
}

// FindOneByUserIDForUpdate finds a wallet by user ID and locks it for update.
func (d *Wallet) FindOneByUserIDForUpdate(ctx context.Context, userID int64) (*model.Wallet, e.Codes) {
	var wallet model.Wallet

	query := `SELECT "id", "balance" FROM "wallet" WHERE "user_id" = $1 FOR UPDATE;`
	row := d.dao.WithContext(ctx).QueryRow(query, userID)

	err := row.Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NotFound
		}
		logger.Errorf("failed to find wallet: %w", err)
		return nil, e.SQLErr
	}

	return &wallet, nil
}

// UpdateBalanceByUserID Update wallet balance by user
func (d *Wallet) UpdateBalanceByUserID(ctx context.Context, userID int64, newBalance float64) e.Codes {
	query := `UPDATE "wallet" SET "balance" = $1 WHERE "user_id" = $2;`
	result, err := d.dao.WithContext(ctx).Exec(query, newBalance, userID)

	if err != nil {
		logger.Errorf("failed to update wallet: %w", err)
		return e.SQLErr
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("failed to get rows affected: %w", err)
		return e.SQLErr
	}
	if rowsAffected == 0 {
		logger.Errorf("no wallet found with userID %d", userID)
		return e.SQLErr
	}

	return nil
}
