package service

import (
	"log"

	"github.com/kuangcp/gobase/myth-bookkeeping/db"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
)

func QueryAllAccounts() []domain.Account {
	connection := db.GetConnection()
	defer connection.Close()

	rows, err := connection.DB.Query("select * from account")
	if err != nil {
		log.Fatal(err)
	}

	var result []domain.Account
	for rows.Next() {
		account := domain.Account{}
		err := rows.Scan(&account.Id, &account.Name, &account.InitAmount,
			&account.CreateTime, &account.UpdateTime, &account.IsDeleted)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, account)
	}
	return result
}
