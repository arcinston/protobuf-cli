package main

import (
	"context"
	"log"
	"net"
	"sync"

	chat "github.com/arcinston/protobuf-cli/protobufs/chat"

	"google.golang.org/grpc"
)

type server struct {
	chat.ChatServer
	messages []*chat.StreamMessagesResponse
	mu       sync.Mutex
}

func (s *server) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	msg := &chat.StreamMessagesResponse{
		Message: req.Message,
		User:    req.User,
	}
	log.Println("msg: ", req.Message)
	log.Println("user: ", req.User)

	s.messages = append(s.messages, msg)

	return &chat.SendMessageResponse{Success: true}, nil
}

func (s *server) StreamMessages(_ *chat.StreamMessagesRequest, stream chat.Chat_StreamMessagesServer) error {
	s.mu.Lock()
	messages := make([]*chat.StreamMessagesResponse, len(s.messages))
	copy(messages, s.messages)
	s.mu.Unlock()

	for _, msg := range messages {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SendWelcomeMessage(stream chat.Chat_StreamMessagesServer) error {
	welcomeMsg := &chat.StreamMessagesResponse{
		Message: "Welcome to the chat!",
		User:    "Server",
	}

	if err := stream.Send(welcomeMsg); err != nil {
		return err
	}

	return nil
}

func main() {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on tcp :50051")

	chat.RegisterChatServer(grpcServer, &server{})
	log.Printf("grpcServer made , and chatserver registered")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("grpcServer served")
}
