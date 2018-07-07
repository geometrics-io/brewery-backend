package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/antschmidt/brewery-backend"
	pb "github.com/antschmidt/brewery-backend/grpc_go"
	pg "github.com/antschmidt/brewery-backend/postgres"
	"google.golang.org/grpc"
)

type server struct{}

func main() {

	var _ pb.Empty
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8888))
	if err != nil {
		log.Fatalf("failed to listen to brewerysrv: %v", err)
	}

	grpcsrv := grpc.NewServer()
	server := &server{}

	pb.RegisterBreweryServiceServer(grpcsrv, server)
	if err := grpcsrv.Serve(lis); err != nil {
		log.Fatalf("Couldn't serve it up: %v", err)
	}
}

func (s *server) AutoCompleteRequest(ctx context.Context, empty *pb.Empty) (*pb.AutoCompleteData, error) {
	var pbac pb.AutoCompleteData
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		log.Fatalf("failed to open the db connection: %v", err)
	}
	acdata, err := db.AutoComplete()
	if err != nil {
		log.Fatalf("Failed to grab autocomplete data from postgres: %v", err)
	}
	fmt.Println(acdata)

	for _, d := range acdata {
		var ac pb.AutoComplete
		ac.Membernumber = int32(d.MemberNumber)
		ac.MembershipID = int32(d.MembershipID)
		ac.AutoComplete = []byte(d.Value)
		pbac.Data = append(pbac.Data, &ac)
	}
	fmt.Println(pbac)
	return &pbac, nil
}

func (s *server) NewMember(ctx context.Context, m *pb.NewMemberData) (*pb.AutoComplete, error) {
	var ac *pb.AutoComplete
	db := pg.NewClient()
	//defer db.Close()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	var newmem *brewery.Member
	newmem.Id = m.Member.Id
	newmem.Names = m.Member.Names
	newmem.Number = int(m.Member.Membernumber)
	newmem.Contact = m.Member.Contact.Contact
	ms := db.MemberService()
	err = ms.Add(newmem)
	if err != nil {
		return nil, err
	}
	acs, err := db.AutoComplete()
	if err != nil {
		return nil, err
	}
	for _, v := range acs {
		if v.MemberNumber == int(m.Member.Membernumber) {
			ac.AutoComplete = []byte(v.Value)
			ac.Membernumber = int32(v.MemberNumber)
			ac.MembershipID = int32(v.MembershipID)
		}
	}
	return ac, nil
}

func (s *server) DeleteTransaction(ctx context.Context, str *pb.StoreTransactionRequest) (*pb.Transactions, error) {
	var btrans brewery.Transaction
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339, str.Transaction.Timestamp)
	if err != nil {
		return nil, err
	}
	tserv := db.TransactionService()
	btrans.ID = int(str.MembershipID)
	btrans.RawUnits = str.Transaction.RawUnits
	btrans.Timestamp = t
	err = tserv.Remove(&btrans)
	if err != nil {
		return nil, err
	}
	trans, err := tserv.Transactions(int(str.MembershipID))
	if err != nil {
		return nil, err
	}

	return marshalTransactions(trans), nil

}

func (s *server) MemberByID(ctx context.Context, mid *pb.MemberID) (*pb.Member, error) {
	id := mid.MemberID
	return memberToProtoByID(id)
}

func memberToProtoByID(id string) (*pb.Member, error) {
	var pMember *pb.Member
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	ms := db.MemberService()
	mss := db.MembershipService()
	ts := db.TransactionService()
	bMember, err := ms.MemberByID(id)
	pMember.Contact.Contact = bMember.Contact
	pMember.Id = bMember.Id
	pMember.Names = bMember.Names
	bMemberships, err := mss.MembershipsByID(id)
	if err != nil {
		return nil, err
	}
	for _, membership := range bMemberships {
		var pMembership pb.Membership
		pMembership.MembershipID = int32(membership.ID)
		pMembership.Active = membership.Active
		pMembership.Type = membership.Type
		pMembership.TotalRawUnits = membership.TotalRawUnits
		pMembership.StartDate = membership.StartDate.String()
		//pMembership.UnitBase
		membershipTransactions, err := ts.Transactions(membership.ID)
		if err != nil {
			return nil, err
		}
		pMembership.Transactions = marshalTransactions(membershipTransactions)
		pMember.Memberships = append(pMember.Memberships, &pMembership)
	}
	return pMember, nil

}

func marshalTransactions(btr []brewery.Transaction) *pb.Transactions {
	var pbts pb.Transactions
	for _, v := range btr {
		var pbt pb.Transaction
		pbt.RawUnits = v.RawUnits
		pbt.Timestamp = v.Timestamp.String()
		pbts.Transactions = append(pbts.Transactions, &pbt)
	}
	return &pbts
}

func (s *server) TransactionsByID(ctx context.Context, id *pb.MembershipID) (*pb.Transactions, error) {
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	ts := db.TransactionService()
	transactions, err := ts.Transactions(int(id.MembershipID))
	if err != nil {
		return nil, err
	}
	return marshalTransactions(transactions), nil
}

func (s *server) MembershipsByID(ctx context.Context, id *pb.MemberID) (*pb.Memberships, error) {
	var pMemberships pb.Memberships
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	mss := db.MembershipService()
	memberships, err := mss.MembershipsByID(id.MemberID)
	if err != nil {
		return nil, err
	}
	ts := db.TransactionService()
	for _, membership := range memberships {
		var pMembership pb.Membership
		pMembership.MembershipID = int32(membership.ID)
		pMembership.Active = membership.Active
		pMembership.Type = membership.Type
		pMembership.TotalRawUnits = membership.TotalRawUnits
		pMembership.StartDate = membership.StartDate.String()
		//pMembership.UnitBase
		membershipTransactions, err := ts.Transactions(membership.ID)
		if err != nil {
			return nil, err
		}
		pMembership.Transactions = marshalTransactions(membershipTransactions)
		pMemberships.Memberships = append(pMemberships.Memberships, &pMembership)
	}
	return &pMemberships, nil
}

// AutoCompleteRequest(context.Context, *Empty) (*AutoCompleteData, error)
// NewMember(context.Context, *NewMemberData) (*AutoComplete, error)
// MemberByID(context.Context, *MemberID) (*Member, error)
// MembershipsByID(context.Context, *MemberID) (*Memberships, error)
// TransactionsByID(context.Context, *MembershipID) (*Transactions, error)
// StoreTransaction(context.Context, *StoreTransactionRequest) (*Transactions, error)
// DeleteTransaction(context.Context, *StoreTransactionRequest) (*Transactions, error)
