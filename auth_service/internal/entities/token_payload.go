package entities

type TokenPayload struct {
	accountID string
	userRole  Role
}

func newTokenPayload(accountID string, userRole Role) *TokenPayload {
	return &TokenPayload{
		accountID: accountID,
		userRole:  userRole,
	}
}

func NewTokenPayloadFromRawDatabaseDocument(accountID string, userRole Role) *TokenPayload {
	return newTokenPayload(accountID, userRole)
}

func NewTokenPayloadFromRawTokenClaims(accountID string, userRole Role) *TokenPayload {
	return newTokenPayload(accountID, userRole)
}

func (p *TokenPayload) GetAccountID() string {
	return p.accountID
}

func (p *TokenPayload) GetUserRole() Role {
	return p.userRole
}
