//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative chat.proto
package chat

import (
	context "context"
	"errors"
	"log"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server represents the gRPC server for chat service.
type Server struct {
	UnimplementedChatServiceServer
	// A map to store client streams with a Mutex for synchronization.
	streams sync.Map
	// A map to keep track of group subscriptions.
	groupSubscriptions sync.Map // map[uint64][]uint64 : groupID -> []playerIDs
}

// Subscribe adds a player to a chat group.
func (s *Server) Subscribe(ctx context.Context, req *SubscriptionRequest) (*SubscriptionResponse, error) {
	playerID := req.Player.Id
	groupID := req.Group.Id

	if playerID == 0 || groupID == 0 {
		return nil, errors.New("invalid player or group ID")
	}

	s.subscribePlayerToGroup(playerID, groupID)
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

	if playerID == 0 || groupID == 0 {
		return nil, errors.New("invalid player or group ID")
	}

	s.unsubscribePlayerFromGroup(playerID, groupID)
	log.Printf("Player %d unsubscribed from group %d", playerID, groupID)

	return &UnsubscriptionResponse{
		Status: &ResponseStatus{Code: 200, Message: "Unsubscribed successfully"},
	}, nil
}

// subscribePlayerToGroup adds a player to a group's subscription list.
func (s *Server) subscribePlayerToGroup(playerID, groupID uint64) {
	var subscribers []uint64
	if existing, ok := s.groupSubscriptions.Load(groupID); ok {
		subscribers = existing.([]uint64)
	}
	subscribers = append(subscribers, playerID)
	s.groupSubscriptions.Store(groupID, subscribers)
}

// unsubscribePlayerFromGroup removes a player from a group's subscription list.
func (s *Server) unsubscribePlayerFromGroup(playerID, groupID uint64) {
	existing, ok := s.groupSubscriptions.Load(groupID)
	if !ok {
		return // Group not found or no subscribers
	}
	subscribers := existing.([]uint64)
	for i, id := range subscribers {
		if id == playerID {
			s.groupSubscriptions.Store(groupID, append(subscribers[:i], subscribers[i+1:]...))
			break
		}
	}
}

// StreamMessages handles bidirectional streaming for chat messages.
func (s *Server) StreamMessages(stream ChatService_StreamMessagesServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			return err
		}

		switch msg := in.MessageType.(type) {
		case *MessageStream_Message:
			// Handle incoming chat message
			log.Printf("Received message from player %d in group %d", msg.Message.Player.Id, msg.Message.Group.Id)
			s.broadcastMessageToGroup(msg.Message)

		case *MessageStream_StreamRequest:
			// Handle stream request, store stream for later use
			playerID := msg.StreamRequest.Player.Id
			s.streams.Store(playerID, stream)

		default:
			return status.Errorf(codes.InvalidArgument, "Unknown message type")
		}
	}
}

// broadcastMessageToGroup sends a given message to all subscribed players of a group.
func (s *Server) broadcastMessageToGroup(message *Message) {
	groupID := message.Group.Id
	players, ok := s.groupSubscriptions.Load(groupID)
	if !ok {
		log.Printf("No players subscribed to group %d", groupID)
		return
	}

	subscribedPlayers, ok := players.([]uint64)
	if !ok {
		log.Printf("Invalid type for group subscribers")
		return
	}

	for _, playerID := range subscribedPlayers {
		stream, ok := s.streams.Load(playerID)
		if !ok {
			log.Printf("Stream for player %d not found", playerID)
			continue
		}

		playerStream, ok := stream.(ChatService_StreamMessagesServer)
		if !ok {
			log.Printf("Invalid stream type for player %d", playerID)
			continue
		}

		if err := playerStream.Send(&MessageStream{MessageType: &MessageStream_Message{Message: message}}); err != nil {
			log.Printf("Error sending message to player %d: %v", playerID, err)
		}
	}
}
