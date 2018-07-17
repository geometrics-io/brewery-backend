package postgres

import (
	"github.com/antschmidt/brewery-backend"
)

type MembershipLevelService struct {
	client *Client
}

// Add returns the available membership levels after adding the provided one to the database
func (mls *MembershipLevelService) Add(ml brewery.MembershipLevel) ([]*brewery.MembershipLevel, error) {
	err := mls.client.Open()
	if err != nil {
		return nil, err
	}

	defer mls.client.db.Close()

	query := "INSERT INTO membership_levels (name,unit_type,units,unit_base) VALUES ($1, $2, $3, $4);"
	_, err = mls.client.db.Exec(query, ml.Name, ml.UnitType, ml.Units, ml.UnitBase)
	if err != nil {
		return nil, err
	}

	return mls.MembershipLevels()
}

// Remove returns remaining available membership levels after removing one from the database
func (mls *MembershipLevelService) Remove(ml brewery.MembershipLevel) ([]*brewery.MembershipLevel, error) {
	err := mls.client.Open()
	if err != nil {
		return nil, err
	}
	defer mls.client.db.Close()

	query := "DELETE FROM membership_level where name=$1;"
	_, err = mls.client.db.Exec(query, ml.Name)
	if err != nil {
		return nil, err
	}
	return mls.MembershipLevels()
}

// Update returns a membership level after updating it in the database
func (mls *MembershipLevelService) Update(name string, ml brewery.MembershipLevel) (*brewery.MembershipLevel, error) {
	err := mls.client.Open()
	if err != nil {
		return nil, err
	}
	defer mls.client.db.Close()

	query := "UPDATE membership_levels set name=$1,unit_type=$2,units=$3,unit_base=4 where name=$5;"
	_, err = mls.client.db.Exec(query, ml.Name, ml.UnitType, ml.Units, ml.UnitBase, name)
	if err != nil {
		return nil, err
	}

	return mls.MembershipLevel(ml.Name)
}

// MembershipLevels returns an array of the available membership levels
func (mls *MembershipLevelService) MembershipLevels() ([]*brewery.MembershipLevel, error) {
	err := mls.client.Open()
	if err != nil {
		return nil, err
	}
	defer mls.client.db.Close()

	query := "select name, unit_type, units, unit_base from membership_levels;"
	var levels []*brewery.MembershipLevel

	rows, err := mls.client.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var level brewery.MembershipLevel
		err := rows.Scan(&level.Name, &level.UnitType, &level.Units, &level.UnitBase)
		if err != nil {
			return nil, err
		}
		levels = append(levels, &level)
	}
	return levels, nil
}

// MembershipLevel returns the membership level data for the provided name
func (mls *MembershipLevelService) MembershipLevel(n string) (*brewery.MembershipLevel, error) {
	err := mls.client.Open()
	if err != nil {
		return nil, err
	}
	defer mls.client.db.Close()

	var level brewery.MembershipLevel
	query := "SELECT name, unit_type, units, unit_base from memebrship_levels where name=$1;"
	err = mls.client.db.QueryRow(query, n).Scan(&level.Name, &level.UnitType, &level.Units, &level.UnitBase)
	if err != nil {
		return nil, err
	}
	return &level, nil
}
