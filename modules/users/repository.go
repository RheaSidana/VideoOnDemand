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
}

func NewRepository(client *gorm.DB) Repository {
	return &repository{client: client}
}

func (r *repository) Create(user model.User) (model.User, error) {
	encryptPassword, err := GenerateFromPassword(user.Password)
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

	err := compareHashAndPassword(userForLogin.Password, userInDB.Password)
	if err != nil {
		return model.User{}, err
	}

	return userInDB, nil
}
