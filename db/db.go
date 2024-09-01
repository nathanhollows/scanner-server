package db

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

func MustOpen() {
	var sqldb *sql.DB
	var err error

	dataSourceName := "scans.db"
	sqldb, err = sql.Open(sqliteshim.ShimName, dataSourceName)
	DB = bun.NewDB(sqldb, sqlitedialect.New())

	if err != nil {
		panic(err)
	}

	DB.AddQueryHook(bundebug.NewQueryHook(
		// disable the hook
		bundebug.WithEnabled(false),

		// BUNDEBUG=1 logs failed queries
		// BUNDEBUG=2 logs all queries
		bundebug.FromEnv("BUNDEBUG"),
	))
}
