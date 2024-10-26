package model

import (
	"time"
)

const TableNameWalletRecord = "wallet_record"

// WalletRecord mapped from table <wallet_record>
type WalletRecord struct {
	ID              int64     `json:"id"`
	WalletID        int64     `json:"wallet_id"`
	Amount          float64   `json:"amount"`
	TransactionType int64     `json:"transaction_type"` // 0=unknown 1=income 2=expend
	Describe        *string   `json:"describe"`
	Deleted         int64     `json:"deleted"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// TableName WalletRecord's table name
func (*WalletRecord) TableName() string {
	return TableNameWalletRecord
}
