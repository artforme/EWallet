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

// New creates database in storagePath and returns storage with db if successfully
func New(storagePath string) (*Storage, error) {
	//this constant shows the path where we're working, so we use it if we will get mistake
	const op = "storage.sqlite.New"

	//create dataBase with the driver name sqlite3
	dataBase, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//create first request to create table wallets
	firSqlRequest, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS wallets (
	    walletID TEXT NOT NULL UNIQUE,
	    balance DECIMAL(10, 4))	
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	//create second request to create table transactions
	secSqlRequest, err := dataBase.Prepare(`
	CREATE TABLE IF NOT EXISTS transactions (
	transactionTime TEXT,
	fromWallet TEXT NOT NULL, 
	toWallet TEXT NOT NULL, 
	amount DECIMAL(10, 4),
	FOREIGN KEY(fromWallet) REFERENCES wallets(walletID), 
	FOREIGN KEY(toWallet) REFERENCES wallets(walletID))
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	//execute first request
	_, err = firSqlRequest.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	//execute second request
	_, err = secSqlRequest.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return &Storage{dataBase: dataBase}, nil

}

// CreateWallet create new createWallet and returns id of createWallet if successfully
func (s *Storage) CreateWallet() (string, error) {
	const op = "storage.sqlite.CreateWallet"

	//prepare sql request to our db
	SqlRequest, err := s.dataBase.Prepare("INSERT INTO wallets(walletID, balance) VALUES(?, 100)")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	//we use func NewRandomString from package randomStr to create unique id for createWallet
	walletID := randomStr.NewRandomString(WalletIDlenght)
	_, err = SqlRequest.Exec(walletID)
	if err != nil {
		return "", fmt.Errorf("%s: execute statement:  %w", op, err)
	}
	return walletID, nil

}

// Transfer make transactions between two wallets
func (s *Storage) Transfer(fromWallet, toWallet, amount string) error {
	const op = "storage.sqlite.Transfer"
	// check if exists and balance
	SqlRequest, err := s.dataBase.Prepare(`
		SELECT balance
		FROM wallets
		WHERE walletID = ?
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	//in this part of code we compare exists balance and amount money that we are going to get from balance
	//if balance < amount
	//we don't do transaction
	var resBalance string
	err = SqlRequest.QueryRow(fromWallet).Scan(&resBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("outgoing wallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	//formatting string to float
	fAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	fResBalance, err := strconv.ParseFloat(resBalance, 64)
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	//compare balance and amount
	if fResBalance < fAmount {
		return fmt.Errorf("%s: execute statement: %w", op, errors.New("not enough money "+
			"to complete transfer"))
	}

	// check if exists
	SqlRequest, err = s.dataBase.Prepare(`
		SELECT walletID
		FROM wallets
		WHERE walletID = ?
	`)
	err = SqlRequest.QueryRow(toWallet).Scan(&resBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("target wallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// prepare transfer
	RedSqlRequest, err := s.dataBase.Prepare(`
		UPDATE wallets
		SET balance = ROUND(balance - ?, 4)
		WHERE walletID = ?;	
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	AddSqlRequest, err := s.dataBase.Prepare(`
		UPDATE wallets
		SET balance = ROUND(balance + ?, 4)
		WHERE walletID = ?;
	`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	// complete transfer
	_, err = RedSqlRequest.Exec(amount, fromWallet)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	_, err1 := AddSqlRequest.Exec(amount, toWallet)
	if err1 != nil {
		//if got mistake we undo previous step
		_, err = AddSqlRequest.Exec(amount, fromWallet)
		if err != nil {
			return errors.New("fatal error")
		}

		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("createWallet not found")
		}
		return fmt.Errorf("%s: execute statement: %w", op, err1)
	}

	// in this part of code we add info to the transactions table
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

// ShowHistory shows the history of all transactions with specific createWallet
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
	// create slice that contains response structures
	var historyResponse []response.RespTransaction
	// get rows form request
	rows, err := SqlRequest.Query(walletID, walletID)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	// until rows run out
	for rows.Next() {
		var rec response.RespTransaction
		// record to temp struct rec info if successfully
		if err = rows.Scan(&rec.TransactionTime, &rec.FromWallet, &rec.ToWallet, &rec.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		historyResponse = append(historyResponse, rec)
	}
	return historyResponse, nil
}

// ShowWallet shows specific createWallet
func (s *Storage) ShowWallet(walletID string) (response.Wallet, error) {
	const op = "storage.sqlite.ShowWallet"

	SqlRequest, err := s.dataBase.Prepare(`
		SELECT *
		FROM wallets
		WHERE  walletID = ?
	`)
	var resWallet response.Wallet
	row, err := SqlRequest.Query(walletID)
	if err != nil {
		return response.Wallet{}, fmt.Errorf("%s: query statement: %w", op, err)
	}
	row.Next()
	//scan row and record it to resWallet
	if err = row.Scan(&resWallet.WalletID, &resWallet.Balance); err != nil {
		return response.Wallet{}, fmt.Errorf("%s: %w", op, err)
	}
	return resWallet, nil
}
