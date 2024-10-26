package model

import (
	"time"
)

const TableNameWallet = "wallet"

// Wallet mapped from table <wallet>
type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	Deleted   int64     `json:"deleted"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName Wallet's table name
func (*Wallet) TableName() string {
	return TableNameWallet
}
