package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	dataBase *sql.DB
}

// New creates database in storagePath and returns it if successfully
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	dataBase, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS wallet (
	    id INTEGER PRIMARY KEY,
	    walletID TEXT NOT NULL UNIQUE,
	    balance FLOAT);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{dataBase: dataBase}, nil

}
