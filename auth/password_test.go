package auth

import "testing"

func TestHashPassowrd(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password %v", err)
	}
	if hash == "" {
		t.Error("Error hashing can not be empty")
	}
	if hash == "password" {
		t.Error("Expected hash to be different from password")
	}
}

func TestComparePasswords(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("Error hashing password %v", err)
	}
	if hash == "" {
		t.Error("Error hashing can not be empty")
	}
	if hash == "password" {
		t.Error("Expected hash to be different from password")
	}

	if !ComparePasswords(hash, []byte("password")) {
		t.Error("Expected password to match hash")
	}
	if ComparePasswords(hash, []byte("not-password")) {
		t.Error("Expected password to not match hash")
	}
}
