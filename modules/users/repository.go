package users

import (
	"vod/model"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user model.User) (model.User, error)
	Find(user model.User) (model.User, error)
}

type repository struct {
	client *gorm.DB
	encryption IEncryption
}

func NewRepository(client *gorm.DB) Repository {
	return &repository{client: client}
}

func (r *repository) Create(user model.User) (model.User, error) {
	r.encryption = NewEncryption()
	encryptPassword, err := r.encryption.GenerateFromPassword(user.Password)
	if err != nil {
		return model.User{}, err
	}
	user.Password = encryptPassword

	role, err := setRole(user.Role)
	if err != nil {
		return model.User{}, err
	}
	user.Role = role

	result := r.client.Create(&user)

	if result.Error != nil {
		return model.User{}, result.Error
	}

	return user, nil
}

func (r *repository) Find(userForLogin model.User) (model.User, error) {
	var userInDB model.User
	res := r.client.Where("users.email=?", userForLogin.Email).Find(&userInDB)
	if res.Error != nil {
		return model.User{}, res.Error
	}

	r.encryption = NewEncryption()
	err := r.encryption.CompareHashAndPassword(userForLogin.Password, userInDB.Password)
	if err != nil {
		return model.User{}, err
	}

	return userInDB, nil
}
