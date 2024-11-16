package util

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPassword(plain string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(plain), 11)
	return string(bs), err
}

func BcryptCompare(passwordHash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plain)) == nil
}

func PasswordComplexityVerify(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
