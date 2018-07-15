package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/antschmidt/brewery-backend"
	pb "github.com/antschmidt/brewery-backend/grpc_go"
	pg "github.com/antschmidt/brewery-backend/postgres"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{}

func main() {
	cert := os.Getenv("GRPC_CRT_TSH")
	key := os.Getenv("GRPC_KEY_TSH")

	lis, err := net.Listen("tcp", fmt.Sprintf("grpctestdev.theschmidt.house:%d", 8888))
	if err != nil {
		log.Fatalf("failed to listen to brewerysrv: %v", err)
	}

	screds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	grpcsrv := grpc.NewServer(grpc.Creds(screds))
	server := &server{}

	pb.RegisterBreweryServiceServer(grpcsrv, server)
	go grpcsrv.Serve(lis)

	runHTTP("grpctestdev.theschmidt.house:8888")
}

func runHTTP(clientAddr string) {
	cert := os.Getenv("GRPC_CRT_TSH")
	key := os.Getenv("GRPC_KEY_TSH")
	//runtime.HTTPError = CustomHTTPError

	addr := "grpctestdev.theschmidt.house:6001"
	ccreds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalf("gateway cert load error: %s", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(ccreds)}
	//opts := []grpc.DialOption{grpc.WithInsecure()}

	mux := runtime.NewServeMux()
	if err := pb.RegisterBreweryServiceHandlerFromEndpoint(context.Background(), mux, clientAddr, opts); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
	log.Printf("HTTP Listening on %s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, cert, key, mux))
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

	for _, d := range acdata {
		var ac pb.AutoComplete
		ac.Membernumber = int32(d.MemberNumber)
		ac.MembershipID = int32(d.MembershipID)
		ac.AutoComplete = d.Value
		pbac.Data = append(pbac.Data, &ac)
	}
	return &pbac, nil
}

func (s *server) MemberService(ctx context.Context, m *pb.MemberRequest) (*pb.MemberResponse, error) {
	var mr pb.MemberResponse
	fmt.Println(m.Action)
	fmt.Println("Now going to try to print the m.Member")
	fmt.Printf("===============\n %v \n ====================", m.Member)
	switch m.Action {
	case 0:
		return nil, fmt.Errorf("This type of member query is not yet implemented")

	case 1:
		db := pg.NewClient()
		err := db.Open()
		if err != nil {
			mr.Status = 2
			return &mr, err
		}
		member := protoToBreweryMember(*m.Member)
		ms := db.MemberService()
		mss := db.MembershipService()
		ts := db.TransactionService()
		id, err := ms.Add(member)
		fmt.Println("New Member ID is: ", id)
		var newms brewery.Membership

		newms.StartDate, err = time.Parse("2006-01-02 15:04:05 +0000 +0000", m.Member.Memberships[0].StartDate)
		if err != nil {
			mr.Status = 2
			err = ms.RemoveByID(id)
			if err != nil {
				return &mr, fmt.Errorf("I couldn't understand that date and I was unable to remove the incomplete member with id: %v\nwith error: %v", id, err)
			}
			return &mr, fmt.Errorf("I couldn't understand that date, please try again")
		}
		newms.TotalRawUnits = m.Member.Memberships[0].TotalRawUnits
		newms.Type = m.Member.Memberships[0].Type
		msid, err := mss.Add(id, &newms)
		if err != nil {
			mr.Status = 2
			return &mr, fmt.Errorf("The Member was added successfully though the Membership was not, this will likely cause other problems. Contact Tony at tony@geometrics.io and alerts@geometrics.io\nError: %v", err)
		}
		_, err = ts.Add(msid, 0)
		if err != nil {
			mr.Status = 2
			return &mr, fmt.Errorf("The member data and membership were logged but there was a problem with the initial dummy transaction which is necessary for some reason please contact tony@geometrics.io and alerts@geometrics.io (if the system is no longer up)\nerror: %v\nid: %v", err, id)
		}
		acs, err := db.AutoComplete()
		if err != nil {
			mr.Status = 0
			return &mr, fmt.Errorf("The member was added successfully though the Autocomplete data was not updated, please refresh the page\nerror: %v", err)
		}
		mr.Member, err = s.MemberByID(ctx, &pb.MemberID{MemberID: id})
		if err != nil {
			log.Println(err)
		}

		mr.AutoComplete = breweryToProtoAC(acs).Data
		mr.Status = 1
		return &mr, nil

	case 2:
		mr.Status = 2
		return &mr, fmt.Errorf("This type of deletions is not yet implemented")

	case 3:
		mr.Status = 2
		return &mr, fmt.Errorf("This type of update not yet implemented")
	}
	mr.Status = 2
	return &mr, fmt.Errorf("There seems to be no Action declared")
}

func breweryToProtoAC(acs []pg.AutoComplete) *pb.AutoCompleteData {
	var pbac pb.AutoCompleteData
	for _, d := range acs {
		var ac pb.AutoComplete
		ac.Membernumber = int32(d.MemberNumber)
		ac.MembershipID = int32(d.MembershipID)
		ac.AutoComplete = d.Value
		pbac.Data = append(pbac.Data, &ac)
	}
	return &pbac
}

func protoToBreweryMember(pm pb.Member) *brewery.Member {
	var member brewery.Member
	fmt.Println("I made it into the protoToBreweryMember now going to try to print the &pm")
	fmt.Println("The Member to be added is \n", pm)
	fmt.Println("The m.Id is: ", pm.Id)
	member.Number = int(pm.Membernumber)
	member.Id = pm.Id
	member.Names = protoNamesToByte(pm.Names)
	member.Number = int(pm.Membernumber)
	member.Contact = protoContactToBytes(pm.Contact)
	return &member
}

func protoNamesToByte(pnames []*pb.Name) []byte {
	var names []map[string]string
	for _, n := range pnames {
		name := make(map[string]string)
		name["firstname"] = n.First
		name["lastname"] = n.Last
		names = append(names, name)
	}
	bytes, err := json.Marshal(names)
	if err != nil {
		return []byte{}
	}
	return bytes
}

func bytesToProtoName(b []byte) []*pb.Name {
	var names []map[string]string
	fmt.Println("Inside b2pn: ", string(b))

	err := json.Unmarshal(b, &names)
	if err != nil {
		fmt.Println("Couldn't unmarshal the bytesToProtoName: ", err)
		return nil
	}

	var protoNames []*pb.Name
	for _, n := range names {
		var protoName pb.Name
		protoName.First = n["firstname"]
		protoName.Last = n["lastname"]
		protoNames = append(protoNames, &protoName)
	}
	return protoNames

}

func bytesToProtoContact(b []byte) pb.Contact {
	protoContact := pb.Contact{}
	var contact map[string]string

	err := json.Unmarshal(b, &contact)
	if err != nil {
		log.Printf("Something went wrong with the contact info: %v\n\n But let's keep the beer flowing", err)
		return protoContact
	}

	protoContact.City = contact["city"]
	protoContact.Email = contact["email"]
	protoContact.Phone = contact["phone"]
	protoContact.State = contact["state"]
	protoContact.Street = contact["street"]
	protoContact.Zip = contact["zip"]

	return protoContact
}

func protoContactToBytes(pc *pb.Contact) []byte {
	// fmt.Println("The proto contact data coming into the protoC2B is: ", pc)
	// var contact struct {
	// 	zip    string
	// 	city   string
	// 	email  string
	// 	phone  string
	// 	state  string
	// 	street string
	// }
	// var contact map[string]string
	// contact.city = pc.City
	// contact.zip = pc.Zip
	// contact.email = pc.Email
	// contact.phone = pc.Phone
	// contact.state = pc.State
	// contact.street = pc.Street
	// fmt.Println("Struct to be marshalled into a byte slice is: ", contact)
	contactBytes, err := json.Marshal(pc)
	if err != nil {
		log.Println("Something went wrong with the contact marshaling, oh well, here's blank contact info.. keep the beer flowing!")
		return []byte("{\"zip\": \" \", \"city\": \" \", \"email\": \" \", \"phone\": \" \", \"state\": \" \", \"street\": \" \"}")
	}
	fmt.Println("Contact bytes to string within protoContactToBytes is: ", string(contactBytes))
	return contactBytes
}

func (s *server) DeleteTransaction(ctx context.Context, t *pb.TransactionRequest) ([]*pb.Transaction, error) {
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	tserv := db.TransactionService()

	btrans, err := protoToBreweryTransaction(int(t.MembershipID), t.Transaction)

	err = tserv.Remove(btrans)
	if err != nil {
		return nil, err
	}
	trans, err := tserv.Transactions(int(t.MembershipID))
	if err != nil {
		return nil, err
	}

	return marshalTransactions(trans), nil

}

func (s *server) StoreTransaction(ctx context.Context, t *pb.TransactionRequest) (*pb.TransactionResponse, error) {
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return nil, err
	}
	trans := db.TransactionService()
	var tr pb.TransactionResponse
	switch t.Action {
	case 0:
		return nil, nil
	case 1:
		transaction, err := trans.Add(int(t.MembershipID), t.Transaction.RawUnits)
		if err != nil {
			return nil, err
		}
		var ptrans pb.Transaction
		ptrans.RawUnits = transaction.RawUnits
		ptrans.Timestamp = transaction.Timestamp.String()
		tr.Transaction = &ptrans
		tr.Status = 1
		return &tr, nil
	case 2:
		btrans, err := protoToBreweryTransaction(int(t.MembershipID), t.Transaction)
		if err != nil {
			tr.Status = 2
			return &tr, err
		}
		err = trans.Remove(btrans)
		if err != nil {
			tr.Status = 2
			return &tr, err
		}
		tr.Status = 1
		return &tr, nil
	case 3:
		btrans, err := protoToBreweryTransaction(int(t.MembershipID), t.Transaction)
		if err != nil {
			return nil, err
		}
		err = trans.Update(btrans)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, fmt.Errorf("you seem to be missing an Action")

}

func (s *server) MemberByID(ctx context.Context, mid *pb.MemberID) (*pb.Member, error) {
	id := mid.MemberID
	pm, err := memberToProtoByID(id)
	if err != nil {
		log.Println("memberToProtoByID function is at fault:", err)
		return nil, err
	}
	return &pm, err
}

func (s *server) MemberByNumber(ctx context.Context, n *pb.Membernumber) (*pb.Member, error) {
	number := int(n.Membernumber)
	pm, err := memberToProtoByNumber(number)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func memberToProtoByNumber(n int) (pb.Member, error) {
	var pMember pb.Member
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return pMember, err
	}
	ms := db.MemberService()
	mss := db.MembershipService()
	ts := db.TransactionService()
	bMember, err := ms.MemberByNumber(n)

	contact := bytesToProtoContact(bMember.Contact)

	pMember.Contact = &contact

	protoName := bytesToProtoName(bMember.Names)
	pMember.Id = bMember.Id
	pMember.Names = protoName
	bMemberships, err := mss.MembershipsByID(bMember.Id)
	if err != nil {
		return pMember, err
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
			return pMember, err
		}
		pMembership.Transactions = marshalTransactions(membershipTransactions)
		pMember.Memberships = append(pMember.Memberships, &pMembership)
	}
	fmt.Println(pMember)
	return pMember, nil

}

func memberToProtoByID(id string) (pb.Member, error) {
	var pMember pb.Member
	db := pg.NewClient()
	err := db.Open()
	if err != nil {
		return pMember, err
	}
	ms := db.MemberService()
	mss := db.MembershipService()
	ts := db.TransactionService()
	bMember, err := ms.MemberByID(id)

	contact := bytesToProtoContact(bMember.Contact)

	pMember.Contact = &contact

	protoName := bytesToProtoName(bMember.Names)
	pMember.Id = id
	pMember.Names = protoName
	bMemberships, err := mss.MembershipsByID(id)
	if err != nil {
		return pMember, err
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
			return pMember, err
		}
		pMembership.Transactions = marshalTransactions(membershipTransactions)
		pMember.Memberships = append(pMember.Memberships, &pMembership)
	}
	return pMember, nil

}

func marshalTransactions(btr []brewery.Transaction) []*pb.Transaction {
	var pbts []*pb.Transaction
	for _, v := range btr {
		var pbt pb.Transaction
		pbt.RawUnits = v.RawUnits
		pbt.Timestamp = v.Timestamp.String()
		pbts = append(pbts, &pbt)
	}
	return pbts
}

func protoToBreweryTransaction(id int, t *pb.Transaction) (*brewery.Transaction, error) {
	var bt brewery.Transaction
	bt.ID = id
	bt.RawUnits = t.RawUnits
	transtime, err := time.Parse("2006-01-02 15:04:05 +0000 +0000", t.Timestamp)
	if err != nil {
		transtime = time.Now()
	}
	bt.Timestamp = transtime
	return &bt, nil
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
	var pbt pb.Transactions
	pbt.Transactions = marshalTransactions(transactions)
	return &pbt, nil
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

//curl -d '{"Action": 1, "MemberData":{"Id":"","Names":[{"First":"The TESTY","Last":"McTester"},{"First":"Testi","Last":"McTester"}],"Contact":{"zip":"62606","city":"BLUBADUB","email":"blubaduber@gmail.com","phone":"111-111-1111","state":"Conscious","street":"1820 For Now"},"Memberships":[{"MembershipID":0,"Type":"social","StartDate":"2018-07-13 00:00:00 +0000 +0000","TotalRawUnits":50,"Active":true}]}}' -H "Content-Type: application/json" -X POST https://grpctestdev.theschmidt.house/v1/brewery/member
