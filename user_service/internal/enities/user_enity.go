package enities

type User struct {
	name    string
	surname string
}

func NewUser(name string, surname string) *User {
	return &User{
		name:    name,
		surname: surname,
	}
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetSurname() string {
	return u.surname
}
