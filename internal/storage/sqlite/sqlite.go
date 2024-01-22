package sqlite

import (
	"EWallet/internal/lib/api/response"
	"EWallet/internal/lib/randomStr"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
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
	    balance DECIMAL(10, 2))	
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	secSqlRequest, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS transactions (
	transactionTime TEXT,
	fromWallet TEXT NOT NULL, 
	toWallet TEXT NOT NULL, 
	amount DECIMAL(10, 2),
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

	SqlRequest, err := s.dataBase.Prepare("INSERT INTO wallets(walletID, balance) VALUES(?, 100)")
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

func (s *Storage) Transfer(fromWallet, toWallet, amount string) error {
	const op = "storage.sqlite.Transfer"

	SqlRequest, err := s.dataBase.Prepare(`
		SELECT balance
		FROM wallets
		WHERE walletID = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	var resBalance float64
	err = SqlRequest.QueryRow(fromWallet).Scan(&resBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	fAmount, err := strconv.ParseFloat(amount, 4)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	if resBalance < fAmount {
		return fmt.Errorf("%s: execute statement: %w", op, errors.New("not enough money "+
			"to do transfer"))
	}

	RedSqlRequest, err := s.dataBase.Prepare(`
		UPDATE wallets
		SET balance = ROUND(balance - ?, 2)
		WHERE walletID = ?;	
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	AddSqlRequest, err := s.dataBase.Prepare(`
		UPDATE wallets
		SET balance = ROUND(balance + ?, 2)
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

	SqlRequest, err = s.dataBase.Prepare(`
		INSERT INTO transactions(transactionTime, fromWallet, toWallet, amount) 
		VALUES(?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	_, err = SqlRequest.Exec(time.Now().Format(time.RFC3339), fromWallet, toWallet, amount)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}

func (s *Storage) ShowHistory(walletID string) ([]response.RespTransaction, error) {
	const op = "storage.sqlite.ShowHistory"

	SqlRequest, err := s.dataBase.Prepare(`
		SELECT *
		FROM transactions
		WHERE  fromWallet = ? OR toWallet = ?
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var historyResponse []response.RespTransaction
	rows, err := SqlRequest.Query(walletID, walletID)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	for rows.Next() {
		var rec response.RespTransaction
		if err = rows.Scan(&rec.TransactionTime, &rec.FromWallet, &rec.ToWallet, &rec.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		historyResponse = append(historyResponse, rec)
	}
	return historyResponse, nil
}

func (s *Storage) ShowWallet(walletID string) (response.Wallet, error) {
	const op = "storage.sqlite.ShowWallet"

	SqlRequest, err := s.dataBase.Prepare(`
		SELECT *
		FROM wallets
		WHERE  walletID = ?
	`)
	var resWallet response.Wallet
	rows, err := SqlRequest.Query(walletID)
	if err != nil {
		return response.Wallet{}, fmt.Errorf("%s: query statement: %w", op, err)
	}
	rows.Next()
	if err = rows.Scan(&resWallet.WalletID, &resWallet.Balance); err != nil {
		return response.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}
	return resWallet, nil
}
