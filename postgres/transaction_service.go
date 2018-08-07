package postgres

import (
	"fmt"
	"log"

	"github.com/lib/pq"

	brewery "github.com/antschmidt/brewery-backend"
)

// TransactionService implements the brewery.TransactionStorage interface
type TransactionService struct {
	client *Client
}

// Add creates a transaction in the database for the provided membership with the given units and returns the transaction data with the timestamp
func (ts *TransactionService) Add(id int, units float64) (*brewery.Transaction, error) {
	var transaction brewery.Transaction
	err := ts.client.Open()
	if err != nil {
		return nil, err
	}
	defer ts.client.db.Close()
	query := "INSERT INTO member_transactions (memstat_id, raw_units) VALUES ($1, $2) returning timestamp;"
	err = ts.client.db.QueryRow(query, id, units).Scan(&transaction.Timestamp)
	if err != nil {
		log.Printf("Failed to insert transaction or grab the timestamp from said transaction from\nID: %v\nRawUnits: %v\n", id, units)
		return nil, err
	}
	transaction.ID = id
	transaction.RawUnits = units
	return &transaction, nil

}

// Remove deletes a given transaction from the database
func (ts *TransactionService) Remove(t *brewery.Transaction) error {
	err := ts.client.Open()
	if err != nil {
		return err
	}
	defer ts.client.db.Close()
	query := "DELETE FROM member_transactions where memstat_id=$1 and timestamp=$2 and raw_units=$3;"
	res, err := ts.client.db.Exec(query, t.ID, pq.FormatTimestamp(t.Timestamp), t.RawUnits)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

// Update allows you to alter the units of particular transaction in the database when given the id and timestamp of an existing transaction
func (ts *TransactionService) Update(t *brewery.Transaction) error {
	err := ts.client.Open()
	if err != nil {
		return err
	}
	defer ts.client.db.Close()
	query := "update member_transactions set raw_units=$1 where memstat_id=$2 and timestamp=$3;"
	_, err = ts.client.db.Exec(query, t.RawUnits, t.ID, t.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// Transactions returns the transactions for the membership of the privided membership id
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

// TransactionsBytes returns teh transactions of the membership connected to the given id and returns the transaction data as a slice of byte
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
