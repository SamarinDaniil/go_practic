package server

import (
	"context"
	"fmt"
	"log"

	eventpb "../pkg/api/protobuf"
)

func makeEvent(client eventpb.EventManagerClient, senderID int64, time int64, name string) {
	req := &eventpb.MakeEventRequest{
		SenderID: senderID,
		Time:     time,
		Name:     name,
	}

	resp, err := client.MakeEvent(context.Background(), req)
	if err != nil {
		log.Fatalf("MakeEvent error: %v", err)
	}

	fmt.Printf("Created{eventId:%d}\n", resp.EventID)
}
