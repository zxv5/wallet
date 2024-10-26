package service

import (
	"wallet/internal/config"
	"wallet/internal/dao"
)

type ServiceContext struct {
	Config          *config.Config
	Dao             *dao.Dao
	UserDao         *dao.User
	WalletDao       *dao.Wallet
	WalletRecordDao *dao.WalletRecord
}

func NewServiceContext(c *config.Config) *ServiceContext {
	d := dao.NewDao(c)
	return &ServiceContext{
		Config:          c,
		Dao:             d,
		UserDao:         dao.NewUser(d),
		WalletDao:       dao.NewWallet(d),
		WalletRecordDao: dao.NewWalletRecord(d),
	}
}
