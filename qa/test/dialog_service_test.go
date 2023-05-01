//go:build test
// +build test

package test

// import (
// 	"context"
// 	"testing"

// 	"github.com/YFatMR/go_messenger/protocol/pkg/proto"
// 	"github.com/stretchr/testify/suite"
// )

// type DialogTestSuite struct {
// 	userManager   UserManager
// 	dialogManager DialogManager
// 	suite.Suite
// }

// func TestDialogTestSuite(t *testing.T) {
// 	suite.Run(t, new(DialogTestSuite))
// }

// func (s *DialogTestSuite) TestCreateDialog() {
// 	ctx := context.Background()
// 	require := s.Require()

// 	_, token, err := s.userManager.NewAuthorizedUser(ctx)
// 	require.NoError(err)
// 	ctx = s.userManager.NewContextWithToken(ctx, token)

// 	dialog, _, err := s.dialogManager.CreateDialogWithNewUser(ctx)
// 	require.NoError(err)
// 	require.NotNil(dialog)
// }

// func (s *DialogTestSuite) TestGetDialogs() {
// 	ctx := context.Background()
// 	require := s.Require()

// 	_, token, err := s.userManager.NewAuthorizedUser(ctx)
// 	require.NoError(err)
// 	ctx = s.userManager.NewContextWithToken(ctx, token)
// 	dialogsCount := 5

// 	for i := 0; i < dialogsCount; i++ {
// 		dialog, _, err := s.dialogManager.CreateDialogWithNewUser(ctx)
// 		require.NoError(err)
// 		require.NotNil(dialog)
// 	}

// 	dialogs, err := frontServicegRPCClient.GetDialogs(ctx, &proto.GetDialogsRequest{
// 		Offset: 0,
// 		Limit:  10,
// 	})
// 	require.NoError(err)
// 	require.NotNil(dialogs)

// 	require.Equal(dialogsCount, len(dialogs.Dialogs))
// }
