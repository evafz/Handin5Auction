package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	pb "Handin5Auction/grpc"

	"google.golang.org/grpc"
)

type AuctionServer struct {
	pb.UnimplementedAuctionServer

	mu              sync.Mutex
	highestBid      int32
	winningBidder   string
	auctionOver     bool
	bidders         map[string]int32 // Map to store bidders and their bids
	timeToComplete  time.Duration    // Timeframe for the auction to complete
	auctionDeadline time.Time        // Auction deadline
}

func NewAuctionServer(timeToComplete time.Duration) *AuctionServer {
	return &AuctionServer{
		bidders:         make(map[string]int32),
		timeToComplete:  timeToComplete,
		auctionDeadline: time.Now().Add(timeToComplete),
	}
}

func (s *AuctionServer) Bid(ctx context.Context, req *pb.BidRequest) (*pb.BidResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the auction is still active
	if time.Now().Before(s.auctionDeadline) {
		bidderID := req.GetBidderId()
		currentBid := req.GetAmount()

		// Check if the bidder has bid before and if the new bid is higher
		if prevBid, ok := s.bidders[bidderID]; !ok || currentBid > int64(prevBid) {
			s.bidders[bidderID] = int32(currentBid)

			if int64(currentBid) > int64(s.highestBid) {
				s.highestBid = int32(currentBid)
				s.winningBidder = bidderID
			}

			return &pb.BidResponse{Result: pb.BidResponse_BID_SUCCESS}, nil
		}
		return &pb.BidResponse{Result: pb.BidResponse_BID_FAIL}, nil
	}

	return &pb.BidResponse{Result: pb.BidResponse_BID_EXCEPTION}, nil
}

func (s *AuctionServer) Result(ctx context.Context, req *pb.ResultRequest) (*pb.ResultResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the auction is still active
	if time.Now().Before(s.auctionDeadline) {
		return &pb.ResultResponse{Outcome: pb.ResultResponse_AUCTION_NOT_OVER, HighestBid: s.highestBid}, nil
	}

	// Auction is over, provide the result
	return &pb.ResultResponse{Outcome: pb.ResultResponse_AUCTION_SUCCESS, HighestBid: s.highestBid, HighestBidder: s.winningBidder}, nil
}

func startServer(server *AuctionServer) {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started auction server at port 50051\n")

	pb.RegisterAuctionServer(grpcServer, server)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	server := NewAuctionServer(100 * time.Second) // Set the auction duration (e.g., 100 seconds)
	startServer(server)
}

func (s *AuctionServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	logMessage := "Bidder with id " + req.BidderId + " joined"
	log.Printf(logMessage)

	n := strconv.Itoa(int(req.NodeId))
	message := "You joined node " + n + " and can now bid"
	return &pb.JoinResponse{WelcomeMessage: message}, nil
}

func (s *AuctionServer) Leave(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	message := "Bidder with id " + req.BidderId + " left"
	log.Printf(message)

	s.mu.Lock()
	delete(s.bidders, req.BidderId) // Remove client from active clients list
	s.mu.Unlock()

	return &pb.LeaveResponse{}, nil
}
