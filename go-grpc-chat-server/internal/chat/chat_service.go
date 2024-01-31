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

	// Check if the player and group IDs are valid.
	if playerID == 0 || groupID == 0 {
		return nil, errors.New("invalid player or group ID")
	}

	// Subscribe the player to the group.
	s.subscribePlayerToGroup(playerID, groupID)
	log.Printf("Player %d subscribed to group %d", playerID, groupID)

	// Return a successful response.
	return &SubscriptionResponse{
		Status: &ResponseStatus{Code: 200, Message: "Subscribed successfully"},
		Group:  req.Group,
	}, nil
}

// Unsubscribe removes a player from a chat group.
func (s *Server) Unsubscribe(ctx context.Context, req *UnsubscriptionRequest) (*UnsubscriptionResponse, error) {
	playerID := req.Player.Id
	groupID := req.Group.Id

	// Check if the player and group IDs are valid.
	if playerID == 0 || groupID == 0 {
		return nil, errors.New("invalid player or group ID")
	}

	// Unsubscribe the player from the group.
	s.unsubscribePlayerFromGroup(playerID, groupID)
	log.Printf("Player %d unsubscribed from group %d", playerID, groupID)

	// Return a successful response.
	return &UnsubscriptionResponse{
		Status: &ResponseStatus{Code: 200, Message: "Unsubscribed successfully"},
	}, nil
}

// subscribePlayerToGroup adds a player to a group's subscription list.
func (s *Server) subscribePlayerToGroup(playerID, groupID uint64) {
	// Lock the groupSubscriptions map for safe concurrent access.
	s.groupSubscriptions.Store(groupID, append(s.getSubscribedPlayers(groupID), playerID))
}

// unsubscribePlayerFromGroup removes a player from a group's subscription list.
func (s *Server) unsubscribePlayerFromGroup(playerID, groupID uint64) {
	subscribedPlayers := s.getSubscribedPlayers(groupID)
	for i, id := range subscribedPlayers {
		if id == playerID {
			s.groupSubscriptions.Store(groupID, append(subscribedPlayers[:i], subscribedPlayers[i+1:]...))
			break
		}
	}
}

// getSubscribedPlayers returns a slice of player IDs subscribed to a group.
func (s *Server) getSubscribedPlayers(groupID uint64) []uint64 {
	players, ok := s.groupSubscriptions.Load(groupID)
	if !ok {
		return []uint64{}
	}
	return players.([]uint64)
}

// StreamMessages handles bidirectional streaming for chat messages.
func (s *Server) StreamMessages(stream ChatService_StreamMessagesServer) error {
	// Temporary storage for client's stream
	var clientStream ChatService_StreamMessagesServer

	for {
		// Receive a message from the stream
		in, err := stream.Recv()
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			return err
		}

		// Handle different message types
		switch msg := in.MessageType.(type) {
		case *MessageStream_Message:
			// Handle incoming chat message
			log.Printf("Received message: %v", msg)
			// [Your message handling logic here]

		case *MessageStream_StreamRequest:
			// Handle stream request
			log.Printf("Stream request from player: %v", msg.StreamRequest.Player.Id)
			clientStream = stream
			s.streams.Store(msg.StreamRequest.Player.Id, stream)

		default:
			// Unknown message type
			return status.Errorf(codes.InvalidArgument, "Unknown message type")
		}

		// If clientStream is set, you can send messages back to the client
		if clientStream != nil {
			// [Your logic to send messages to the client]
		}
	}
}

// Additional methods for Subscribe, Unsubscribe, etc., can be added here.
