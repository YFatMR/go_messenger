package entity

type Account struct {
	UserID         UserID
	Email          string
	HashedPassword string
	Role           UserRole
	Nickname       string
	Name           string
	Surname        string
	Github         string
	Linkedin       string
	PublicEmail    string
	Status         string
}
