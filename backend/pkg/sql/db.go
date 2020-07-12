package sql

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DBer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
}

type db struct {
	connection *sql.DB
}

func NewDB(dsn string) db {
	connection, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to articleRepository : " + err.Error())
	}
	return db{connection}
}

func (db db) Exec(query string, args ...interface{}) (sql.Result, error) {
	results, err := db.connection.Exec(query, args...)
	return results, err
}
func (db db) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.connection.Query(query, args...)
	return rows, err
}
func (db db) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := db.connection.Prepare(query)
	return stmt, err
}
