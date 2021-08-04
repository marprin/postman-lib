package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidEmail(t *testing.T) {
	t.Run("should return false as not valid email", func(t *testing.T) {
		isValid := IsValidEmail("ABC")
		assert.False(t, isValid)
	})

	t.Run("should return true as valid email", func(t *testing.T) {
		isValid := IsValidEmail("testing@gmail.com")
		assert.True(t, isValid)
	})
}
