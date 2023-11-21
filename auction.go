package main

import (
	"fmt"
	"sync"
)

type AuctionNode struct {
	id            int
	highestBid    int
	highestBidder string
	bidders       map[string]int
	mu            sync.Mutex
}


func (n *AuctionNode) Bid(request *BidRequest) *BidResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the bidder is registered
	if _, exists := n.bidders[bidderID]; !exists {
		n.bidders[bidderID] = 0
	}

	// Check if the bid is higher than the previous one
	if amount <= n.bidders[bidderID] {
		return BidFail
	}

	// Update the bid
	n.bidders[bidderID] = amount

	return BidSuccess
}

func main(){
	node := &AuctionNode{
		id:	1,
		bidders: make(map[string]int),
	}

 //some bids:
	bidRequest1 := &BidRequest{Amount: 50, BidderId: "Alex"}
	bidRequest2 := &BidRequest{Amount: 55, BidderId: "Bjarne"}
	bidRequest3 := &BidRequest{Amount: 150, BidderId: "ChristHimself"}
	
	response1 := node.Bid(&bidRequest1)
	response2 := node.Bid(&bidRequest2)
	response3 := node.Bid(&bidRequest3)

	resultRequest := &ResultRequest{}
	resultResponse := node.Result(resultRequest)

	switch resultResponse.Outcome {
	case ResultResponse_AUCTION_NOT_OVER:
		fmt.Println("Auction is still ongoing")
	case ResultResponse_AUCTION_SUCCESS:
		fmt.Printf("Auction won by %s with a bid of %d\n", resultResponse.HighestBidder, resultResponse.HighestBid)
	case ResultResponse_AUCTION_FAIL:
		fmt.Println("Auction failed")
	}
}
