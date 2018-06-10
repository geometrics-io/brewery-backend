package postgres

type StorageService struct {
	client *Client
}

// import (
// 	"log"

// 	brewery "github.com/brewery-grpc"
// )

// type StorageService struct {
// 	client *Client
// }

// // MemberByID pulls a members basic data with the member's id string
// func (c *Client) MemberByID(id string) (*brewery.Member, error) {
// 	var member brewery.Member
// 	err := c.Open()
// 	if err != nil {
// 		//fix this to handle this error with some sort of message concerning database connectivity and who to contact
// 		log.Fatal(err)
// 	}
// 	defer c.db.Close()
// 	query := "select members.id,members.membernumber,members.names,contact.contact from members left join contact on contact.id = members.id where members.id=$1;"
// 	err = c.db.QueryRow(query, id).Scan(&member.Id, &member.Number, &member.Names, &member.Contact)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &member, nil
// }

// // Add Adds a
// func (c *Client) Add(m *brewery.Member) error {
// 	err := c.Open()
// 	if err != nil {
// 		//fix this to handle this error with some sort of message concerning database connectivity and who to contact
// 		log.Fatal(err)
// 	}
// 	defer c.db.Close()
// 	query := "INSERT INTO members (membernumber,names) values ($1, $2) returning id;"
// 	err = c.db.QueryRow(query, m.Number, m.Names).Scan(&m.Id)
// 	if err != nil {
// 		return err
// 	}
// 	query = "INSERT INTO contact (id, contact) values ($1, $2);"
// 	_, err = c.db.Exec(query, m.Id, m.Contact)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *Client) RemoveById(id string) error {
// 	err := c.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer c.db.Close()
// 	query := "delete from contact where id=$1;"
// 	_, err = c.db.Exec(query, id)
// 	if err != nil {
// 		log.Printf("%v Contact nor Member deleted for id: %v", c.Now(), id)
// 		return err
// 	}
// 	query = "delete from members where id=$1;"
// 	_, err = c.db.Exec(query, id)
// 	if err != nil {
// 		log.Printf("%v Contact deleted but not member with id: %v", c.Now(), id)
// 		return err
// 	}

// 	return nil
// }

// func (c *Client) Update(m *brewery.Member) error {
// 	err := c.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer c.db.Close()
// 	query := "update members set membernumber=$1 where id=$2"
// 	_, err = c.db.Exec(query, m.Number, m.Id)
// 	if err != nil {
// 		log.Printf("%v Member nor Contact updated for id: %v", c.Now(), m.Id)
// 		return err
// 	}
// 	query = "update contact set contact=$1 where id=$2"
// 	_, err = c.db.Exec(query, m.Contact, m.Id)
// 	if err != nil {
// 		log.Printf("%v Member updated but not contact for id: %v", c.Now(), m.Id)
// 	}
// 	return nil
// }

// func (c *Client) MemberService() brewery.MemberService {
// 	return &c.memberService
// }
