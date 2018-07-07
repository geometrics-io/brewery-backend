package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/antschmidt/brewery-backend/postgres"
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
	// jsonblob, err := transclient.AllJson()
	// if err != nil {
	// 	fmt.Println("Fuck")
	// }
	// fmt.Println(string(jsonblob))
	ac, err := transclient.AutoComplete()
	if err != nil {
		fmt.Println("Fuckin' AC")
	}

	acbs, err := json.Marshal(ac)
	fmt.Println(string(acbs))
	tclient := transclient.TransactionService()
	transactions, err := tclient.TransactionsBytes(28)
	if err != nil {
		log.Println("F-ing Transactions", err)
	}
	//ts, err := json.Marshal(transacitons)

	fmt.Println(string(transactions))
	//fmt.Println(member.Number)
}
