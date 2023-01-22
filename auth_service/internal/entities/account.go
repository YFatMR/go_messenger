package entities

import proto "github.com/YFatMR/go_messenger/protocol/pkg/proto"

type AccountID struct {
	accountID string
}

func newAccountID(accountID string) *AccountID {
	return &AccountID{
		accountID: accountID,
	}
}

func NewAccountIDFromProtobuf(accountID *proto.AccountID) (*AccountID, error) {
	if accountID == nil || accountID.GetID() == "" {
		return nil, ErrWrongRequestFormat
	}
	return newAccountID(accountID.GetID()), nil
}

func NewAccountIDFromRawDatabaseDocument(accountID string) *AccountID {
	return newAccountID(accountID)
}

func (a *AccountID) GetID() string {
	return a.accountID
}
