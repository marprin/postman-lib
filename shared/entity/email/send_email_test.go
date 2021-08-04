package email

import (
	"testing"

	"github.com/marprin/postman-lib/shared/constants"
	"github.com/stretchr/testify/assert"
)

func Test_SendEmailRequest(t *testing.T) {
	t.Run("should return error as missing params", func(t *testing.T) {
		p := SendEmailRequest{}
		err := p.Validate()
		assert.Equal(t, constants.ErrFromEmailIsRequired, err)

		p.FromEmail = "ABC"
		err = p.Validate()
		assert.Equal(t, constants.ErrFromEmailIsNotValid, err)

		p.FromEmail = "abc@gmail.com"
		err = p.Validate()
		assert.Equal(t, constants.ErrToEmailIsRequired, err)

		p.ToEmail = "xyz"
		err = p.Validate()
		assert.Equal(t, constants.ErrToEmailIsNotValid, err)

		p.ToEmail = "xyz@gmail.com"
		err = p.Validate()
		assert.Equal(t, constants.ErrSubjectIsRequired, err)

		p.Subject = "Welcome party"
		err = p.Validate()
		assert.Equal(t, constants.ErrBodyIsRequired, err)
	})

	t.Run("should nil error", func(t *testing.T) {
		p := SendEmailRequest{
			FromEmail: "vel@gmail.com",
			ToEmail:   "bli@gmail.com",
			Subject:   "Welcome here",
			Body:      "Welcome body",
		}
		err := p.Validate()
		assert.Nil(t, err)
	})
}
