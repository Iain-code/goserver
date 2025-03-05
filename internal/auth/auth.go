package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	passSlice := []byte(password)
	pass, err := bcrypt.GenerateFromPassword(passSlice, 10)
	if err != nil {
		return "", err
	}

	itm := string(pass)
	return itm, nil

}

func CheckPasswordHash(password, hash string) error {

	passSlice := []byte(password)
	hashSlice := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashSlice, passSlice)
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims)

	return "", nil
}
