package book

import (
	"mybook/app/common/dal"

	"github.com/wonderivan/logger"
)

func AddBookkeeping(book *BookKeeping) {
	db := dal.GetDB()

	create := db.Create(book)
	logger.Info(create)
}
