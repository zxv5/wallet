package dao

import (
	"context"
	"wallet/internal/consts"
	"wallet/internal/model"
	"wallet/pkg/e"
)

type WalletRecord struct {
	dao *Dao
}

func NewWalletRecord(dao *Dao) *WalletRecord {
	return &WalletRecord{dao: dao}
}

func (d *WalletRecord) FindListAndCountByWalletID(ctx context.Context, walletID, offset, limit int64) ([]*model.WalletRecord, int64, e.Codes) {
	query := `SELECT "id", "wallet_id", "amount", "transaction_type", "describe", "created_at" FROM "wallet_record"
	WHERE "wallet_id" = $1 and "deleted" = $2 ORDER BY "created_at" DESC OFFSET $3 LIMIT $4;`

	rows, err := d.dao.WithContext(ctx).Query(query, walletID, consts.NotDeleted, offset, limit)
	if err != nil {
		return nil, 0, e.SQLErr
	}
	defer rows.Close()

	var walletRecords []*model.WalletRecord
	for rows.Next() {
		var walletRecord model.WalletRecord
		if err := rows.Scan(
			&walletRecord.ID,
			&walletRecord.WalletID,
			&walletRecord.Amount,
			&walletRecord.TransactionType,
			&walletRecord.Describe,
			&walletRecord.CreatedAt,
		); err != nil {
			return nil, 0, e.SQLErr
		}
		walletRecords = append(walletRecords, &walletRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, e.SQLErr
	}

	var totalCount int64
	countQuery := `SELECT count(*) as count FROM "wallet_record" WHERE "wallet_id" = $1 AND "deleted" = $2`

	err = d.dao.WithContext(ctx).QueryRow(countQuery, walletID, consts.NotDeleted).Scan(&totalCount)
	if err != nil {
		return nil, 0, e.SQLErr
	}

	return walletRecords, totalCount, nil
}

func (d *WalletRecord) Create(ctx context.Context, walletRecord *model.WalletRecord) e.Codes {
	query := `INSERT INTO "wallet_record" ("wallet_id", "amount", "transaction_type", "describe") VALUES ($1, $2, $3, $4)`

	_, err := d.dao.WithContext(ctx).Exec(query, walletRecord.WalletID, walletRecord.Amount, walletRecord.TransactionType, walletRecord.Describe)
	if err != nil {
		return e.SQLErr
	}

	return nil
}
