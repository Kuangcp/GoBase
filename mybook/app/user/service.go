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
