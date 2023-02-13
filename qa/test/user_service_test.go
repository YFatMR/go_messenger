package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	userManager UserManager
	suite.Suite
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
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
	_, err = frontServicegRPCClient.GetUserByID(ctx, autorizedUserID)
	require.NoError(err)
}

func (s *UserTestSuite) TestGetUserInfoWithoutToken() {
	ctx := context.Background()
	require := s.Require()

	autorizedUserID, _, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)

	_, err = frontServicegRPCClient.GetUserByID(ctx, autorizedUserID)
	require.Error(err)
}
