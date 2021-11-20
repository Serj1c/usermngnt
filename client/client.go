package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Serj1c/usermngnt/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// hardcoded for testing
	newUsers := make(map[string]int32)
	newUsers["Alice"] = 43
	newUsers["Bob"] = 30
	for name, age := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User details: Name: %s, Age: %d, Id: %d", r.GetName(), r.GetAge(), r.GetId())
	}
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nr.GetUsers: %v", r.GetUsers())
	fmt.Printf("\nr: %v", r)
}
