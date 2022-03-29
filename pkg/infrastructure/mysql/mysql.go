package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// blank import for MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const driverName = "mysql"

// Conn 各repositoryで利用するDB接続(Connection)情報
type SQLHandler struct {
	Conn *sql.DB
}

func NewSQLHandler() SQLHandler {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Printf("Failed to load env file: %v", err)
	}
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	conn, err := sql.Open(driverName,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database))
	if err != nil {
		log.Fatal(err)
	}

	return SQLHandler{
		Conn: conn,
	}
}

// a interface that implements all functions in sql.DB and sql.Tx
type DBTx interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
