package auth

import "errors"

var (
	ErrInvalidAccessToken    = errors.New("got invalid access token")
	ErrTokenGenerationFailed = errors.New("token generation failed")
)
