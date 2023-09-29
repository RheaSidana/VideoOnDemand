package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type IEncryption interface{
	CompareHashAndPassword(passwordForLogin, passwordInDB string) error
	GenerateFromPassword(password string) (string, error)
}

type encryption struct{}

func NewEncryption() IEncryption{
	return &encryption{}
}

func (r *encryption) GenerateFromPassword(password string) (string, error) {
	encryptPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("encryption unsuccessful")
	}

	password = string(encryptPassword)
	return password, nil
}

func (e *encryption) CompareHashAndPassword(passwordForLogin, passwordInDB string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordInDB),
		[]byte(passwordForLogin))
	if err != nil {
		return errors.New("incorrect credentials")
	}
	return nil
}
