package brewery

import "time"

// A Membership represents a member's status as a member for a given period of time and how many units are available to the given membership.
type Membership struct {
	ID            int       `db:"memstat_id"`
	StartDate     time.Time `db:"start_date"`
	Type          string    `db:"membership"`
	TotalRawUnits float64   `db:"total_raw_units"`
	Active        bool      `db:"active"`
}

// A MembershipLevel represents a particular membership type's limits with the UnitBase being how many pints are included with each transaction.
type MembershipLevel struct {
	Name     string
	UnitType string
	Units    float64
	UnitBase int
}

// A Transaction represents a beer transaction for a customers membership.
type Transaction struct {
	ID        int
	Timestamp time.Time
	RawUnits  float64
}

// The TransactionStorage interface is for communicating with a form of storage to retain the Transaction data
type TransactionStorage interface {
	Add(id int, units float64) (*Transaction, error)
	Remove(t *Transaction) error
	Update(t *Transaction) error
	Transactions(id int) ([]Transaction, error)
	TransactionsBytes(id int) ([]byte, error)
}

// MembershipStorage is for communicating with a form of storage to retain a member's membership data
type MembershipStorage interface {
	Add(id string, m *Membership) (int, error)
	Remove(m *Membership) error
	Update(id int, m Membership) error
	MembershipsByID(id string) ([]*Membership, error)
	Memberships() ([]*Membership, error)
}

// MembershipLevelStorage is for communicating with a form of storage to retain Membership Levels
type MembershipLevelStorage interface {
	Add(ml MembershipLevel) ([]*MembershipLevel, error)
	Remove(ml MembershipLevel) ([]*MembershipLevel, error)
	Update(name string, ml MembershipLevel) (*MembershipLevel, error)
	MembershipLevels() ([]*MembershipLevel, error)
	MembershipLevel(n string) (*MembershipLevel, error)
}
