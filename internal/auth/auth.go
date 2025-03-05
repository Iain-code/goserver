package auth

import (
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
