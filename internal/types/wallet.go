package types

import "time"

type DepositReq struct {
	Amount float64 `json:"amount"` // Amount to deposit
}

type DepositInput struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

type WithdrawReq struct {
	Amount float64 `json:"amount"` // Amount to withdraw
}

type WithdrawInput struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransferReq struct {
	TargetUserID int64   `json:"target_user_id"` // Transfer target user ID
	Amount       float64 `json:"amount"`         // Amount to transfer
}

type TransferInput struct {
	UserID       int64   `json:"user_id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	TargetUserID int64   `json:"target_user_id"`
	Amount       float64 `json:"amount"`
}

type BalanceInput struct {
	UserID int64 `json:"user_id"`
}

type BalanceOutput struct {
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}

type RecordReq struct {
	Page int64 `json:"page"` // page
	Size int64 `json:"size"` // size
}

type RecordInput struct {
	UserID int64 `json:"user_id"`
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type RecordOutput struct {
	ID              int64     `json:"id"`               // wallet record id
	TransactionType int64     `json:"transaction_type"` // transaction type 0=unknown 1=income 2=expend
	Amount          float64   `json:"amount"`           // amount
	Describe        *string   `json:"describe"`         // description
	CreatedAt       time.Time `json:"created_at"`       // Creation time
}
