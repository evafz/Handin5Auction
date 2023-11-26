package main

import (
	"fmt"
	"sync"
	//proto "Handin5Auction/protoFile"
	//"google.golang.org/grpc"
)

type AuctionNode struct {
	id            int
	highestBid    int
	highestBidder string
	bidders       map[string]struct {
		amount  int64
		lamTime int64
	}
	mu      sync.Mutex
	result  string
	lamTime int64
}

type BidRequest struct {
	bidderID string
	amount   int64
	lamTime  int64
}

type bidResult int64

type BidResponse struct {
	bidResult bidResult
}

const (
	BidFail bidResult = iota
	BidSuccess
	BidException
)

type ResultRequest struct {
	bidderID string
}

type Outcome int64

const (
	AUCTION_NOT_OVER Outcome = iota
	AUCTION_SUCCESS
	AUCTION_FAIL
)

type ResultResponse struct {
	outcome       Outcome
	highestBid    int64
	highestBidder string
}

func (n *AuctionNode) Bid(request *BidRequest) *BidResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Check if the bidder is registered
	if _, exists := n.bidders[request.bidderID]; !exists {
		n.bidders[request.bidderID] = struct {
			amount  int64
			lamTime int64
		}{0, 0}
	}

	// Check if the bid is higher than the previous one
	if request.amount <= n.bidders[request.bidderID].amount {
		return &BidResponse{bidResult: BidFail}
	}

	// Update the bid
	n.bidders[request.bidderID] = struct {
		amount  int64
		lamTime int64
	}{request.amount, request.lamTime}

	n.lamTime = max(n.lamTime, request.lamTime) + 1
	return &BidResponse{bidResult: BidSuccess}
}

func (n *AuctionNode) ResultMethod() Outcome {
	resultResponse := n.result

	switch resultResponse {
	case "AUCTION_NOT_OVER":
		fmt.Printf("Auction is still ongoing and the highest bid is %d", n.highestBid)
		return AUCTION_NOT_OVER
	case "AUCTION_SUCCESS":
		fmt.Printf("Auction won by %s with a bid of %d\n", n.highestBidder, n.highestBid)
		return AUCTION_SUCCESS
	case "AUCTION_FAIL":
		fmt.Println("Auction failed")
		return AUCTION_FAIL
	}
}

func main() {
	node := &AuctionNode{
		id: 1,
		bidders: make(map[string]struct {
			amount  int64
			lamTime int64
		}),
		result: "AUCTION_NOT_OVER",
	}

	//some bids:
	bidRequest1 := &BidRequest{amount: 50, bidderID: "Alex"}
	bidRequest2 := &BidRequest{amount: 55, bidderID: "Bjarne"}
	bidRequest3 := &BidRequest{amount: 150, bidderID: "ChristHimself"}

	node.Bid(bidRequest1)
	node.Bid(bidRequest2)
	node.Bid(bidRequest3)

	node.ResultMethod()
}
