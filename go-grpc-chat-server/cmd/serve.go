/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/vinemesh/go-grpc-chat-server/internal/chat"
	"github.com/vinemesh/go-grpc-chat-server/internal/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor gRPC",
	Long:  `Inicia o servidor gRPC que lida com mensagens.`,
	Run: func(cmd *cobra.Command, args []string) {
		kp, err := kafka.NewProducer("kafka:19092")
		if err != nil {
			log.Fatalf("Erro ao iniciar o produtor do kafka")
		}
		log.Printf("Produtor do kafka iniciado na porta :9092")

		grpcServer := grpc.NewServer()
		chat.RegisterChatServiceServer(grpcServer, &chat.Server{
			Producer: kp,
		})

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
