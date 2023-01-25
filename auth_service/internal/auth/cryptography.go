package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword), err
}

func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type (
	Hasher   = func(string) (string, error)
	Verifier = func(hashedPassword string, password string) error
)

type PasswordValidator struct {
	hasher   Hasher
	verifier Verifier
}

func NewPasswordValidator(hasher Hasher, verifier Verifier) *PasswordValidator {
	return &PasswordValidator{
		hasher:   hasher,
		verifier: verifier,
	}
}

func NewDefaultPasswordValidator() *PasswordValidator {
	return NewPasswordValidator(HashPassword, VerifyPassword)
}

func (v *PasswordValidator) GetHasher() Hasher {
	return v.hasher
}

func (v *PasswordValidator) GetVerifier() Verifier {
	return v.verifier
}
