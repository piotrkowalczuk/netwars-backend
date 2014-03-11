package database

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)


func InitializeGorp() *gorp.DbMap {
	db, _ := sql.Open("postgres", "postgres://postgres:postgres.123@localhost/netwars-org?sslmode=disable")

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	return dbMap
}
