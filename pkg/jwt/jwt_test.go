package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignAndParse(t *testing.T) {
	// Setup a configuration with a secret and expiration time
	config := Config{
		Secret: "mysecretkey",
		Exp:    1, // 1 hour
	}

	// Create claims
	claims := jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	// Create a new ClaimsWrap
	claimsWrap := New(claims, config)

	// Sign the token
	tokenString, err := claimsWrap.Sign()
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Parse the token
	parsedToken, err := claimsWrap.Parse(tokenString)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	// Check if the token is valid
	if !parsedToken.Valid {
		t.Fatal("Parsed token is not valid")
	}

	// Check claims
	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to cast claims")
	}

	if parsedClaims["foo"] != "bar" {
		t.Errorf("Expected claim foo to be bar, got %v", parsedClaims["foo"])
	}

	// Check expiration claim
	exp, ok := parsedClaims["exp"].(float64)
	if !ok {
		t.Fatal("Expected exp claim to be present")
	}

	if int64(exp) < time.Now().Unix() {
		t.Error("Token is expired")
	}
}

func TestInvalidToken(t *testing.T) {
	config := Config{Secret: "mysecretkey"}
	claims := jwt.MapClaims{}

	claimsWrap := New(claims, config)

	// Try parsing an invalid token
	_, err := claimsWrap.Parse("invalid.token")
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}
}

func TestExpiredToken(t *testing.T) {
	config := Config{
		Secret: "mysecretkey",
		Exp:    1, // 1 hour
	}

	// Create claims with an expiration in the past
	claims := jwt.MapClaims{
		"foo": "bar",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}

	claimsWrap := New(claims, config)

	// Sign the token
	tokenString, err := claimsWrap.Sign()
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	// Parse the token
	parsedToken, err := claimsWrap.Parse(tokenString)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}

	if parsedToken != nil {
		t.Fatalf("Expected nil token for expired token, got %v", parsedToken)
	}
}
