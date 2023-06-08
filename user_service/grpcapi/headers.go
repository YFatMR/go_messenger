package grpcapi

type Headers struct {
	UserID string
}

func NewHeaders(userID string) Headers {
	return Headers{UserID: userID}
}
