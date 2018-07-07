package postgres

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/antschmidt/brewery-backend"
	_ "github.com/lib/pq"
)

type Client struct {
	Now                func() time.Time
	memberService      MemberService
	membershipService  MembershipService
	transactionService TransactionService
	db                 *sql.DB
}

type AutoComplete struct {
	ID           string
	MemberNumber int
	MembershipID int
	Value        string
}

func NewClient() *Client {
	c := &Client{Now: time.Now}
	c.memberService.client = c
	c.membershipService.client = c
	c.transactionService.client = c
	return c
}

func (c *Client) Open() error {
	db, err := sql.Open("postgres", dbInfo())
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

func dbInfo() string {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASS"), os.Getenv("POSTGRES_DB"))
	return dbinfo
}

func (c *Client) AllJson() ([]byte, error) {
	var jsonblob bytes.Buffer
	err := c.Open()
	if err != nil {
		return nil, err
	}
	defer c.db.Close()
	rows, err := c.db.Query("select data from alljson;")
	if err != nil {
		fmt.Println("It's the row query")
		return nil, err
	}
	jsonblob.Write([]byte("["))
	for rows.Next() {
		var blob []byte
		err = rows.Scan(&blob)
		if err != nil {
			return nil, err
		}
		jsonblob.Write(blob)

		jsonblob.Write([]byte(","))

	}
	tmpblob := bytes.TrimRight(jsonblob.Bytes(), ",")
	jsonblob.Reset()
	jsonblob.Write(tmpblob)
	jsonblob.Write([]byte("]"))
	return jsonblob.Bytes(), nil
}

func (c *Client) AutoComplete() ([]AutoComplete, error) {
	var acs []AutoComplete
	err := c.Open()
	if err != nil {
		return nil, err
	}
	defer c.db.Close()
	query := "select ms.id,m.membernumber,ms.memstat_id,ms.membership,m.names from member_status ms left join  members m on ms.id = m.id where ms.active=true order by m.membernumber;"
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var ac AutoComplete
		var namesb []byte
		var names []map[string]string
		var typeb []byte
		err = rows.Scan(&ac.ID, &ac.MemberNumber, &ac.MembershipID, &typeb, &namesb)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(namesb, &names)
		if err != nil {
			return nil, err
		}
		for _, n := range names {
			ac.Value = fmt.Sprintf("%v %v %v - %v", ac.MemberNumber, n["firstname"], n["lastname"], string(typeb))
			acs = append(acs, ac)
		}

	}
	return acs, nil
}

func (c *Client) TransactionService() brewery.TransactionStorage { return &c.transactionService }
func (c *Client) MembershipService() brewery.MembershipStorage   { return &c.membershipService }
func (c *Client) MemberService() brewery.MemberStorage           { return &c.memberService }

//func (c *Client) TransactionService() brewery.TransactionStorage { return &c.transactionService }

//dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
