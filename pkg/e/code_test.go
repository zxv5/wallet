package e

import (
	"errors"
	"testing"
)

func TestRegisterAndMessage(t *testing.T) {
	code := NewCode(1001)
	message := "Test error message"

	Register("en-us", code, message)

	retrievedMessage := code.Message()
	if retrievedMessage != message {
		t.Errorf("Expected message '%s', got '%s'", message, retrievedMessage)
	}
}

func TestSetDefault(t *testing.T) {
	code := NewCode(1002)
	message := "Another test error message"

	Register("en-us", code, message)
	SetDefault("fr-fr") // Set to a non-registered language

	retrievedMessage := code.Message()
	if retrievedMessage != "1002" { // Fallback since "fr-fr" is not registered
		t.Errorf("Expected message '1002', got '%s'", retrievedMessage)
	}

	// Register the same code for the default language
	Register("en-us", code, "Test message in English")
	SetDefault("en-us")

	retrievedMessage = code.Message()
	if retrievedMessage != "Test message in English" {
		t.Errorf("Expected message 'Test message in English', got '%s'", retrievedMessage)
	}
}

func TestNewCode(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when creating code with value <= 0, but did not panic")
		}
	}()

	NewCode(0) // This should panic
}

func TestCause(t *testing.T) {
	err := errors.New("simple error")
	code := String("4001")
	cause := Cause(err)

	// Testing with Codes type
	customErr := New(code, err)
	cause = Cause(customErr)

	if cause.Code() != code.Code() {
		t.Errorf("Expected cause code %d, got %d", code.Code(), cause.Code())
	}
}
