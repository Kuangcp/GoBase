package service

import (
	"github.com/kuangcp/gobase/mybook/dal"
	"github.com/kuangcp/gobase/mybook/domain"
	"github.com/wonderivan/logger"
)

func AddBookkeeping(book *domain.BookKeeping) {
	db := dal.GetDB()

	create := db.Create(book)
	logger.Info(create)
}
