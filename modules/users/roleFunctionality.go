package users

import "errors"

var UserRole = []string{
	"ADMIN",
	"CUSTOMER",
}

func UserRoleAdmin() string {
	return UserRole[0]
}

func UserRoleCustomer() string {
	return UserRole[1]
}

func setRole(role string) (string, error) {
	if role == "admin" {
		return UserRoleAdmin(), nil
	} else if role == "customer" {
		return UserRoleCustomer(), nil
	}

	return "", errors.New("inavlid role provided")
}