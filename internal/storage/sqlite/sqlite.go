package sqlite

import (
	"EWallet/internal/lib/randomStr"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	WalletIDlenght = 16
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
	firSqlRequest, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS wallets (
	    walletID TEXT NOT NULL UNIQUE,
	    balance FLOAT)	
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	secSqlRequest, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS transactions (
	transactionTime TEXT,
	fromWallet TEXT NOT NULL UNIQUE, 
	toWallet TEXT NOT NULL UNIQUE, 
	amount FLOAT,
	FOREIGN KEY(fromWallet) REFERENCES wallets(walletID), 
	FOREIGN KEY(toWallet) REFERENCES wallets(walletID))
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = firSqlRequest.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	_, err = secSqlRequest.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return &Storage{dataBase: dataBase}, nil

}

func (s *Storage) CreateWallet() (string, error) {
	const op = "storage.sqlite.CreateWallet"

	SqlRequest, err := s.dataBase.Prepare("INSERT INTO wallets(wallet, balance) VALUES(?, 100)")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement:  %w", op, err)
	}

	walletID := randomStr.NewRandomString(WalletIDlenght)
	_, err = SqlRequest.Exec(walletID)
	if err != nil {
		return "", fmt.Errorf("%s: execute statement:  %w", op, err)
	}
	return walletID, nil

}

func (s *Storage) Transfer(fromWallet, toWallet string, amount float32) error {
	const op = "storage.sqlite.Transfer"

	SqlRequest, err := s.dataBase.Prepare(`
		SELECT balance
		FROM wallet
		WHERE walletID = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	var resBalance float32
	err = SqlRequest.QueryRow(fromWallet).Scan(&resBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	if resBalance < amount {
		return fmt.Errorf("%s: execute statement: %w", op, errors.New("not enough money "+
			"to do transfer"))
	}

	RedSqlRequest, err := s.dataBase.Prepare(`
		UPDATE table_name wallet
		SET balance = balance - ?
		WHERE walletID = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	AddSqlRequest, err := s.dataBase.Prepare(`
		UPDATE table_name wallet
		SET balance = balance + ?
		WHERE walletID = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	_, err = RedSqlRequest.Exec(amount, fromWallet)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	_, err1 := AddSqlRequest.Exec(amount, toWallet)
	if err1 != nil {

		_, err = AddSqlRequest.Exec(amount, fromWallet)
		if err != nil {
			return errors.New("fatal error")
		}

		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err1)
	}
	return nil
}
