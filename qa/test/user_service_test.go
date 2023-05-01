//go:build test
// +build test

package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/configs/cviper"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	userManager UserManager
	client      UserServiceHTTPClient
	suite.Suite
}

func TestUserTestSuite(t *testing.T) {
	config := cviper.New()
	config.SetConfigFile(envFile)
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
		panic(err)
	}

	qaHost := config.GetStringRequired("QA_HOST")
	restFrontServiceAddress := "http://" + qaHost + ":" + config.GetStringRequired("PUBLIC_REST_FRONT_SERVICE_PORT")

	client := UserServiceHTTPClient{
		HttpClient: HttpClient{
			BaseUrl: restFrontServiceAddress,
			Client: http.Client{
				Timeout: time.Second * 10,
			},
		},
	}

	obj := &UserTestSuite{
		userManager: UserManager{
			client: client,
		},
		client: client,
	}
	suite.Run(t, obj)
}

func (s *UserTestSuite) TestUserCreation() {
	ctx := context.Background()
	require := s.Require()

	_, _, err := s.userManager.NewUnauthorizedUser(ctx)
	require.NoError(err)
}

func (s *UserTestSuite) TestUserTokenGeneration() {
	ctx := context.Background()
	require := s.Require()

	_, _, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)
}

func (s *UserTestSuite) TestGetUserInfoWithValidToken() {
	ctx := context.Background()
	require := s.Require()

	autorizedUserID, token, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)

	ctx = s.userManager.NewContextWithToken(ctx, token)
	_, err = s.client.GetUserByID(autorizedUserID)
	require.NoError(err)
}

func (s *UserTestSuite) TestGetUserInfoWithoutToken() {
	ctx := context.Background()
	require := s.Require()

	autorizedUserID, _, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)

	_, err = s.client.GetUserByID(autorizedUserID)
	require.Error(err)
}
