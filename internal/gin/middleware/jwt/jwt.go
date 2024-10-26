package jwt

import (
	"strconv"
	"strings"
	"time"
	"wallet/internal/config"
	"wallet/internal/gin/core"
	"wallet/internal/types"
	"wallet/pkg/jwt"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// UserInfoKey User information key
const UserInfoKey string = "__userinfo"

// Jwt gin middleware
func Jwt(c *config.Jwt) gin.HandlerFunc {
	config := jwt.Config{Secret: c.Secret}

	return func(c *gin.Context) {
		ctx := core.New(c)

		authorization := c.GetHeader("Authorization")
		if len(authorization) == 0 {
			ctx.SendNotLogin()
			return
		}

		arr := strings.Split(authorization, " ")
		if len(arr) != 2 {
			ctx.SendNotLogin()
			return
		}

		tokenS := arr[1]
		wrap := jwt.New(&types.UserJwtClaims{}, config)

		token, err := wrap.Parse(tokenS)
		if err != nil {
			ctx.SendNotLogin()
			return
		}

		if cl, ok := token.Claims.(*types.UserJwtClaims); ok {
			c.Set(UserInfoKey, cl.UserInfo)
			c.Next()
			return
		}

		ctx.SendNotLogin()
	}
}

// Sign .
func Sign(c *config.Jwt, userInfo *types.UserInfo) (string, error) {
	config := jwt.Config{Secret: c.Secret}

	claims := &types.UserJwtClaims{
		RegisteredClaims: jwtv5.RegisteredClaims{
			Subject:   strconv.Itoa(int(userInfo.ID)),
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Hour * time.Duration(c.Exp))),
		},
		UserInfo: *userInfo,
	}

	wrap := jwt.New(claims, config)

	str, err := wrap.Sign()
	if err != nil {
		return "", err
	}
	return str, nil
}
