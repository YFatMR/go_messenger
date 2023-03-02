package userservice_test

import (
	"context"
	"testing"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/core/pkg/jwtmanager"
	"github.com/YFatMR/go_messenger/user_service/entities/account"
	"github.com/YFatMR/go_messenger/user_service/entities/unsafecredential"
	"github.com/YFatMR/go_messenger/user_service/entities/userrole"
	"github.com/YFatMR/go_messenger/user_service/passwordmanager"
	"github.com/YFatMR/go_messenger/user_service/repositories"
	"github.com/YFatMR/go_messenger/user_service/userservice"
)

func BenchmarkTokenGeneration(b *testing.B) {
	var repository repositories.UserRepository = &MockUserRepository{
		GetAccountByLoginResponse: GetAccountByLoginResponseData{
			Account: account.New("id", "login", "hashed_password", &userrole.User, "nickname", "name", "surname"),
			Error:   nil,
		},
	}

	var passwordManager passwordmanager.Manager = &MockPasswordManager{
		HashPasswordResponse: HashPasswordResponseData{
			HashedPassword: "hashed_password",
			Error:          nil,
		},
		VerifyPasswordResponse: VerifyPasswordResponseData{
			Error: nil,
		},
	}

	// Create real jwtManager to bench execution
	jwtManager := jwtmanager.New("SuperSecretKey", time.Second*10, czap.NewNop())

	userService := userservice.New(repository, passwordManager, jwtManager, czap.NewNop())
	unsafeCredential := unsafecredential.New("login", "hashed_password", &userrole.User)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userService.GenerateToken(context.Background(), unsafeCredential)
		if err != nil {
			panic(err)
		}
	}
}
