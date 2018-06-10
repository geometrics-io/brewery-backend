package main

import (
	"fmt"

	"github.com/brewery-grpc/postgres"
)

func main() {
	transclient := postgres.NewClient()
	// transservice := transclient.TransactionService()
	// transactions, err := transservice.Transactions(28)
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(transactions)

	//members := transclient.MemberService()
	//member, err := members.MemberByID("7d6c6e2a-ab5c-450e-949d-adbb4a8140f8")
	// if err != nil {
	// 	log.Println(err)
	// }
	jsonblob, err := transclient.AllJson()
	if err != nil {
		fmt.Println("Fuck")
	}
	fmt.Println(string(jsonblob))
	//fmt.Println(member.Number)
}
