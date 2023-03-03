package user

import (
	"github.com/YFatMR/go_messenger/user_service/apientity"
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

type manager struct {
	hasher   Hasher
	verifier Verifier
}

func NewPasswordManager(hasher Hasher, verifier Verifier) apientity.PasswordManager {
	return &manager{
		hasher:   hasher,
		verifier: verifier,
	}
}

func DefaultPasswordManager() apientity.PasswordManager {
	return NewPasswordManager(HashPassword, VerifyPassword)
}

func (m *manager) HashPassword(password string) (string, error) {
	return m.hasher(password)
}

func (m *manager) VerifyPassword(hashedPassword string, password string) error {
	return m.verifier(hashedPassword, password)
}
