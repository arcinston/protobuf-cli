// cmd/client/main.go
package main

import (
	"context"
	"fmt"
	"log"

	chat "github.com/arcinston/protobuf-cli/protobufs/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatClient(conn)
	var user string
	log.Printf("Enter User name ")
	fmt.Scanln(&user)

	ctx := context.Background()

	// Start listening to messages
	stream, err := client.StreamMessages(ctx, &chat.StreamMessagesRequest{})
	if err != nil {
		log.Fatalf("failed to start streaming messages: %v", err)
	}

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				log.Println("EOF reached , restartign stream")
			}
			fmt.Printf("[%s] %s\n", res.User, res.Message)
		}
	}()

	// Send messages
	for {
		fmt.Print("Enter message: ")
		var msg string
		fmt.Scanln(&msg)
		_, err = client.SendMessage(ctx, &chat.SendMessageRequest{Message: msg, User: user})
		if err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}
