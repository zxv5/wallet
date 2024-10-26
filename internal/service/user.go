package service

import (
	"wallet/internal/consts"
	"wallet/internal/gin/middleware/jwt"
	"wallet/internal/types"
	"wallet/pkg/e"

	"github.com/gin-gonic/gin"
)

type userService struct {
	ctx    *gin.Context
	svcCtx *ServiceContext
}

func NewUserService(ctx *gin.Context, svcCtx *ServiceContext) *userService {
	return &userService{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (svc *userService) Login(email, password string) (*types.UserLoginOutput, e.Codes) {
	userInfo, err := svc.svcCtx.UserDao.FindOneByEmail(svc.ctx, email)
	if err != nil {
		return nil, err
	}

	if userInfo.Status == consts.UserStatusInactive {
		return nil, e.UserNotActive
	}

	if userInfo.Password != password {
		return nil, e.UsernameOrPasswordErr
	}

	outputInfo := &types.UserInfo{
		ID:        userInfo.ID,
		Email:     userInfo.Email,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
		Gender:    userInfo.Gender,
		Status:    userInfo.Status,
		CreatedAt: userInfo.CreatedAt,
	}

	token, err2 := jwt.Sign(&svc.svcCtx.Config.Jwt, outputInfo)
	if err2 != nil {
		return nil, e.Cause(err2)
	}

	output := &types.UserLoginOutput{
		Token: token,
		Info:  outputInfo,
	}

	return output, nil
}
