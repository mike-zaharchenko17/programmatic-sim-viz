package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func BidRequestProducer(ctx context.Context, bidRequestChan chan *BidRequest, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			close(bidRequestChan)
			return
		case bidRequestChan <- generateBidRequest():
			time.Sleep(time.Second)
		}
	}
}

// consumer of BidRequest, producer of BidResponse- intermediate step
func BidRequestResponsePipe(ctx context.Context, bidRequestChan chan *BidRequest, bidResponseChan chan []*BidResponse, wg *sync.WaitGroup) {
	var x *BidRequest
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case x = <-bidRequestChan:
			bidResponseChan <- generateBidResponses(x)
			data, _ := json.MarshalIndent(&x, "", "   ")
			fmt.Println("Received data on bidRequestChannel:\n" + string(data))
		}
	}
}

func BidResponseConsumer(ctx context.Context, bidResponseChan chan []*BidResponse, auctionResultChan chan AuctionResult, wg *sync.WaitGroup) {
	var resArray []*BidResponse

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case resArray = <-bidResponseChan:
			data, _ := json.MarshalIndent(&resArray, "", "  ")
			fmt.Printf("Received data on bidResponseChannel:\n%s\n", string(data))

			auctionResult := simulateAuction(resArray)

			// gracefully log and discard if nothing is listening (i.e., ws not open)
			select {
			case auctionResultChan <- auctionResult:
				auctionResultData, _ := json.MarshalIndent(&auctionResult, "", " ")
				fmt.Printf("(successfully sent) Auction Result: %s\n", auctionResultData)
			default:
				auctionResultData, _ := json.MarshalIndent(&auctionResult, "", " ")
				fmt.Printf("(default block; discarded) Auction Result: %s\n", auctionResultData)
			}
		}
	}
}
