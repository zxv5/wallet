package controller

import (
	"wallet/internal/gin/core"
	"wallet/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svcCtx *service.ServiceContext
}

func NewUserController(svcCtx *service.ServiceContext) *UserController {
	return &UserController{svcCtx: svcCtx}
}

// Login User Login
// @Summary User Login
// @Description User Login
// @Tags User
// @Accept application/json
// @Produce application/json
// @Param body body types.UserLoginReq true "parameter"
// @Security ApiKeyAuth
// @Success 200 {object} types.UserLoginOutput
// @Router /user/login [POST]
func (ctrl *UserController) Login(c *gin.Context) {
	var form struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required"`
	}
	ctx := core.New(c).Bind(&form)

	userService := service.NewUserService(c, ctrl.svcCtx)
	output, err := userService.Login(form.Email, form.Password)
	if err != nil {
		ctx.SendErr(err)
		return
	}

	ctx.SendOk(output)
}
