package entity

type Credential struct {
	Email          string
	HashedPassword string
	Role           UserRole
}

func CredentialFromUnsafeCredential(usafeCredential *UnsafeCredential, hashedPassword string) *Credential {
	return &Credential{
		Email:          usafeCredential.Email,
		HashedPassword: hashedPassword,
		Role:           usafeCredential.Role,
	}
}
