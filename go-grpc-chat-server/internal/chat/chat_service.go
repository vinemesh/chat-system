//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat.proto
package chat

import (
	context "context"
	"fmt"
	"log"
	"strconv"
	sync "sync"

	kp "github.com/vinemesh/go-grpc-chat-server/internal/kafka"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Server represents the gRPC server for chat service.
type Server struct {
	UnimplementedChatServiceServer
	Producer *kp.Producer
	ErrChan  chan error
	// A map to store client streams with a Mutex for synchronization.
	streams sync.Map
}

// Subscribe adds a player to a chat group.
func (s *Server) Subscribe(ctx context.Context, req *SubscriptionRequest) (*SubscriptionResponse, error) {
	playerID := req.Player.Id
	groupID := req.Group.Id
	log.Printf("Received subscription request from player %d to group %d", playerID, groupID)

	if playerID == 0 || groupID == 0 {
		err := status.Errorf(codes.NotFound, "invalid player or group ID")
		s.ErrChan <- err
		return nil, err
	}

	err := s.Producer.ProduceMessage(
		fmt.Sprintf("chat_events_%d", groupID),
		strconv.FormatUint(playerID, 10),
		"subscribe",
	)
	if err != nil {
		s.ErrChan <- err
		return nil, status.Errorf(codes.Internal, "error producing kafka message: %v", err)
	}
	log.Printf("Player %d subscribed to group %d", playerID, groupID)

	return &SubscriptionResponse{
		Status: &ResponseStatus{Code: 200, Message: "Subscribed successfully"},
		Group:  req.Group,
	}, nil
}

// Unsubscribe removes a player from a chat group.
func (s *Server) Unsubscribe(ctx context.Context, req *UnsubscriptionRequest) (*UnsubscriptionResponse, error) {
	playerID := req.Player.Id
	groupID := req.Group.Id
	log.Printf("Received unsubscription request from player %d to group %d", playerID, groupID)

	if playerID == 0 || groupID == 0 {
		err := status.Errorf(codes.NotFound, "invalid player or group ID")
		s.ErrChan <- err
		return nil, err
	}

	err := s.Producer.ProduceMessage(
		fmt.Sprintf("chat_events_%d", groupID),
		strconv.FormatUint(playerID, 10),
		"unsubscribe",
	)
	if err != nil {
		s.ErrChan <- err
		return nil, status.Errorf(codes.Internal, "error producing kafka message: %v", err)
	}
	log.Printf("Player %d unsubscribed from group %d", playerID, groupID)

	return &UnsubscriptionResponse{
		Status: &ResponseStatus{Code: 200, Message: "Unsubscribed successfully"},
	}, nil
}

// StreamMessages handles bidirectional streaming for chat messages.
func (s *Server) StreamMessages(stream ChatService_StreamMessagesServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			s.ErrChan <- err
			return status.Errorf(codes.DataLoss, "Error receiving from stream: %v", err)
		}
		go s.handleMessage(stream, in)
	}
}

// handleMessage handles the incoming message concurrently.
func (s *Server) handleMessage(stream ChatService_StreamMessagesServer, in *MessageStream) error {
	switch msg := in.MessageType.(type) {
	case *MessageStream_Message:
		// Handle incoming chat message
		playerID := msg.Message.Player.Id
		groupID := msg.Message.Group.Id
		content := msg.Message.Content
		log.Printf("Group %d | Player %d: %s", groupID, playerID, content)

		if playerID == 0 || groupID == 0 {
			err := status.Errorf(codes.NotFound, "invalid player or group ID in chat message")
			s.ErrChan <- err
			return err
		}

		err := s.Producer.ProduceMessage(
			fmt.Sprintf("chat_messages_%d", groupID),
			strconv.FormatUint(playerID, 10),
			content,
		)
		if err != nil {
			s.ErrChan <- err
			return status.Errorf(codes.Internal, "error producing kafka message: %v", err)
		}
		return nil

	case *MessageStream_StreamRequest:
		playerId := msg.StreamRequest.Player.Id
		log.Printf("Received stream request from player %d", playerId)

		// Handle stream request, store stream for later use
		s.streams.Store(playerId, stream)
		return nil

	default:
		err := status.Errorf(codes.InvalidArgument, "Unknown message type")
		s.ErrChan <- err
		return err
	}
}

// // broadcastMessageToGroup sends a given message to all subscribed players of a group.
// func (s *Server) broadcastMessageToGroup(message *Message) {
// 	groupID := message.Group.Id
// 	players, ok := s.groupSubscriptions.Load(groupID)
// 	if !ok {
// 		log.Printf("No players subscribed to group %d", groupID)
// 		return
// 	}

// 	subscribedPlayers, ok := players.([]uint64)
// 	if !ok {
// 		log.Printf("Invalid type for group subscribers")
// 		return
// 	}

// 	for _, playerID := range subscribedPlayers {
// 		stream, ok := s.streams.Load(playerID)
// 		if !ok {
// 			log.Printf("Stream for player %d not found", playerID)
// 			continue
// 		}

// 		playerStream, ok := stream.(ChatService_StreamMessagesServer)
// 		if !ok {
// 			log.Printf("Invalid stream type for player %d", playerID)
// 			continue
// 		}

// 		if err := playerStream.Send(&MessageStream{MessageType: &MessageStream_Message{Message: message}}); err != nil {
// 			log.Printf("Error sending message to player %d: %v", playerID, err)
// 		}
// 	}
// }
