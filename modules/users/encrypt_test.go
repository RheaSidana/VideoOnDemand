package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	t.Run("TestNewEncryption", func(t *testing.T) {
		enc := NewEncryption()
		assert.NotNil(t, enc)
	})

	t.Run("TestGenerateFromPassword", func(t *testing.T) {
		enc := NewEncryption()
		password := "secret"

		encryptedPassword, err := enc.GenerateFromPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedPassword)
		assert.NotEqual(t, password, encryptedPassword)
	})

	t.Run("TestCompareHashAndPasswordValid", func(t *testing.T) {
		enc := NewEncryption()
		password := "secret"

		encryptedPassword, err := enc.GenerateFromPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedPassword)

		err = enc.CompareHashAndPassword(password, encryptedPassword)
		assert.NoError(t, err)
	})

	t.Run("TestCompareHashAndPasswordInvalid", func(t *testing.T) {
		enc := NewEncryption()
		password := "secret"
		incorrectPassword := "incorrect"

		encryptedPassword, err := enc.GenerateFromPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, encryptedPassword)

		err = enc.CompareHashAndPassword(incorrectPassword, encryptedPassword)
		assert.Error(t, err)
		assert.Equal(t, "incorrect credentials", err.Error())
	})
}
