package user

import (
	"mybook/app/common/dal"
	"mybook/app/common/util"
)

func ListUsers() []*User {
	db := dal.GetDB()

	var accounts []*User
	e := db.Find(&accounts).Error
	util.RecordError(e)

	return accounts
}

func QueryUserMap() map[uint]User {
	users := ListUsers()
	cache := make(map[uint]User)
	for _, user := range users {
		cache[user.ID] = *user
	}
	return cache
}
