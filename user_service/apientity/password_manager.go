package apientity

type PasswordManager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword string, password string) error
}
