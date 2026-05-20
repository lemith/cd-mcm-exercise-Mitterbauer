package main

import (
	"os"
	"testing"
)

func TestGetEnvWithValue(t *testing.T) {
	os.Setenv("TEST_KEY", "hello")
	defer os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "fallback")
	if result != "hello" {
		t.Errorf("expected hello, got %s", result)
	}
}

func TestGetEnvWithFallback(t *testing.T) {
	os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback, got %s", result)
	}
}

func TestGetEnvWithEmptyValue(t *testing.T) {
	os.Setenv("TEST_KEY", "")
	defer os.Unsetenv("TEST_KEY")
	result := getEnv("TEST_KEY", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback for empty value, got %s", result)
	}
}

func TestGetEnvWithWhitespaceValue(t *testing.T) {
	os.Setenv("TEST_KEY", "   ")
	defer os.Unsetenv("TEST_KEY")
	result := getEnv("TEST_KEY", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback for whitespace value, got %s", result)
	}
}

func TestGetEnvWithWhitespaceFallback(t *testing.T) {
	os.Unsetenv("TEST_KEY")
	result := getEnv("TEST_KEY", "   ")
	if result != "   " {
		t.Errorf("expected whitespace fallback, got %s", result)
	}
}

func TestGetEnvWithEmptyFallback(t *testing.T) {
	os.Unsetenv("TEST_KEY")
	result := getEnv("TEST_KEY", "")
	if result != "" {
		t.Errorf("expected empty fallback, got %s", result)
	}
}

func TestGetEnvWithWhitespaceKey(t *testing.T) {
	os.Setenv("   ", "value")
	defer os.Unsetenv("   ")
	result := getEnv("   ", "fallback")
	if result != "value" {
		t.Errorf("expected value for whitespace key, got %s", result)
	}
}

func TestGetEnvWithEmptyKey(t *testing.T) {
	os.Setenv("", "value")
	defer os.Unsetenv("")
	result := getEnv("", "fallback")
	if result != "value" {
		t.Errorf("expected value for empty key, got %s", result)
	}
}

func TestGetEnvWithWhitespaceKeyAndFallback(t *testing.T) {
	os.Unsetenv("   ")
	result := getEnv("   ", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback for unset whitespace key, got %s", result)
	}
}

func TestGetEnvWithEmptyKeyAndFallback(t *testing.T) {
	os.Unsetenv("")
	result := getEnv("", "fallback")
	if result != "fallback" {
		t.Errorf("expected fallback for unset empty key, got %s", result)
	}
}

func TestGetEnvWithWhitespaceKeyAndEmptyFallback(t *testing.T) {
	os.Unsetenv("   ")
	result := getEnv("   ", "")
	if result != "" {
		t.Errorf("expected empty fallback for unset whitespace key, got %s", result)
	}
}

func TestGetEnvWithEmptyKeyAndEmptyFallback(t *testing.T) {
	os.Unsetenv("")
	result := getEnv("", "")
	if result != "" {
		t.Errorf("expected empty fallback for unset empty key, got %s", result)
	}
}
