package database

import (
	"database/sql"
)

type Repository struct {
	TableName string
}

func (self *Repository) FindOne(idField string, id int) (row *sql.Row) {
	row = Postgre.getDb().QueryRow( "SELECT * FROM "+self.TableName+" WHERE "+idField+" = $1", id)

	return
}

