package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)


func InitPostgre(config PostgreConfig) *sql.DB {
	db, err := sql.Open("postgres", config.ConnectionString)

	if err != nil {
		panic(err)
	}

	return db
}
