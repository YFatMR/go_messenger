package passwordmanager

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

type Manager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword string, password string) error
}

type manager struct {
	hasher   Hasher
	verifier Verifier
}

func New(hasher Hasher, verifier Verifier) Manager {
	return &manager{
		hasher:   hasher,
		verifier: verifier,
	}
}

func Default() Manager {
	return New(HashPassword, VerifyPassword)
}

func (m *manager) HashPassword(password string) (string, error) {
	return m.hasher(password)
}

func (m *manager) VerifyPassword(hashedPassword string, password string) error {
	return m.verifier(hashedPassword, password)
}
