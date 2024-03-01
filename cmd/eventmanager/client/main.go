package main

import (
	eventpb "event"
	"log"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := eventpb.NewEventManagerClient(conn)
	var id int64 = 1
	var time int64 = 132312
	var name string = "Party"
	makeEvent(client, id, time, name)
}
