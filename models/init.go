package models

import (
	"context"
	"log"

	"github.com/nathanhollows/scanner-server/db"
)

func CreateTables() {
	var models = []interface{}{
		(*Scan)(nil),
		(*Tag)(nil),
	}

	for _, model := range models {
		_, err := db.DB.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}
