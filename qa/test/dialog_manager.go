package test

type DialogManager struct {
	userManager UserManager
}

// func (m *DialogManager) CreateDialogWithNewUser(ctx context.Context) (*proto.Dialog, context.Context, error) {
// 	userID, token, err := m.userManager.NewAuthorizedUser(ctx)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	retCtx := m.userManager.NewContextWithToken(context.Background(), token)
// 	dialog, err := frontServicegRPCClient.CreateDialogWith(ctx, userID)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return dialog, retCtx, nil
// }
