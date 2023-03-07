package entity

type Account struct {
	UserID         string
	Login          string
	HashedPassword string
	Role           *UserRole
	Nickname       string
	Name           string
	Surname        string
}
