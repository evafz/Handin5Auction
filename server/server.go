package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"

	proto "Handin5Auction/protoFile"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedAuctionServer

	name        string
	port        int
	lamport     int64
	auctionNode client.AuctionNode
	clientLock  sync.RWMutex
}

func NewServer(name string, port int) *Server {
	return &Server{
		name: name,
		port: port,
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

func (s *Server) Bid(ctx context.Context, request *proto.BidRequest) (*proto.BidResponse, error){
	s.clientLock.Lock()
	defer s.clientLock.Unlock()

	//New BidRequest
	bidRequest := &proto.BidRequest{
		bidderID: request.BidderId,
		amount: request.Amount,
		lamTime: s.auctionNode.lamTime,
	}

	response := s.auctionNode.Bid(bidRequest)

	return &proto.BidResponse{Result: proto.BidResponse_BidResult(response.result)}, nil
}

func main() {
	server := NewServer("AuctionServer", 8080) // Change port if needed
	startServer(server)
}
