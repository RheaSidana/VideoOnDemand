package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func generateFromPassword(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("encryption unsuccessful")
	}

	password = string(encryptPassword)
	return password, nil
}

func compareHashAndPassword(passwordForLogin , passwordInDB string) (error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordInDB), 
		[]byte(passwordForLogin))
	if err != nil {
		return errors.New("incorrect credentials")
	}
	return nil
}