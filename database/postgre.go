package database

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)


func InitializeGorp(config PostgreConfig) *gorp.DbMap {
	db, _ := sql.Open("postgres", config.ConnectionString)

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	return dbMap
}
