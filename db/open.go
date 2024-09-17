package db

import (
	"database/sql"

	"github.com/Twintat/randomExams/config"
	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	return sql.Open("sqlite3", config.DBpath())
}
