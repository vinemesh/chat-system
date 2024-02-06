/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/cobra"
	"github.com/vinemesh/go-grpc-chat-server/internal/chat"
	kp "github.com/vinemesh/go-grpc-chat-server/internal/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor gRPC",
	Long:  `Inicia o servidor gRPC que lida com mensagens.`,
	Run: func(cmd *cobra.Command, args []string) {
		kp, err := kp.NewProducer("kafka:19092")
		if err != nil {
			log.Fatalf("Erro ao iniciar o produtor do kafka")
			os.Exit(1)
		}
		log.Printf("Produtor do kafka iniciado na porta :9092")

		server := &chat.Server{
			Producer: kp,
			ErrChan:  make(chan error, 100),
		}

		grpcServer := grpc.NewServer()
		chat.RegisterChatServiceServer(grpcServer, server)

		go func() {
			for e := range kp.Producer.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						log.Printf("Delivery failed: %v", ev.TopicPartition.Error)
					} else {
						log.Printf("Delivered message to topic %s [%d] at offset %v",
							*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
					}
				}
			}
		}()

		go func() {
			for err := range server.ErrChan {
				// Trate o erro, por exemplo, registrando-o
				log.Printf("Error from chat server: %v", err)
			}
		}()

		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)

		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
		log.Printf("Servidor gRPC iniciado na porta :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Erro ao servir: %v", err)
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			log.Println("Signal interrupt received. Shutting down...")
			kp.Producer.Close()
			grpcServer.GracefulStop()
			close(server.ErrChan)
			os.Exit(0)
		}()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
