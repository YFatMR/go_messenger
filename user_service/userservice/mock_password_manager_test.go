package userservice_test

type HashPasswordResponseData struct {
	HashedPassword string
	Error          error
}

type VerifyPasswordResponseData struct {
	HashedPassword string
	Error          error
}

type MockPasswordManager struct {
	HashPasswordResponse   HashPasswordResponseData
	VerifyPasswordResponse VerifyPasswordResponseData
}

func (m *MockPasswordManager) HashPassword(password string) (string, error) {
	return m.HashPasswordResponse.HashedPassword, m.HashPasswordResponse.Error
}

func (m *MockPasswordManager) VerifyPassword(hashedPassword string, password string) error {
	return m.VerifyPasswordResponse.Error
}
