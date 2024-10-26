package e

import (
	"errors"
	"testing"
)

func TestNewE(t *testing.T) {
	code := NewCode(9999991)
	err := errors.New("test error")

	customErr := New(code, err)

	if customErr == nil {
		t.Fatal("Expected a valid error, got nil")
	}

	if customErr.Code() != code.Code() {
		t.Errorf("Expected code %d, got %d", code.Code(), customErr.Code())
	}

	expectedMessage := "test error"
	if customErr.Message() != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, customErr.Message())
	}

	expectedErrorString := "Error: test error, Code: 9999991, Message: test error"
	if customErr.Error() != expectedErrorString {
		t.Errorf("Unexpected error string: %s", customErr.Error())
	}
}

func TestNewI(t *testing.T) {
	code := 9999992
	err := errors.New("another test error")

	customErr := NewI(code, err)

	if customErr == nil {
		t.Fatal("Expected a valid error, got nil")
	}

	if customErr.Code() != code {
		t.Errorf("Expected code %d, got %d", code, customErr.Code())
	}

	expectedMessage := "another test error"
	if customErr.Message() != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, customErr.Message())
	}
}

func TestMessageWithNoError(t *testing.T) {
	code := Int(9999993)
	e := New(code, nil)

	if e.Code() != 9999993 {
		t.Errorf("expected code 9999993, got %d", e.Code())
	}

	expectedMessage := "9999993"
	if e.Message() != expectedMessage {
		t.Errorf("expected message '%s', got '%s'", expectedMessage, e.Message())
	}
}

func TestUS(t *testing.T) {
	code := NewCode(9999994)
	message := "User-friendly message"

	Register("en-us", code, message)

	if code.Message() != message {
		t.Errorf("Expected message '%s', got '%s'", message, code.Message())
	}

	// Update the message
	code.US("Updated user-friendly message")
	retrievedMessage := code.Message()
	if retrievedMessage != "Updated user-friendly message" {
		t.Errorf("Expected updated message '%s', got '%s'", "Updated user-friendly message", retrievedMessage)
	}
}
