package consts

const (
	NotDeleted int64 = iota // Not deleted
	IsDeleted               // Deleted
)

const (
	UserStatusActive   = iota // active
	UserStatusInactive        // inactive
)

const (
	WalletRecordTypeUnknown = iota // unknown
	WalletRecordTypeIncome         // income
	WalletRecordTypeExpend         // expend
)
