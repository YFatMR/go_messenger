package entity

type Credential struct {
	Login          string
	HashedPassword string
	Role           *UserRole
}

func CredentialFromUnsafeCredential(usafeCredential *UnsafeCredential, hashedPassword string) *Credential {
	return &Credential{
		Login:          usafeCredential.Login,
		HashedPassword: hashedPassword,
		Role:           usafeCredential.Role,
	}
}
