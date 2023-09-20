package users

import "errors"

var UserRole = []string{
	"ADMIN",
	"CUSTOMER",
}

func userRoleAdmin() string {
	return UserRole[1]
}

func userRoleCustomer() string {
	return UserRole[2]
}

func setRole(role string) (string, error) {
	if role == "admin" {
		return userRoleAdmin(), nil
	} else if role == "customer" {
		return userRoleCustomer(), nil
	}

	return "", errors.New("inavlid role provided")
}