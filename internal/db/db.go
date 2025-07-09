package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	os.MkdirAll("./data", os.ModePerm)
	return sql.Open("sqlite3", path)
}
