package postgres

import (
	"fmt"
	"log"

	brewery "github.com/antschmidt/brewery-backend"
)

type TransactionService struct {
	client *Client
}

func (ts *TransactionService) Add(id int, units float64) (*brewery.Transaction, error) {
	var transaction brewery.Transaction
	err := ts.client.Open()
	if err != nil {
		return nil, err
	}
	defer ts.client.db.Close()
	query := "INSERT INTO transactions (memstat_id, raw_units) VALUES ($1, $2) returning timestamp;"
	err = ts.client.db.QueryRow(query, id, units).Scan(&transaction.Timestamp)
	if err != nil {
		log.Printf("Failed to insert transaction or grab the timestamp from said transaction from\nID: %v\nRawUnits: %v\n", id, units)
		return nil, err
	}
	transaction.ID = id
	transaction.RawUnits = units
	return &transaction, nil

}

func (ts *TransactionService) Remove(t *brewery.Transaction) error {
	err := ts.client.Open()
	if err != nil {
		return err
	}
	defer ts.client.db.Close()
	query := "DELETE FROM member_transactions where memstat_id=$1 and timestamp=$2 and raw_units=$3;"
	_, err = ts.client.db.Exec(query, t.ID, t.Timestamp, t.RawUnits)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TransactionService) Transactions(id int) ([]brewery.Transaction, error) {
	var transactions []brewery.Transaction
	err := ts.client.Open()
	if err != nil {
		return nil, err
	}
	defer ts.client.db.Close()
	query := "SELECT memstat_id,timestamp,raw_units FROM member_transactions WHERE memstat_id=$1;"
	rows, err := ts.client.db.Query(query, id)
	defer rows.Close()
	if err != nil {
		fmt.Println("It was with the Transactions Query")
		return nil, err
	}
	for rows.Next() {
		var transaction brewery.Transaction
		err = rows.Scan(&transaction.ID, &transaction.Timestamp, &transaction.RawUnits)
		if err != nil {
			fmt.Println("It happened while scanning the transaction")
			return nil, err
		}

		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (ts *TransactionService) TransactionsBytes(id int) ([]byte, error) {
	var tt []byte
	err := ts.client.Open()
	if err != nil {
		return nil, err
	}
	defer ts.client.db.Close()
	err = ts.client.db.QueryRow("select transactions from transactions where memstat_id=$1", id).Scan(&tt)
	if err != nil {
		return nil, err
	}
	return tt, nil
}
