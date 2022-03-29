package mysql

import (
	"database/sql"
)

type Tx struct {
	Conn *sql.DB
}

func (tr *Tx) GetDBConn() *sql.DB {
	return tr.Conn
}
