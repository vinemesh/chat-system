/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"github.com/vinemesh/go-grpc-chat-server/internal/chat"
	"google.golang.org/grpc"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Inicia o servidor gRPC",
	Long:  `Inicia o servidor gRPC que lida com mensagens.`,
	Run: func(cmd *cobra.Command, args []string) {
		grpcServer := grpc.NewServer()
		chat.RegisterChatServiceServer(grpcServer, &chat.Server{})

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
