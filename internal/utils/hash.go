package utils

import (
	"golang.org/x/crypto/bcrypt"
)
// Auto-generated swagger comments for HashPassword
// @Summary Auto-generated summary for HashPassword
// @Description Auto-generated description for HashPassword — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func HashPassword(pw string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(h), err
}

func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
