package accountservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/YFatMR/go_messenger/auth_service/internal/auth"
	"github.com/YFatMR/go_messenger/auth_service/internal/auth/jwtmanager"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/accountid"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/credential"
	"github.com/YFatMR/go_messenger/auth_service/internal/entities/tokenpayload"
	"github.com/YFatMR/go_messenger/auth_service/internal/services"
	"github.com/YFatMR/go_messenger/auth_service/internal/services/accountservice"
	"github.com/YFatMR/go_messenger/core/pkg/ulo"
)

type MockAccountRepository struct {
	accountID      string
	role           entities.Role
	hashedPassword string
}

func (r *MockAccountRepository) CreateAccount(ctx context.Context, credential *credential.Entity, role entities.Role) (
	accountID *accountid.Entity, logtash ulo.LogStash, err error,
) {
	return accountid.New(r.accountID), nil, nil
}

func (r *MockAccountRepository) GetTokenPayloadWithHashedPasswordByLogin(ctx context.Context, login string) (
	tokenPayload *tokenpayload.Entity, hashedPassword string, logtash ulo.LogStash, err error,
) {
	return tokenpayload.New(r.accountID, r.role), r.hashedPassword, nil, nil
}

func setup() (services.AccountService, *credential.Entity, error) {
	login := "123"
	password := "123"
	accountID := "63cd84a07244f00b9e54e41b"
	accountRole := entities.UserRole

	passwordValidator := auth.NewDefaultPasswordValidator()
	hashedPassword, err := passwordValidator.GetHasher()(password)
	if err != nil {
		return nil, nil, err
	}

	repository := &MockAccountRepository{
		accountID:      accountID,
		role:           accountRole,
		hashedPassword: hashedPassword,
	}
	var authManager jwtmanager.Manager = jwtmanager.New("secret", time.Hour)
	userservice := accountservice.New(repository, authManager)
	credential := credential.New(login, password, hashedPassword, passwordValidator.GetVerifier())
	return userservice, credential, nil
}

func BenchmarkTokenGeneration(b *testing.B) {
	userservice, credential, err := setup()
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := userservice.GetToken(context.Background(), credential)
		if err != nil {
			panic(err)
		}
	}
}
