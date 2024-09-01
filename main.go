package main

import (
	"github.com/nathanhollows/scanner-server/db"
	"github.com/nathanhollows/scanner-server/handlers"
	"github.com/nathanhollows/scanner-server/models"
)

func main() {
	db.MustOpen()
	models.CreateTables()
	handlers.Start()
}
