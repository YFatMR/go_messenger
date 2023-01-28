package entities

type TokenPayload struct {
	accountID string
	userRole  Role
}

func NewTokenPayload(accountID string, userRole Role) *TokenPayload {
	return &TokenPayload{
		accountID: accountID,
		userRole:  userRole,
	}
}

func (p *TokenPayload) GetAccountID() string {
	return p.accountID
}

func (p *TokenPayload) GetUserRole() Role {
	return p.userRole
}
