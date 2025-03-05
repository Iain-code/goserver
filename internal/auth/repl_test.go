package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {

	password := "pissword"
	pass, err := HashPassword(password)
	if err != nil {
		t.Error("hash password failed")
	}
	err = CheckPasswordHash(password, pass)
	if err != nil {
		t.Error("passwords didnt match")
	}
}
