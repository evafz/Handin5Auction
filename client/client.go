package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	pb "Handin5Auction/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuctionClient(conn)

	main.Join
	/*fmt.Print("Enter your bidder ID: ")
	scanner := bufio.NewReader(os.Stdin)
	bidderID, _ := scanner.ReadString('\n')
	bidderID = strings.TrimSpace(bidderID) */

	// Join the auction
	joinResp, err := client.Join(context.Background(), &pb.JoinRequest{BidderId: bidderID})
	if err != nil {
		log.Fatalf("Error while joining the auction: %v", err)
	}
	fmt.Println(joinResp.WelcomeMessage)

	// Start bidding process
	for {
		fmt.Print("Enter your bid amount (or type 'exit' to leave): ")
		amountStr, _ := scanner.ReadString('\n')
		amountStr = strings.TrimSpace(amountStr)

		if amountStr == "exit" {
			break
		}

		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			log.Println("Invalid bid amount. Please enter a valid number.")
			continue
		}

		// Bid in the auction
		bidResp, err := client.Bid(context.Background(), &pb.BidRequest{Amount: int64(amount), BidderId: bidderID})
		if err != nil {
			log.Fatalf("Error while bidding: %v", err)
		}

		switch bidResp.Result {
		case pb.BidResponse_BID_SUCCESS:
			fmt.Println("Bid successful!")
		case pb.BidResponse_BID_FAIL:
			fmt.Println("Bid failed. Try a higher amount.")
		case pb.BidResponse_BID_EXCEPTION:
			fmt.Println("Exception occurred during bidding.")
		}

		// Query auction result
		resultResp, err := client.Result(context.Background(), &pb.ResultRequest{BidderId: bidderID})
		if err != nil {
			log.Fatalf("Error while querying result: %v", err)
		}

		switch resultResp.Outcome {
		case pb.ResultResponse_AUCTION_NOT_OVER:
			fmt.Printf("Auction is not over. Current highest bid: %d\n", resultResp.HighestBid)
		case pb.ResultResponse_AUCTION_SUCCESS:
			fmt.Printf("Auction result: Winning bid is %d\n", resultResp.HighestBid)
		}
	}

	// Leave the auction
	leaveResp, err := client.Leave(context.Background(), &pb.LeaveRequest{BidderId: bidderID})
	if err != nil {
		log.Fatalf("Error while leaving the auction: %v", err)
	}
	fmt.Println(leaveResp.GoodbyeMessage)
}
