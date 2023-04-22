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
	chat.UnimplementedChatServer
	messages []chat.StreamMessagesResponse
	mu       sync.Mutex
}

func (s *server) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*chat.SendMessageResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := chat.StreamMessagesResponse{
		Message: req.Message,
		User:    req.User,
	}

	s.messages = append(s.messages, msg)

	return &chat.SendMessageResponse{Success: true}, nil
}

func (s *server) StreamMessages(_ *chat.StreamMessagesRequest, stream chat.Chat_StreamMessagesServer) error {
	s.mu.Lock()
	messages := s.messages
	s.mu.Unlock()

	for _, msg := range messages {
		if err := stream.Send(&msg); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chat.RegisterChatServer(grpcServer, &server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
