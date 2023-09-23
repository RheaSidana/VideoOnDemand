package data

import (
	"strconv"
	"vod/initializer"
	"vod/model"
	"vod/modules/users"
)

func userData() []model.User {
	pswd, _ := users.GenerateFromPassword("apsswRd#2")
	role := users.UserRoleAdmin()

	userToAdd := model.User{
		Name:     "Name",
		Email:    "name",
		Password: pswd,
		Role:     role,
	}

	var usersList []model.User
	for i := 1; i <= 2; i++ {
		user := userToAdd
		user.Name = user.Name + " " + strconv.Itoa(i)
		user.Email = user.Email + strconv.Itoa(i) + "@example.com"

		usersList = append(usersList, user)
	}

	return usersList
}

func AddAdminUsersToDB() []model.User {
	var usersData []model.User
	for _, user := range userData() {
		if initializer.Db.Where(
			"email=?", user.Email,
		).Find(
			&user,
		).RowsAffected == 1 {
			continue
		}
		initializer.Db.Create(&user)

		usersData = append(usersData, user)
	}

	return usersData
}
