package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRole(t *testing.T) {
	t.Run("ValidAdminRole", func(t *testing.T) {
		role, err := setRole("admin")
		assert.NoError(t, err)
		assert.Equal(t, UserRoleAdmin(), role)
	})

	t.Run("ValidCustomerRole", func(t *testing.T) {
		role, err := setRole("customer")
		assert.NoError(t, err)
		assert.Equal(t, UserRoleCustomer(), role)
	})

	t.Run("InvalidRole", func(t *testing.T) {
		role, err := setRole("invalid_role")
		assert.Error(t, err)
		assert.Empty(t, role)
	})

	t.Run("EmptyRole", func(t *testing.T) {
		role, err := setRole("")
		assert.Error(t, err)
		assert.Empty(t, role)
	})
}
