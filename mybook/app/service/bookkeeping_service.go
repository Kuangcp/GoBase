package service

import (
	"mybook/app/common/dal"
	"mybook/app/domain"
	"github.com/wonderivan/logger"
)

func AddBookkeeping(book *domain.BookKeeping) {
	db := dal.GetDB()

	create := db.Create(book)
	logger.Info(create)
}
