package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/antschmidt/brewery-backend/grpc_go"
	"google.golang.org/grpc"
)

var _ pb.Empty

func main() {
	var conn *grpc.ClientConn

	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not conn to grpc: %v", err)
	}
	defer conn.Close()
	if conn != nil {
		fmt.Println("It seems the conneciton was made")
	}
	b := pb.NewBreweryServiceClient(conn)
	if b != nil {
		fmt.Println("Got the PB service client", b)
	}
	data, err := b.AutoCompleteRequest(context.Background(), &pb.Empty{})
	fmt.Println("Data is: ", data)
	if err != nil {
		log.Fatalf("I was unable to pull the data: %v", err)
	}
	for _, d := range data.Data {
		fmt.Println(d)
	}
}
