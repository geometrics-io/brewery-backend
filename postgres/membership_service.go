package postgres

import (
	brewery "github.com/antschmidt/brewery-backend"
)

// MembershipService implements the brewery.MembershipStorage interface
type MembershipService struct {
	client *Client
}

// Add creates a new membership based on the member_id and membership data provided
func (mss *MembershipService) Add(id string, m *brewery.Membership) (int, error) {
	err := mss.client.Open()
	if err != nil {
		return 0, err
	}
	defer mss.client.db.Close()
	query := "INSERT INTO member_status (id,start_date,membership,total_raw_units) VALUES ($1, $2, $3, $4) RETURNING memstat_id;"
	var MemstatID int
	err = mss.client.db.QueryRow(query, id, m.StartDate, m.Type, m.TotalRawUnits).Scan(&MemstatID)
	if err != nil {
		return 0, err
	}
	return MemstatID, nil
}

// Remove deletes a single specified membership from the database
func (mss *MembershipService) Remove(m *brewery.Membership) error {
	err := mss.client.Open()
	if err != nil {
		return err
	}
	defer mss.client.db.Close()
	query := "DELETE FROM member_status WHERE memstat_id=$1;"
	_, err = mss.client.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Update alters a membership's entry in the database with the data provided
func (mss *MembershipService) Update(id int, m brewery.Membership) error {
	err := mss.client.Open()
	if err != nil {
		return err
	}
	defer mss.client.db.Close()
	query := "UPDATE member_status SET start_date=$1, membership=$2, total_raw_units=$3, active=$4 WHERE memstat_id=$5;"
	_, err = mss.client.db.Exec(query, m.StartDate, m.Type, m.TotalRawUnits, m.Active, m.ID)
	if err != nil {
		return err
	}
	return nil
}

// MembershipsByID returns the memberships held by a particular member with the given ID
func (mss *MembershipService) MembershipsByID(id string) ([]*brewery.Membership, error) {
	var memberships []*brewery.Membership
	err := mss.client.Open()
	if err != nil {
		return nil, err
	}
	defer mss.client.db.Close()
	query := "SELECT memstat_id,start_date,membership,total_raw_units,active FROM member_status where id=$1;"
	rows, err := mss.client.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var membership brewery.Membership
		err := rows.Scan(&membership.ID, &membership.StartDate, &membership.Type, &membership.TotalRawUnits, &membership.Active)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, &membership)
	}
	return memberships, nil
}

// Memberships returns all memberships held in the connected database
// (I'm not sure why I created this, perhaps it should be membership by membernumber...)
func (mss *MembershipService) Memberships() ([]*brewery.Membership, error) {
	var memberships []*brewery.Membership
	err := mss.client.Open()
	if err != nil {
		return nil, err
	}
	defer mss.client.db.Close()
	query := "SELECT memstat_id,start_date,membership,total_raw_units,active FROM member_status;"
	rows, err := mss.client.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var membership brewery.Membership
		err = rows.Scan(&membership.ID, &membership.StartDate, &membership.Type, &membership.TotalRawUnits, &membership.Active)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, &membership)
	}
	return memberships, nil
}
