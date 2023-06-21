package internal

import "errors"

var (
	ErrSigSet           = errors.New("signature already set")
	ErrChallengeInvalid = errors.New("challenge is expired or invalid")
)
