package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	proto "Replication/protoFile"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAuctionServer

	name       string
	port       int
	lamport    int64
	auctionNode AuctionNode
	clientLock sync.RWMutex
}

func NewServer(name string, port int) *Server {
	return &Server{
		name:    name,
		port:    port,
	}
}

func startServer(server *Server) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	proto.RegisterAuctionServer(grpcServer, server)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	server := NewServer("ChittyChatServer", 8080) // Change port if needed
	startServer(server)
}