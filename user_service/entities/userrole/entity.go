package userrole

type Entity struct {
	name string
}

func (e *Entity) GetName() string {
	return e.name
}

var (
	User  = Entity{name: "user"}
	Admin = Entity{name: "admin"}
)

func FromString(role string) (*Entity, error) {
	switch role {
	case User.GetName():
		return &User, nil
	case Admin.GetName():
		return &Admin, nil
	default:
		return nil, ErrUndefinedRole
	}
}
