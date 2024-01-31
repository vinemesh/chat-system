//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative messages.proto

package messages

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	sync "sync"

	"google.golang.org/grpc/peer"
)

type Server struct {
	MessageServiceServer
	// Mapa para armazenar streams de clientes com um Mutex para sincronização
	streams sync.Map
}

func (s *Server) SendMessage(ctx context.Context, message *PlayerMessage) (*ResponseMessage, error) {
	// Aqui, você pode implementar a lógica para processar a mensagem recebida.

	// Por enquanto, apenas retornaremos a mensagem recebida como confirmação.
	log.Printf("Received message from player %s: %s", message.PlayerId, message.Content)

	response := &ResponseMessage{
		Message: fmt.Sprintf("Olá, %s! Tudo bem?", message.PlayerId),
	}

	// Retorne a mensagem recebida como resposta
	return response, nil
}

func (s *Server) StreamMessages(stream MessageService_StreamMessagesServer) error {
	log.Printf("StreamMessages starting")

	// Obter o endereço IP do cliente como chave
	peer, ok := peer.FromContext(stream.Context())
	if !ok {
		return errors.New("não foi possível obter informações do peer")
	}
	clientKey := peer.Addr.String()

	// Logar quando um cliente se conecta
	log.Printf("Cliente conectado ao StreamMessages do endereço: %v", peer.Addr)

	// Armazenar o stream usando o endereço IP como chave
	s.streams.Store(clientKey, stream)
	defer s.streams.Delete(clientKey)

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// Final do stream
			// Cliente desconectou
			log.Println("Cliente desconectado")
			return nil
		}
		if err != nil {
			log.Printf("Erro ao receber mensagem: %v", err)
			return err
		}

		log.Printf("Mensagem recebida de %s: %s", in.PlayerId, in.Content) // Log quando uma mensagem é recebida
		// Processar a mensagem recebida
		s.processMessage(in)

		// Opcional: enviar uma resposta para o remetente, se necessário
		// ...
	}
}

func (s *Server) processMessage(message *PlayerMessage) {
	s.streams.Range(func(key, value interface{}) bool {
		stream, ok := value.(MessageService_StreamMessagesServer)
		if ok {
			if err := stream.Send(&ResponseMessage{Message: message.Content}); err != nil {
				log.Printf("Erro ao enviar mensagem para o cliente %v: %v", key, err)
				// Decida como lidar com o erro. Por exemplo, você pode querer parar de iterar:
				return false
			}
		}
		return true // Continuar iterando
	})
}
