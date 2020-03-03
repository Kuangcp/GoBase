package service

import (
	"github.com/kuangcp/gobase/mybook/app/dal"
	"github.com/kuangcp/gobase/mybook/app/domain"
	"github.com/wonderivan/logger"
)

func AddBookkeeping(book *domain.BookKeeping) {
	db := dal.GetDB()

	create := db.Create(book)
	logger.Info(create)
}
