package jwt

import (
	"crypto/md5"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Config conf
type Config struct {
	Secret string `json:"secret"`
	Exp    int    `json:"exp"` // hours
}

// ClaimsWrap .
type ClaimsWrap struct {
	jwt.Claims
	Config
}

// ClaimsWraper .
type ClaimsWraper interface {
	GetSignKey() []byte
	Sign() (string, error)
	Parse(tokenString string) (*jwt.Token, error)
}

// New ClaimsWraper
func New(claims jwt.Claims, config Config) ClaimsWraper {
	return &ClaimsWrap{claims, config}
}

// GetSignKey .
func (c *ClaimsWrap) GetSignKey() []byte {
	s := fmt.Sprintf("%x", md5.Sum([]byte(c.Secret)))
	return []byte(s)
}

// Sign .
func (c *ClaimsWrap) Sign() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c.Claims)
	s, err := token.SignedString(c.GetSignKey())
	return s, err
}

// Parse .
func (c *ClaimsWrap) Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, c.Claims, func(token *jwt.Token) (interface{}, error) {
		return c.GetSignKey(), nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return nil, err
}
