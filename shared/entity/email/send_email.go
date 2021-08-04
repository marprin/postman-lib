package email

import (
	"github.com/marprin/postman-lib/pkg/strings"
	"github.com/marprin/postman-lib/shared/constants"
)

type (
	SendEmailRequest struct {
		FromEmail string
		FromName  string
		ToEmail   string
		ToName    string
		Subject   string
		Body      string
		ReplyTo   string
	}
)

func (s *SendEmailRequest) Validate() error {
	if s.FromEmail == "" {
		return constants.ErrFromEmailIsRequired
	}

	if !strings.IsValidEmail(s.FromEmail) {
		return constants.ErrFromEmailIsNotValid
	}

	if s.ToEmail == "" {
		return constants.ErrToEmailIsRequired
	}

	if !strings.IsValidEmail(s.ToEmail) {
		return constants.ErrToEmailIsNotValid
	}

	if s.Subject == "" {
		return constants.ErrSubjectIsRequired
	}

	if s.Body == "" {
		return constants.ErrBodyIsRequired
	}

	return nil
}
