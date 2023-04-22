// cmd/client/main.go
package main

import (
	"context"
	"fmt"
	"log"

	"your_project_path/protobufs/chat"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatClient(conn)
	user := "YourUsername"
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
				log.Fatalf("failed to receive message: %v", err)
			}
			fmt.Printf("[%s] %s\n", res.User, res.Message)
		}
	}()

	// Send messages
	for {
		fmt.Print("Enter message: ")
		msg := ""
		fmt.Scanln(&msg)

		_, err = client.SendMessage(ctx, &chat.SendMessageRequest{Message: msg, User: user})
		if err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
}
