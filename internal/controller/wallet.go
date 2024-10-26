package controller

import (
	"wallet/internal/gin/core"
	"wallet/internal/service"
	"wallet/internal/types"
	"wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

type walletController struct {
	svcCtx *service.ServiceContext
}

func NewWalletController(svcCtx *service.ServiceContext) *walletController {
	return &walletController{svcCtx: svcCtx}
}

// Deposit User deposits
// @Summary User deposits
// @Description User deposits
// @Tags Wallet
// @Accept application/json
// @Produce application/json
// @Param body body types.DepositReq true "parameter"
// @Security ApiKeyAuth
// @Success 200
// @Router /wallet/deposit [POST]
func (ctrl *walletController) Deposit(c *gin.Context) {
	var form struct {
		Amount float64 `form:"amount" json:"amount" binding:"required,gt=0,max=9999999999"`
	}
	var userInfo types.UserInfo
	ctx := core.New(c).Bind(&form).BindUserInfo(&userInfo)

	in := &types.DepositInput{
		UserID: userInfo.ID,
		Amount: form.Amount,
	}
	walletService := service.NewWalletService(c, ctrl.svcCtx)
	if err := walletService.Deposit(in); err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.SendOk()
}

// Withdraw User withdraw
// @Summary User withdraw
// @Description User withdraw
// @Tags Wallet
// @Accept application/json
// @Produce application/json
// @Param body body types.WithdrawReq true "parameter"
// @Security ApiKeyAuth
// @Success 200
// @Router /wallet/withdraw [POST]
func (ctrl *walletController) Withdraw(c *gin.Context) {
	var form struct {
		Amount float64 `form:"amount" json:"amount" binding:"required,gt=0,max=9999999999"`
	}
	var userInfo types.UserInfo
	ctx := core.New(c).Bind(&form).BindUserInfo(&userInfo)

	in := &types.WithdrawInput{
		UserID: userInfo.ID,
		Amount: form.Amount,
	}
	walletService := service.NewWalletService(c, ctrl.svcCtx)
	if err := walletService.Withdraw(in); err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.SendOk()
}

// Transfer User Transfer
// @Summary User Transfer
// @Description User Transfer
// @Tags Wallet
// @Accept application/json
// @Produce application/json
// @Param body body types.TransferReq true "parameter"
// @Security ApiKeyAuth
// @Success 200
// @Router /wallet/transfer [POST]
func (ctrl *walletController) Transfer(c *gin.Context) {
	var form struct {
		Amount       float64 `form:"amount" json:"amount" binding:"required,gt=0,max=9999999999"`
		TargetUserID int64   `form:"target_user_id" json:"target_user_id" binding:"required"`
	}
	var userInfo types.UserInfo
	ctx := core.New(c).Bind(&form).BindUserInfo(&userInfo)

	in := &types.TransferInput{
		UserID:       userInfo.ID,
		FirstName:    userInfo.FirstName,
		LastName:     userInfo.LastName,
		TargetUserID: form.TargetUserID,
		Amount:       form.Amount,
	}
	walletService := service.NewWalletService(c, ctrl.svcCtx)
	if err := walletService.Transfer(in); err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.SendOk()
}

// Balance Get the user's balance
// @Summary Get the user's balance
// @Description Get the user's balance
// @Tags Wallet
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} types.BalanceOutput
// @Router /wallet/balance [GET]
func (ctrl *walletController) Balance(c *gin.Context) {
	var userInfo types.UserInfo
	ctx := core.New(c).BindUserInfo(&userInfo)

	in := &types.BalanceInput{
		UserID: userInfo.ID,
	}
	walletService := service.NewWalletService(c, ctrl.svcCtx)
	output, err := walletService.Balance(in)
	if err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.SendOk(output)
}

// Record Get user transaction records
// @Summary Get user transaction records
// @Description Get user transaction records
// @Tags Wallet
// @Accept application/json
// @Produce application/json
// @Param query query types.RecordReq true "parameter"
// @Security ApiKeyAuth
// @Success 200 {array} types.RecordOutput
// @Router /wallet/record [GET]
func (ctrl *walletController) Record(c *gin.Context) {
	var query struct {
		Page int64 `form:"page" binding:"omitempty"`
		Size int64 `form:"size" binding:"omitempty"`
	}
	var userInfo types.UserInfo
	ctx := core.New(c).BindQuery(&query).BindUserInfo(&userInfo)

	offset, limit := utils.GetOffsetLimit(query.Page, query.Size)

	in := &types.RecordInput{
		UserID: userInfo.ID,
		Offset: offset,
		Limit:  limit,
	}
	walletService := service.NewWalletService(c, ctrl.svcCtx)
	output, total, err := walletService.Record(in)
	if err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.Append("total", total).SendOk(output)
}
