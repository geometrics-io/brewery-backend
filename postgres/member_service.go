package postgres

import (
	"fmt"
	"log"

	brewery "github.com/antschmidt/brewery-backend"
)

// MemberService contains the client which implements the brewery.MemberStorage interface
type MemberService struct {
	client *Client
}

// MemberByID pulls a members basic data with the member's id string
func (ms *MemberService) MemberByID(id string) (brewery.Member, error) {
	var member brewery.Member
	err := ms.client.Open()
	if err != nil {
		//fix this to handle this error with some sort of message concerning database connectivity and who to contact
		log.Fatal(err)
	}
	defer ms.client.db.Close()
	query := "select members.id,members.membernumber,members.names,contact.contact from members left join contact on contact.id = members.id where members.id=$1;"
	err = ms.client.db.QueryRow(query, id).Scan(&member.Id, &member.Number, &member.Names, &member.Contact)
	if err != nil {
		return member, err
	}
	fmt.Println(member)
	return member, nil
}

// MemberByNumber returns the member to which the provided member number belongs.
func (ms *MemberService) MemberByNumber(n int) (brewery.Member, error) {
	var member brewery.Member
	err := ms.client.Open()
	if err != nil {
		//fix this to handle this error with some sort of message concerning database connectivity and who to contact
		log.Fatal(err)
	}
	defer ms.client.db.Close()
	query := "select members.id,members.membernumber,members.names,contact.contact from members left join contact on contact.id = members.id where members.membernumber=$1;"
	err = ms.client.db.QueryRow(query, n).Scan(&member.Id, &member.Number, &member.Names, &member.Contact)
	if err != nil {
		return member, err
	}
	return member, nil
}

// Add returns a MemberID after inserting a member into the database
func (ms *MemberService) Add(m *brewery.Member) (string, error) {
	err := ms.client.Open()
	if err != nil {
		//fix this to handle this error with some sort of message concerning database connectivity and who to contact
		log.Fatal(err)
	}
	defer ms.client.db.Close()
	query := "INSERT INTO members (membernumber,names) values ($1, $2) returning id;"
	err = ms.client.db.QueryRow(query, m.Number, m.Names).Scan(&m.Id)
	if err != nil {
		return "", err
	}
	fmt.Println("Contact info is: ", string(m.Contact))
	query = "INSERT INTO contact (id, contact) values ($1, $2);"
	_, err = ms.client.db.Exec(query, m.Id, string(m.Contact))
	if err != nil {
		return "", err
	}
	return m.Id, nil
}

// RemoveByID removes a member and all releated data of that member from the database
func (ms *MemberService) RemoveByID(id string) error {
	err := ms.client.Open()
	if err != nil {
		return err
	}
	defer ms.client.db.Close()
	query := "delete from contact where id=$1;"
	_, err = ms.client.db.Exec(query, id)
	if err != nil {
		log.Printf("%v Contact nor Member deleted for id: %v", ms.client.Now(), id)
		return err
	}
	query = "delete from members where id=$1;"
	_, err = ms.client.db.Exec(query, id)
	if err != nil {
		log.Printf("%v Contact deleted but not member with id: %v", ms.client.Now(), id)
		return err
	}

	return nil
}

// Update will update the database entry of the member provided
func (ms *MemberService) Update(m *brewery.Member) error {
	err := ms.client.Open()
	if err != nil {
		return err
	}
	defer ms.client.db.Close()
	query := "update members set membernumber=$1 where id=$2"
	_, err = ms.client.db.Exec(query, m.Number, m.Id)
	if err != nil {
		log.Printf("%v Member nor Contact updated for id: %v", ms.client.Now(), m.Id)
		return err
	}
	query = "update contact set contact=$1 where id=$2"
	_, err = ms.client.db.Exec(query, m.Contact, m.Id)
	if err != nil {
		log.Printf("%v Member updated but not contact for id: %v", ms.client.Now(), m.Id)
	}
	return nil
}
