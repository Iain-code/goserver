package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestMakeJWT(t *testing.T) {

	userID := uuid.New()
	testSecret := "your-test-secret-key"

	token, err := MakeJWT(userID, testSecret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	ID, err := ValidateJWT(token, testSecret)
	if err != nil {
		t.Fatal(err)
	}
	if ID != userID {
		t.Fatalf("expected user ID %v, got %v", userID, ID)
	}

}
