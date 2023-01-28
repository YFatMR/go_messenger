package entities

import "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type AccountID struct {
	accountID string
}

func NewAccountID(accountID string) *AccountID {
	return &AccountID{
		accountID: accountID,
	}
}

func NewAccountIDFromProtobuf(accountID *proto.AccountID) (*AccountID, error) {
	if accountID == nil || accountID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return NewAccountID(accountID.GetID()), nil
}

func (u *AccountID) GetID() string {
	return u.accountID
}
