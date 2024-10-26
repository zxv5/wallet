package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserJwtClaims struct {
	jwt.RegisteredClaims
	UserInfo
}

type UserInfo struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    int64     `json:"gender"`
	Status    int64     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UserLoginReq struct {
	Email    string `json:"email"`    // login email
	Password string `json:"password"` // login password
}

type UserLoginOutput struct {
	Token string    `json:"token"`
	Info  *UserInfo `json:"info"`
}
