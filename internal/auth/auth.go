package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   userID.String(),
		Issuer:    "chirpy",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// After creating your token with jwt.NewWithClaims
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	// Create a new jwt.RegisteredClaims to hold the parsed claims
	claims := &jwt.RegisteredClaims{}

	// tokenString =  The JWT string to parse
	// claims = Where to put the extracted claims
	// func(token *jwt.Token) = This function provides the key to check the signature
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	sub := claims.Subject // takes the subject field of claims struct (which is the USER.ID)
	if sub == "" {
		return uuid.Nil, err
	}

	notstr, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, err
	}

	return notstr, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	token_str := headers.Get("Authorization")

	if token_str == "" {
		return "", errors.New("invalid token")
	}
	if !strings.HasPrefix(token_str, "Bearer ") {
		return "", errors.New("invalid token")
	}

	tkn := strings.TrimSpace(token_str[7:])

	return tkn, nil
}
