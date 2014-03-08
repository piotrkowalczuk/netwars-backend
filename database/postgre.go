package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var Postgre = PostgreDB{}

type PostgreDB struct {
	db *sql.DB
}

func (self *PostgreDB) setDb(db *sql.DB ) {
	self.db = db
}

func (self *PostgreDB) getDb() *sql.DB {
	return self.db
}

func InitPostgre() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres.123@localhost/netwars-org?sslmode=disable")
	checkError(err, "sql.Open failed")

	Postgre.setDb(db)
}
