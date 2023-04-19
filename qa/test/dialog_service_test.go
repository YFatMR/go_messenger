//go:build test
// +build test

package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DialogTestSuite struct {
	userManager UserManager
	suite.Suite
}

func TestDialogTestSuite(t *testing.T) {
	suite.Run(t, new(DialogTestSuite))
}

func (s *DialogTestSuite) TestCreateDialog() {
	ctx := context.Background()
	require := s.Require()

	userID, _, err := s.userManager.NewUnauthorizedUser(ctx)
	require.NoError(err)

	_, token, err := s.userManager.NewAuthorizedUser(ctx)
	require.NoError(err)
	ctx = s.userManager.NewContextWithToken(ctx, token)

	dialog, err := frontServicegRPCClient.CreateDialogWith(ctx, userID)
	require.NoError(err)
	require.NotNil(dialog)
}
