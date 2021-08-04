package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckHashValidity(t *testing.T) {
	inst := NewHash()
	pwdHash := "$2a$10$xFS4cYiQdxnZrI/qVDM17Oy1LQo0Rr1bj.Lw2QJI3WPnNMNPjJ28i"

	t.Run("should error when check hash validity", func(t *testing.T) {
		err := inst.CheckHashValidity("fake-hash", "password")
		assert.NotNil(t, err)
	})

	t.Run("should return nil when check hash validity", func(t *testing.T) {
		err := inst.CheckHashValidity(pwdHash, "password")
		assert.Nil(t, err)
	})
}

func Test_GenerateHash(t *testing.T) {
	inst := NewHash()

	t.Run("should return nil when check hash validity", func(t *testing.T) {
		hashPwd, err := inst.GenerateHash("password")
		assert.Nil(t, err)
		assert.NotNil(t, hashPwd)
	})
}
