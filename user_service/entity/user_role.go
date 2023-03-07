package entity

type UserRole struct {
	Name string
}

var (
	RoleUser  = UserRole{Name: "user"}
	RoleAdmin = UserRole{Name: "admin"}
)

func UserRoleFromString(role string) (*UserRole, error) {
	switch role {
	case RoleUser.Name:
		return &RoleUser, nil
	case RoleAdmin.Name:
		return &RoleAdmin, nil
	default:
		return nil, ErrUndefinedRole
	}
}
