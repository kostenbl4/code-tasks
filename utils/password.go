package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword принимает пароль и возвращает его хеш
func HashPassword(pass string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

// CheckPassword сравнивает хешированный пароль с введенным
func CheckPassword(hpass, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hpass), []byte(pass))
	return err == nil
}
