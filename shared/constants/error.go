package constants

import "errors"

var (
	ErrFromEmailIsRequired = errors.New("From email is required")
	ErrFromEmailIsNotValid = errors.New("From email is not valid")
	ErrFromAliasIsRequired = errors.New("From alias is required")
	ErrToEmailIsRequired   = errors.New("To email is required")
	ErrToEmailIsNotValid   = errors.New("To email is not valid")
	ErrSubjectIsRequired   = errors.New("Subject is required")
	ErrBodyIsRequired      = errors.New("Body is required")
)
